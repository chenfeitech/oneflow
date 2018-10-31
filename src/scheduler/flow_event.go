package scheduler

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"config"
	"utils/alarm"
	"lua_helper"
	"model"

	log "github.com/cihub/seelog"
)

var (
	init_script = flag.String("init_script", func() string {
		if runtime.GOOS == "darwin" {
			return *config.ServerRoot + "/scripts/init.lua"
		} else {
			return *config.ServerRoot + "/scripts/init.lua"
		}
	}(), "Flow event lua init script.")
)

var (
	g_mutex     = &sync.Mutex{}
	pid_mutexs = make(map[string]*sync.Mutex, 0)
)

func OnStateChange(state_log *model.StateReportLog) {
	log.Info("State change event process.")
	g_mutex.Lock()
	pid_mutex := pid_mutexs[state_log.PId]
	if pid_mutex == nil {
		pid_mutex = &sync.Mutex{}
		pid_mutexs[state_log.PId] = pid_mutex
	}
	g_mutex.Unlock()

	pid_mutex.Lock()
	defer pid_mutex.Unlock()

	switch state_log.State {
	case model.StateSucceed:
		runNextTask(state_log)
	case model.StateFailed:
		retryTask(state_log)
	}

}

// type FlontFlowStateLog struct {
// 	FlowId        string `json:"FlowId"`
// 	PId           string `json:"PId"`
// 	FlowTaskState int    `json:"FlowTaskState"`
// 	RunningTime   string `json:"RunningTime"`
// }

func runNextTask(state_log *model.StateReportLog) {
	lua_helper.RunNextTask(state_log.PId, state_log.Key, state_log.FlowId, state_log.TaskId, state_log.Creator, state_log.Date, state_log.FlowInstId)
}

func retryTask(state_log *model.StateReportLog) {
	task, err := model.GetTaskById(state_log.FlowId, state_log.TaskId)
	if err != nil || task == nil {
		log.Error("Get Task failed:", err)
		return
	}

	task_inst, err := model.GetTaskInstById(state_log.FlowInstId, state_log.TaskId)
	if err != nil || task_inst == nil {
		log.Error("Get Task instance failed:", err)
		return
	}

	if task.MaxRetries <= task_inst.Retries {
		log.Info(state_log.PId, state_log.Key, state_log.FlowId, state_log.Date, model.StateFailed)
		log.Info("------", model.UpdateFlowInstState(state_log.PId, state_log.Key, state_log.FlowId, state_log.Date, model.StateFailed))
		log.Info(state_log.FlowInstId, "-", state_log.TaskId, " failed exceed retries.")
		alarm.RaiseAlarm("FLOW_STATE", "Task execute failed", fmt.Sprintf("PId:%v date:%v flow:%v task:%v", state_log.PId, state_log.Date.Format("2006-01-02"), task.FlowId, task.Name),
			&map[string]interface{}{
				"pid":      state_log.PId,
				"date":      state_log.Date,
				"FlowId":    state_log.FlowId,
				"TaskId":    state_log.TaskId,
				"state_log": state_log,
				"task":      task,
			})
		return
	}
	time.Sleep(10*time.Second + time.Duration(task_inst.Retries)*time.Minute)
	err = model.IncreaseTaskInstRetries(state_log.FlowInstId, state_log.TaskId)
	if err != nil {
		log.Error("Get Task instance failed:", err)
		return
	}

	err = lua_helper.StartFlowTask(state_log.PId, state_log.Key, task, state_log.Creator, state_log.Date, nil, task_inst)
	if err != nil {
		log.Info("------", model.UpdateFlowInstState(state_log.PId, state_log.Key, state_log.FlowId, state_log.Date, model.StateFailed))
		alarm.RaiseAlarm("FLOW_STATE", "Retry task failed.", err.Error(),
			&map[string]interface{}{
				"state_log": state_log,
				"task":      task,
			})
	}
}
