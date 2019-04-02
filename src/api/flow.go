package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"errors"
	"encoding/json"
	"strings"
	"time"

	"config"
	
	"lua_helper"
	"model"
	"model/rpc_data"
	"utils/conv"
	"utils/helper"
	"scheduler"
	"middleware/remote_utils"

	log "github.com/cihub/seelog"
	"github.com/go-sql-driver/mysql"
)

func init() {
	SetRouterRegister(func(router *RouterGroup) {
		flowRouteGroup := router.Group("/oneflow/")
		flowRouteGroup.StdPOST("AddFlow", AddFlow)
		flowRouteGroup.StdPOST("UpdateFlow", UpdateFlow)
		flowRouteGroup.StdPOST("StartFlow", StartFlow)
		flowRouteGroup.StdPOST("RerunTask", RerunTask)
		flowRouteGroup.StdPOST("KillTaskInstance", KillTaskInstance)
		flowRouteGroup.StdPOST("SetTaskInstanceSuccess", SetTaskInstanceSuccess)
		flowRouteGroup.StdPOST("SetFlowInstanceSuccess", SetFlowInstanceSuccess)
		flowRouteGroup.StdPOST("GetRemoteLog", GetRemoteLog)
		
		flowRouteGroup.StdPOST("API", RpcAPI)
		flowRouteGroup.StdGET("GetFlows", GetFlows)
	})
}

/*
{"id":2,"method":"FlowService.UpdateFlow",
"params":[{"id":"NEW_FLOW","name":"新流程","creator":"helight","description":"","startupScript":"",
	"tasks":[{"id":"TASK_1","name":"任务1","description":"","max_retries":"0","script":""},
		{"id":"TASK_2","name":"任务2","description":"","max_retries":"0","script":""}],
		"deleteTaskIds":[]}]}
api	
*/
func UpdateFlow(c *Context) (code int, message string, data interface{}) {
	//1.解析参数
	newFlow, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}

	c.Info("post: --> ", string(newFlow), " <--data")
	var args rpc_data.FlowDataArgs
	err = json.Unmarshal([]byte(newFlow), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 

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

	err = model.UpdateFlow(flow, tasks, args.DeleteTaskIds)
	if err != nil {
		log.Error(err)
		return -1, err.Error(), "err"
	}
	return 0, "ok", nil
}
/*
{"id":2,"method":"FlowService.AddFlow",
"params":[{"id":"NEW_FLOW3","name":"新流程","creator":"helight","description":null,"startupScript":"",
"tasks":[{"id":"TASK_1","name":"任务1","description":"","max_retries":"0","script":""},
{"id":"TASK_2","name":"任务2","description":"","max_retries":"0","script":""},
{"id":"TASK_3","name":"任务3","description":"","max_retries":"0","script":""}],
"deleteTaskIds":[]}]}
*/
// func (h *FlowService) AddFlow(r *http.Request, args *FlowDataArgs, reply *AddFlowReply) error {
func AddFlow(c *Context) (code int, message string, data interface{}) {
	//1.解析参数
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.FlowDataArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
        return -1, err.Error(), "Do Unmarshal err"
	} 

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
	err = model.AddFlow(flow, tasks)
	if err != nil {
		if dberr, match := err.(*mysql.MySQLError); match {
			if dberr.Number == 1062 {
				return -1, "Flow Id " + args.Id + " is exists.", errors.New("Flow Id " + args.Id + " is exists.")
			}
		}
		log.Error(err)
		return 50001, err.Error(), nil
	}
	return 0, "ok", nil
}


// func (h *FlowService) StartFlow(r *http.Request, args *rpc_data.StartFlowArgs, reply *rpc_data.StartFlowReply) error {
func StartFlow(c *Context) (code int, message string, data interface{}) {
	//1.解析参数
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.StartFlowArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	date, ok := conv.Time(args.Date)
	if !ok {
		return  -1, "conv.Time(args.Date) err", log.Error("Parse date(yyyy-MM-dd) failed")
	}
	// var err error
	var reply rpc_data.StartFlowReply
	if len(args.TaskId) == 0 {
		reply.FlowInstId, err = lua_helper.StartFlow(args.Id, args.PId, args.Key, date, args.Creator, &args.StartupScript)
	} else {
		reply.FlowInstId, err = lua_helper.StartFlowFromTask(args.Id, args.TaskId, args.PId, args.Key, date, args.Creator, &args.StartupScript)
	}
	if err != nil {
		return -1, err.Error(), ""
	}
	return 0, "ok", reply
}

// func (h *FlowService) RerunTask(r *http.Request, args *rpc_data.RerunTaskArgs, reply *rpc_data.RerunTaskReply) error {
func RerunTask(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.RerunTaskArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	lua_helper.RerunTask(args.FlowInstId, args.TaskId, args.SingleTask, args.Creator)
	return 0, "ok", ""
}

// begin_date=2019-03-29&end_date=2019-03-31&set_nu=all&process_type=all&process_state_type=-1&pid=all
func GetFlows(c *Context) (code int, message string, data interface{}) {
	begin_date := c.Query("begin_date")
	end_date := c.Query("end_date")
	set_nu := c.Query("set_nu")
	flow_id := c.Query("process_type")
	process_state_type := c.Query("process_state_type")
	pid := c.Query("pid")

	cond1 := " (fi.running_day between '"+begin_date+" 00:00:00' and '"+end_date+" 23:59:59')"
	cond2 := " and ('"+set_nu+"' = 'all') and ('"+flow_id+"' = 'all' or f.id = '"+flow_id+"')"
	cond3 := " and ('"+process_state_type+"' = 'all' or fi.state = "+process_state_type+" or ("+process_state_type+" = -1 and fi.`state` in (0,1,3)))"
	cond4 := " and ('"+pid+"' = 'all' or fi.pid = '"+pid+"') ORDER BY fi.running_day desc, fi.flow_id, last_update_time DESC "

	cond := cond1 + cond2 + cond3 +cond4
	flows, err := model.GetFlowInstAll(cond)
	if err != nil {
		return -1, "get db error", ""
	}
	return -1, "ok", flows
}


