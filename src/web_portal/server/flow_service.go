package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"config"
	"utils/helper"
	"utils/conv"
	"lua_helper"
	"scheduler"
	"model"
	"model/rpc_data"
	"middleware/remote_utils"
	"middleware/microservice"

	// "github.com/gorilla/rpc"
	// "github.com/gorilla/rpc/json"
	log "github.com/cihub/seelog"
	"github.com/go-sql-driver/mysql"
)

type FlowDataArgs struct {
	Id            string
	Name          string
	Description   string
	Creator       string
	StartupScript string
	Tasks         []TaskData
	DeleteTaskIds []string
}

type TaskData struct {
	Id          string
	Name        string
	Script      string
	MaxRetries  interface{} `json:"max_retries"`
	Description string
}

type AddFlowReply struct {
	Message string
}

type FlowService struct{}

type GetFlowInstInfoArgs struct {
	FlowId     string
	PId       string
	RunningDay string
}

func (h *FlowService) AddFlow(r *http.Request, args *FlowDataArgs, reply *AddFlowReply) error {
	flow := &model.Flow{Id: args.Id, Name: args.Name, Description: args.Description, StartupScript: args.StartupScript, Creator: args.Creator}

	tasks := make([]*model.Task, 0, len(args.Tasks))
	for i, taskdata := range args.Tasks {
		task := &model.Task{}
		task.Id = taskdata.Id
		task.Name = taskdata.Name
		task.Description = taskdata.Description
		task.OrderId = i
		task.FlowId = flow.Id
		task.MaxRetries, _ = conv.Int(taskdata.MaxRetries)
		script := taskdata.Script
		task.Script = &script
		tasks = append(tasks, task)
	}
	err := model.AddFlow(flow, tasks)
	if err != nil {
		if dberr, match := err.(*mysql.MySQLError); match {
			if dberr.Number == 1062 {
				return errors.New("Flow Id " + args.Id + " is exists.")
			}
		}
		log.Error(err)
		return err
	}
	return nil
}

