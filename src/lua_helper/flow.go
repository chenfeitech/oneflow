package lua_helper

import (
	"strings"
	"sync"
	"time"
	"errors"

	// "fmt"
	"utils/alarm"
	"utils/helper"
	"model"

	log "github.com/cihub/seelog"
	"github.com/stevedonovan/luar"
)

func init() {
	LuaGlobal["FindTaskByFlowId"] = model.FindTaskByFlowId
	LuaGlobal["GetTaskInstById"] = model.GetTaskInstById
	//LuaGlobal["StartFlow"] = StartFlow
}

var (
	g_start_flow_mutex = &sync.Mutex{}
)

type FlontFlowStateRet struct {
	Code     int    `json:"code"`
	CodeInfo string `json:"codeinfo"`
}

func StartFlowInst(flow_inst *model.FlowInst) error {
	log.Info(flow_inst)
	tasks, err := model.FindTaskByFlowId(flow_inst.FlowId)
	if err != nil {
		log.Error("Find tasks failed:", err)
		return (log.Error("Find tasks failed:", err))
	}
	if len(tasks) == 0 {
		model.SaveState(flow_inst.PId, flow_inst.Key, flow_inst.FlowId, "", 3, flow_inst.Creator, flow_inst.RunningDay, make(map[string]string), 0)
		log.Error("Flow tasks empty.")
		return (log.Error("Flow tasks empty."))
	}

	var task *model.Task
	for _, t := range tasks {
		if t.Id == flow_inst.BeginTask || len(flow_inst.BeginTask) == 0 {
			task = t
			break
		}
	}

	if task == nil {
		model.SaveState(flow_inst.PId, flow_inst.Key, flow_inst.FlowId, "", 3, flow_inst.Creator, flow_inst.RunningDay, make(map[string]string), 0)
		log.Error("Flow task ", flow_inst.BeginTask, " cannot found.")
		return (log.Error("Flow task ", flow_inst.BeginTask, " cannot found."))
	}

	return StartFlowTask(flow_inst.PId, flow_inst.Key, task, flow_inst.Creator, flow_inst.RunningDay, nil, nil)
}

func (s *iState) Lua_StartFlow(flow_id string, pid string, key string, date time.Time, creator string, startup_script *string) (flow_inst_id int) {
	flow_inst_id, err := StartFlow(flow_id, pid, key, date, creator, startup_script)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		s.RaiseError(err.Error())
	}
	return flow_inst_id
}

func StartFlow(flow_id string, pid string, key string, date time.Time, creator string, startup_script *string) (flow_inst_id int, err error) {
	g_start_flow_mutex.Lock()
	defer g_start_flow_mutex.Unlock()

	tasks, err := model.FindTaskByFlowId(flow_id)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Find tasks failed:", err))
	}
	if len(tasks) == 0 {
		return 0, (log.Error("Flow tasks empty."))
	}

	if startup_script == nil {
		flow, err := model.GetFlowByKey(flow_id)
		if err != nil {
			log.Error("err: ", err , " flow_id: ", flow_id)
			return 0, (log.Error("Find flow failed:", err))
		}
		startup_script = &flow.StartupScript
	}

	flow_inst, err := model.GetFlowInst4(pid, key, flow_id, date)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Find flow inst failed:", err))
	}

	if flow_inst != nil && (flow_inst.State == model.StateReady || flow_inst.State == model.StateRunning) {
		return 0, (errors.New("Flow instance is running!"))
	}

	flow_inst_id, err = model.AddFlowInst(&model.FlowInst{
		PId:           pid,
		Key:           key,
		FlowId:        flow_id,
		Creator:       creator,
		RunningDay:    date,
		StartupScript: *startup_script,
	})

	err = model.ResetTaskInstStateSinceOrderId(flow_inst_id, flow_id, -1000)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, err
	}

	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Add flow instance failed:", err))
	}

	helper.Notify("flow.schedule", struct{}{})

	return flow_inst_id, nil
}

func (s *iState) Lua_StartFlowFromTask(flow_id string, task_id string, pid string, key string, date time.Time, creator string, startup_script *string) (flow_inst_id int) {
	flow_inst_id, err := StartFlowFromTask(flow_id, task_id, pid, key, date, creator, startup_script)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		s.RaiseError(err.Error())
	}
	return flow_inst_id
}