// for ssh
// func (h *FlowService) GetRemoteLog(r *http.Request, args *rpc_data.GetRemoteLogArgs, reply *rpc_data.GetRemoteLogReply) error {
func GetRemoteLog(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.GetRemoteLogArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	host := args.Host
	uuid := args.Uuid
	date := args.Date

	// reply.Cmdline = host
	// reply.Output = uuid
	// reply.Error = date
	// return nil

	server, err := model.GetServerByHost(host)
	if err != nil {
		return -1, err.Error(), log.Error("Find host "+host+" failed:", err)
	}
	if server == nil {
		return -1, err.Error(), log.Error("Host " + host + " not found in our database!")
	}
	// password := server.Password
	ssh_client, err := remote_utils.Connect(fmt.Sprintf("%s:%d", server.Host, server.Port), *server.Username, *server.Password)
	if err != nil {
		return -1, err.Error(), log.Error("Connect to host "+host+" failed:", err)
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

	var reply rpc_data.GetRemoteLogReply
	reply.Cmdline = log_strings[0]
	reply.Output = log_strings[1]
	reply.Error = log_strings[2]
	return 0, "ok", reply
}


// func (h *FlowService) RunScript(r *http.Request, args *rpc_data.RunScriptArgs, reply *rpc_data.RunScriptReply) error {
func RunScript(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.RunScriptArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	script := args.Script

	log.Debug("Run script:", script)

	L := lua_helper.GetState()
	defer lua_helper.RevokeState(L)

	err = L.DoString("in_terminal=true")
	if err != nil {
		return -1, err.Error(), ""
	}
	err = L.DoString(script)
	if err != nil {
		return -1, err.Error(), ""
	}

	var reply rpc_data.RunScriptReply
	output := L.GetOutput()
	log.Info("Script output:", output)
	reply.Output = output

	return 0, "ok", reply
}

// func (h *FlowService) SetTaskInstanceSuccess(r *http.Request, args *rpc_data.SetTaskInstanceSuccessArgs, reply *int) error {
func SetTaskInstanceSuccess(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.SetTaskInstanceSuccessArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	task_inst, err := model.GetTaskInstById(args.FlowInstId, args.TaskId)
	if err != nil || task_inst == nil {
		log.Error("Find task intance failed:", err)
		return -1, "Find task intance failed:" + err.Error(), ""
	}
	flow_inst, err := model.GetFlowInstById(args.FlowInstId)
	if err != nil || flow_inst == nil {
		log.Error("Find flow intance failed:", err)
		return -1, "Find flow intance failed:" + err.Error(), ""
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

	return 0, "ok", ""
}

// func (h *FlowService) SetFlowInstanceSuccess(r *http.Request, args *rpc_data.SetFlowInstanceSuccessArgs, reply *int) error {
func SetFlowInstanceSuccess(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.SetFlowInstanceSuccessArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	flow_inst, err := model.GetFlowInstById(args.FlowInstId)
	if err != nil || flow_inst == nil {
		log.Error("Find flow intance failed:", err)
		return -1, "Find flow intance failed:" + err.Error(), ""
	}

	if flow_inst.State != model.StateFailed {
		log.Error("Flow intance is not failed.")
		return -1, "Flow intance is not failed.", ""
	}

	err = model.SetFlowInstStateSucceed(args.FlowInstId)

	helper.Notify("flow.schedule", struct{}{})
	return 0, "ok", ""
}

// func (h *FlowService) KillTaskInstance(r *http.Request, args *rpc_data.KillTaskInstanceArgs, reply *int) error {
func KillTaskInstance(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(reqbody), " <--data")
	var args rpc_data.KillTaskInstanceArgs
	err = json.Unmarshal([]byte(reqbody), &args)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	task_inst, err := model.GetTaskInstById(args.FlowInstId, args.TaskId)
	if err != nil || task_inst == nil {
		log.Error("Find task intance failed:", err)
		return -1, "Find task intance failed:" + err.Error(), ""
	}

	if *task_inst.State != model.StateRunning {
		return 0, "ok", ""
	}

	if task_inst.RunningTime == nil || len(task_inst.RemoteExecHost) == 0 || len(task_inst.RemoteExecUuid) == 0 {
		return 0, "ok", ""
	}

	model.UpdateTaskInstRetries(args.FlowInstId, args.TaskId, 999999)

	// err = lua_helper.RemoteKill(true, task_inst.RemoteExecHost, *task_inst.RunningTime, task_inst.RemoteExecUuid)
	err = lua_helper.RemoteKill(task_inst.RemoteExecHost, *task_inst.RunningTime, task_inst.RemoteExecUuid)

	if err != nil {
		return -1, err.Error(), ""
	}

	err = model.UpdateFlowInstStateById(args.FlowInstId, model.StateFailed)
	model.UpdateTaskInstState(args.FlowInstId, args.TaskId, model.StateFailed)
	return 0, "ok", ""
}

func RpcAPI(c *Context) (code int, message string, data interface{}) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)	
	if err != nil {
		return 50001, err.Error(), nil
	}
	resp, err := http.Post("http://localhost/data_flow/api", "application/json; charset=utf-8", strings.NewReader(string(reqbody)))
	if err != nil {
		return -1, "post err", ""
	}

	defer resp.Body.Close()
	repbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "read post err", ""
	}	
	return 0, "ok", repbody
}