func (h *FlowService) UpdateFlow(r *http.Request, args *FlowDataArgs, reply *AddFlowReply) error {
	flow := &model.Flow{Id: args.Id, Name: args.Name, Description: args.Description, StartupScript: args.StartupScript}

	tasks := make([]*model.Task, 0, len(args.Tasks))
	for i, taskdata := range args.Tasks {
		task := &model.Task{}
		task.Id = taskdata.Id
		task.Name = taskdata.Name
		task.Description = taskdata.Description
		task.OrderId = i
		task.FlowId = flow.Id
		task.MaxRetries, _ = conv.Int(taskdata.MaxRetries)
		script := taskdata.Script
		task.Script = &script
		tasks = append(tasks, task)
	}

	err := model.UpdateFlow(flow, tasks, args.DeleteTaskIds)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (h *FlowService) ReportState1(r *http.Request, args *rpc_data.StateDataArgs, reply *rpc_data.StateReply) error {
	return nil
}

func (h *FlowService) ReportState(r *http.Request, args *rpc_data.StateDataArgs, reply *rpc_data.StateReply) error {
	day, err := time.Parse("2006-01-02", args.Date)
	if err != nil {
		log.Error(err)
		return err
	}
	timestamp := time.Now().UnixNano()
	if r != nil {
		if t, err := strconv.ParseInt(r.Header.Get("X-Timestamp"), 10, 64); err == nil {
			timestamp = t
		}
	}

	report_log := model.StateReportLog{}
	report_log.PId = args.PId
	report_log.Key = args.Key
	report_log.FlowId = args.FlowId
	report_log.TaskId = args.TaskId
	report_log.State = args.State
	report_log.Creator = args.Creator
	report_log.Date = day
	extra_data_json, _ := json.Marshal(args.ExtraData)
	report_log.ExtraData = (string)(extra_data_json)

	log.Info(report_log)
	_, err = model.AddStateReportLog(&report_log)
	if err != nil {
		log.Error(err)
	}

	id, err := model.SaveState(args.PId, args.Key, args.FlowId, args.TaskId, args.State, args.Creator, day, args.ExtraData, timestamp)
	if err != nil {
		log.Error(err)
	}
	report_log.FlowInstId = (int)(id)

	go scheduler.OnStateChange(&report_log)
	return err
}

func (h *FlowService) RerunTask(r *http.Request, args *rpc_data.RerunTaskArgs, reply *rpc_data.RerunTaskReply) error {
	return lua_helper.RerunTask(args.FlowInstId, args.TaskId, args.SingleTask, args.Creator)
}

func (h *FlowService) StartFlow(r *http.Request, args *rpc_data.StartFlowArgs, reply *rpc_data.StartFlowReply) error {
	date, ok := conv.Time(args.Date)
	if !ok {
		return log.Error("Parse date(yyyy-MM-dd) failed")
	}
	var err error
	if len(args.TaskId) == 0 {
		reply.FlowInstId, err = lua_helper.StartFlow(args.Id, args.PId, args.Key, date, args.Creator, &args.StartupScript)
	} else {
		reply.FlowInstId, err = lua_helper.StartFlowFromTask(args.Id, args.TaskId, args.PId, args.Key, date, args.Creator, &args.StartupScript)
	}
	return err
}

// for ssh
func (h *FlowService) GetRemoteLog(r *http.Request, args *rpc_data.GetRemoteLogArgs, reply *rpc_data.GetRemoteLogReply) error {
	host := args.Host
	uuid := args.Uuid
	date := args.Date

	// reply.Cmdline = host
	// reply.Output = uuid
	// reply.Error = date
	// return nil

	server, err := model.GetServerByHost(host)
	if err != nil {
		return log.Error("Find host "+host+" failed:", err)
	}
	if server == nil {
		return log.Error("Host " + host + " not found in our database!")
	}
	// password := server.Password
	ssh_client, err := remote_utils.Connect(fmt.Sprintf("%s:%d", server.Host, server.Port), *server.Username, *server.Password)
	if err != nil {
		return log.Error("Connect to host "+host+" failed:", err)
	}
	defer ssh_client.Close()

	path := fmt.Sprintf("%s/proc/%s/%s", *config.ServerRoot, strings.Replace(date, "-", "", -1), uuid)

	log_filenames := []string{path + "/cmd", path + "/out.log", path + "/err.log"}

	log_datas := remote_utils.ReadFiles(ssh_client, log_filenames)

	log_strings := make([]string, len(log_datas))
	for i, log_data := range log_datas {
		switch val := log_data.(type) {
		case []byte:
			log_strings[i] = (string)(val)
		case error:
			log_strings[i] = "<p class='text-danger'>" + val.Error() + "</p>"
		}
	}

	reply.Cmdline = log_strings[0]
	reply.Output = log_strings[1]
	reply.Error = log_strings[2]
	return nil
}

// for agent
func (h *FlowService) GetRemoteLog1(r *http.Request, args *rpc_data.GetRemoteLogArgs, reply *rpc_data.GetRemoteLogReply) error {
	host := args.Host
	uuid := args.Uuid
	date := args.Date

	// reply.Cmdline = host
	// reply.Output = uuid
	// reply.Error = date
	// return nil
	path := fmt.Sprintf("%s/proc/%s/%s", *config.ServerRoot, strings.Replace(date, "-", "", -1), uuid)

	var err error
	var cmd, stdout, stderr interface{}
	err = microservice.SafeRequest(host, "Script.Exec", "return exec.CombinedOutput('tail', '-n', '6000', '"+path+"/cmd"+"')", &cmd)
	if err != nil {
		reply.Cmdline = "<p class='text-danger'>" + err.Error() + "</p>"
	} else {
		reply.Cmdline = conv.String(cmd)
	}

	err = microservice.SafeRequest(host, "Script.Exec", "return exec.CombinedOutput('tail', '-n', '6000', '"+path+"/out.log"+"')", &stdout)
	if err != nil {
		reply.Output = "<p class='text-danger'>" + err.Error() + "</p>"
	} else {
		reply.Output = conv.String(stdout)
	}

	err = microservice.SafeRequest(host, "Script.Exec", "return exec.CombinedOutput('tail', '-n', '6000', '"+path+"/err.log"+"')", &stderr)
	if err != nil {
		reply.Error = "<p class='text-danger'>" + err.Error() + "</p>"
	} else {
		reply.Error = conv.String(stderr)
	}

	return nil
}

func (h *FlowService) RunScript(r *http.Request, args *rpc_data.RunScriptArgs, reply *rpc_data.RunScriptReply) error {
	script := args.Script

	log.Debug("Run script:", script)

	L := lua_helper.GetState()
	defer lua_helper.RevokeState(L)

	err := L.DoString("in_terminal=true")
	if err != nil {
		return err
	}
	err = L.DoString(script)
	if err != nil {
		return err
	}

	output := L.GetOutput()
	log.Info("Script output:", output)
	reply.Output = output

	return nil
}

func (h *FlowService) KillTaskInstance(r *http.Request, args *rpc_data.KillTaskInstanceArgs, reply *int) error {
	task_inst, err := model.GetTaskInstById(args.FlowInstId, args.TaskId)
	if err != nil || task_inst == nil {
		return log.Error("Find task intance failed:", err)
	}

	if *task_inst.State != model.StateRunning {
		return nil
	}

	if task_inst.RunningTime == nil || len(task_inst.RemoteExecHost) == 0 || len(task_inst.RemoteExecUuid) == 0 {
		return nil
	}

	model.UpdateTaskInstRetries(args.FlowInstId, args.TaskId, 999999)

	// err = lua_helper.RemoteKill(true, task_inst.RemoteExecHost, *task_inst.RunningTime, task_inst.RemoteExecUuid)
	err = lua_helper.RemoteKill(task_inst.RemoteExecHost, *task_inst.RunningTime, task_inst.RemoteExecUuid)

	if err != nil {
		return err
	}

	err = model.UpdateFlowInstStateById(args.FlowInstId, model.StateFailed)
	return model.UpdateTaskInstState(args.FlowInstId, args.TaskId, model.StateFailed)
}

func (h *FlowService) SetTaskInstanceSuccess(r *http.Request, args *rpc_data.SetTaskInstanceSuccessArgs, reply *int) error {
	task_inst, err := model.GetTaskInstById(args.FlowInstId, args.TaskId)
	if err != nil || task_inst == nil {
		return log.Error("Find task intance failed:", err)
	}
	flow_inst, err := model.GetFlowInstById(args.FlowInstId)
	if err != nil || flow_inst == nil {
		return log.Error("Find flow intance failed:", err)
	}

	report_log := model.StateReportLog{}
	report_log.PId = flow_inst.PId
	report_log.Key = flow_inst.Key
	report_log.FlowId = flow_inst.FlowId
	report_log.TaskId = args.TaskId
	report_log.State = model.StateSucceed
	report_log.Creator = "[System]"
	report_log.Date = flow_inst.RunningDay
	extra_data := map[string]string{}
	extra_data_json, _ := json.Marshal(extra_data)
	report_log.ExtraData = (string)(extra_data_json)

	_, err = model.SaveState(flow_inst.PId, flow_inst.Key, flow_inst.FlowId, args.TaskId,
		model.StateSucceed, "[System]", flow_inst.RunningDay, extra_data, time.Now().UnixNano())
	if err != nil {
		log.Error(err)
	}
	report_log.FlowInstId = args.FlowInstId

	go scheduler.OnStateChange(&report_log)

	return nil
}

func (h *FlowService) SetFlowInstanceSuccess(r *http.Request, args *rpc_data.SetFlowInstanceSuccessArgs, reply *int) error {
	flow_inst, err := model.GetFlowInstById(args.FlowInstId)
	if err != nil || flow_inst == nil {
		return log.Error("Find flow intance failed:", err)
	}

	if flow_inst.State != model.StateFailed {
		return log.Error("Flow intance is not failed.")
	}

	err = model.SetFlowInstStateSucceed(args.FlowInstId)

	helper.Notify("flow.schedule", struct{}{})
	return nil
}

func (h *FlowService) TaskAlarm(r *http.Request, args *rpc_data.TaskAlarmArgs, reply *int) error {
	alarm := &model.TaskInstAlarm{FlowInstId: args.FlowInstId, TaskId: args.TaskId, Content: args.Content}
	_, err := model.AddTaskInstAlarm(alarm)
	if err != nil {
		return log.Error("Insert into database failed:", err)
	}

	return nil
}

func (h *FlowService) GetFlowInstInfo(r *http.Request, args *GetFlowInstInfoArgs, reply *model.FlowInst) error {
	flow_inst, err := model.GetFlowInstByName(args.FlowId, args.PId, args.RunningDay)
	if err != nil {
		return log.Error("Find flow intance failed", err)
	}
	if flow_inst == nil {
		return log.Error("Find flow intance failed", err)
	}
	*reply = *flow_inst
	return nil
}