func StartFlowFromTask(flow_id string, task_id string, pid string, key string, date time.Time, creator string, startup_script *string) (flow_inst_id int, err error) {
	g_start_flow_mutex.Lock()
	defer g_start_flow_mutex.Unlock()

	tasks, err := model.FindTaskByFlowId(flow_id)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Find tasks failed:", err))
	}
	if len(tasks) == 0 {
		return 0, (log.Error("Flow tasks empty."))
	}

	var task *model.Task
	for _, t := range tasks {
		if t.Id == task_id {
			task = t
		}
	}

	if task == nil {
		return 0, (log.Error("Flow task not exists."))
	}

	if startup_script == nil {
		flow, err := model.GetFlowByKey(flow_id)
		if err != nil {
			return 0, (log.Error("Find flow failed:", err))
		}
		startup_script = &flow.StartupScript
	}

	flow_inst, err := model.GetFlowInst4(pid, key, flow_id, date)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Find flow inst failed:", err))
	}

	if flow_inst != nil && (flow_inst.State == model.StateReady || flow_inst.State == model.StateRunning) {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (errors.New("Flow instance is running!"))
	}

	flow_inst_id, err = model.AddFlowInst(&model.FlowInst{
		PId:           pid,
		Key:           key,
		FlowId:        flow_id,
		Creator:       creator,
		RunningDay:    date,
		StartupScript: *startup_script,
		BeginTask:     task_id,
	})
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Add flow instance failed:", err))
	}

	model.ResetTaskInstStateSinceOrderId(flow_inst_id, flow_id, task.OrderId)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		return 0, (log.Error("Add flow instance failed:", err))
	}
	helper.Notify("flow.schedule", struct{}{})

	return flow_inst_id, nil
}

func (s *iState) Lua_StartFlowTask(pid string, key string, task *model.Task, creator string, date time.Time, state_log *model.StateReportLog, task_inst *model.TaskInst) {
	err := StartFlowTask(pid, key, task, creator, date, state_log, task_inst)
	if err != nil {
		log.Error("err: ", err , " task: ", task)
		s.RaiseError(err.Error())
	}
}

func StartFlowTask(pid string, key string, task *model.Task, creator string, date time.Time, state_log *model.StateReportLog, task_inst *model.TaskInst) error {
	log.Info(pid, key, task)
	if task == nil {
		return log.Error("Can not start nil task.")
	}

	product, err := model.GetProductsByKey(pid)
	if err != nil {
		log.Error("err: ", err , " task: ", task)
		return log.Error("Find pid", pid, " from tbProducts failed:", err)
	}
	if product == nil {
		return log.Error("Can not found pid ", pid, "  from tbProducts")
	}

	flow_inst, err := model.GetFlowInst4(pid, key, task.FlowId, date)
	if err != nil || flow_inst == nil {
		log.Error("err: ", err , " task: ", task)
		return log.Error("Find flow instance ", task_inst.FlowInstId, " failed:", err)
	}

	if task.Script != nil && len(strings.Trim(*task.Script, " \r\n\t")) > 0 {
		log.Debug("Run script:", *task.Script)

		L := GetState()
		defer RevokeState(L)
		log.Info(pid, key, task, "Register")
		luar.Register(L.State, "", luar.Map{
			"pid":      pid,
			"date":      date,
			"key":       key,
			"creator":   creator,
			"task":      task,
			"state_log": state_log,
			"task_inst": task_inst,
		})

		L.Task = task
		L.FlowInstance = flow_inst
		L.TaskInstance = task_inst

		var output string
		log.Info(pid, key, task, "get pid info")
		err = L.DoString("product=db_query_dict(\"select * from tbProducts where PId=?\", pid)[1]")
		if err != nil {
			output = output + log.Error("Run init script error:", err).Error()
		}
		log.Info(pid, key, task, "SaveState")
		id, err := model.SaveState(pid, key, task.FlowId, task.Id, 0, creator, date, make(map[string]string), 0)
		if err != nil {
			log.Error(err)
			return err
		}

		log.Info(">>>>>>>>>>", flow_inst)
		if flow_inst != nil && len(strings.Trim(flow_inst.StartupScript, " \r\n\t")) > 0 {
			log.Info("Run startup script:", flow_inst.StartupScript)
			err = L.DoString(flow_inst.StartupScript)
			if err != nil {
				output = output + log.Error("Run startup script error:", err).Error()
			}
		}

		err = L.DoString(*task.Script)
		output = output + L.GetOutput()
		log.Info("Script output:", output)
		if err != nil {
			log.Error("err: ", err , " task: ", task)
			model.SaveState(pid, key, task.FlowId, task.Id, 3, creator, date, make(map[string]string), 0)
			model.AddTaskInstScriptLog((int)(id), task.Id, output+"\nERROR:"+err.Error())
			return log.Error("Run script error:", err)
		} else if L.GetRemoteTaskCount() == 0 {
			model.SaveState(pid, key, task.FlowId, task.Id, 2, creator, date, make(map[string]string), 0)
			model.UpdateTaskInstRemoteExec(flow_inst.Id, task.Id, "", "")
			RunNextTask(pid, key, task.FlowId, task.Id, creator, date, flow_inst.Id)
		} else {
			remote_exec_record := L.RemoteExecRecords[len(L.RemoteExecRecords)-1]
			model.UpdateTaskInstRemoteExec(flow_inst.Id, task.Id, remote_exec_record.Host, remote_exec_record.Uuid)
		}

		err = model.AddTaskInstScriptLog((int)(id), task.Id, output)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func RerunTask(task_inst_id int, task_id string, single_task bool, creator string) error {
	flow_inst, err := model.GetFlowInstById(task_inst_id)
	if err != nil {
		log.Error("err: ", err , " task: ", task_id)
		return log.Error("Find flow instance failed:", err)
	}
	if flow_inst == nil {
		return log.Error("Flow instance not exists.")
	}

	task_inst, err := model.GetTaskInstById(task_inst_id, task_id)
	if err != nil {
		log.Error("err: ", err , " task: ", task_id)
		return log.Error("Find task instance failed:", err)
	}
	if task_inst == nil {
		return log.Error("Task instance not exists.")
	}

	task, err := model.GetTaskById(task_inst.FlowId, task_id)
	if err != nil {
		log.Error("err: ", err , " task: ", task_id)
		return log.Error("Find next state failed:", err)
	}
	if task == nil {
		return log.Error("Task not exists.")
	}

	flow_inst.Creator = creator
	flow_inst.BeginTask = task.Id
	flow_inst.EndTask = ""
	if single_task {
		err = model.ResetTaskInstState(task_inst_id, task.Id)
		flow_inst.EndTask = task.Id
	} else {
		err = model.ResetTaskInstStateSinceOrderId(task_inst_id, task.FlowId, task.OrderId)
	}
	if err != nil {
		log.Error("err: ", err , " task: ", task_id)
		return log.Error("Update task state failed:", err)
	}

	_, err = model.AddFlowInst(flow_inst)
	if err != nil {
		log.Error("err: ", err , " task: ", task_id)
		return log.Error("Add flow instance failed:", err)
	}

	helper.Notify("flow.schedule", struct{}{})

	return nil
}

func RunNextTask(pid string, key string, flow_id string, task_id string, creator string, date time.Time, flow_inst_id int) {
	var flow_inst *model.FlowInst

	next_task, err := model.GetNextTask(flow_id, task_id)
	var next_task_inst *model.TaskInst

	if err != nil {
		log.Error("Find next state failed:", err)
		goto FAILED
	}

	if next_task == nil {
		goto SUCCEED
	}

	flow_inst, err = model.GetFlowInstById(flow_inst_id)
	if err != nil {
		log.Error("GetFlowInstById failed:", err)
		goto FAILED
	}

	log.Info("TaskInstId:", flow_inst_id, " TaskId:", next_task.Id)
	next_task_inst, err = model.GetTaskInstById(flow_inst_id, next_task.Id)

	if err != nil {
		log.Error("Find next state failed:", err)
		goto FAILED
	}

	log.Infof("TaskInst:%+v", next_task_inst)

	if flow_inst.EndTask == task_id {
		goto SUCCEED
	}

	if next_task_inst != nil && next_task_inst.State != nil && *next_task_inst.State >= 0 {
		goto SUCCEED
	}

	err = StartFlowTask(pid, key, next_task, creator, date, nil, next_task_inst)
	if err != nil {
		alarm.RaiseAlarm("FLOW_STATE", "Start next task failed.", err.Error(),
			&map[string]interface{}{
				"next_task": next_task,
			})
		goto FAILED
	}
	return

SUCCEED:
	model.UpdateFlowInstState(pid, key, flow_id, date, model.StateSucceed)
	helper.Notify("flow.schedule", struct{}{})
	return

FAILED:
	model.UpdateFlowInstState(pid, key, flow_id, date, model.StateFailed)
	return
}

func (s *iState) Lua_FindTaskByFlowId(flow_id string) (tasks []*model.Task) {
	tasks, err := model.FindTaskByFlowId(flow_id)
	if err != nil {
		log.Error("err: ", err , " flow_id: ", flow_id)
		s.RaiseError(err.Error())
	}
	return tasks
}

func (s *iState) Lua_GetTaskInstById(flow_inst_id int, task_id string) (task_instance *model.TaskInst) {
	task_instance, err := model.GetTaskInstById(flow_inst_id, task_id)
	if err != nil {
		log.Error("err: ", err , " task_id: ", task_id)
		s.RaiseError(err.Error())
	}
	return task_instance
}
