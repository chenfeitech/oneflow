package scheduler

import (
	"sync"
	"time"

	"utils/helper"
	"lua_helper"
	"model"

	log "github.com/cihub/seelog"
)

var (
	schedule_mutex sync.Mutex
)

func ScheduleLoop() {
	schedule_mutex.Lock()
	defer schedule_mutex.Unlock()
	schedule_chan := helper.CreateNotificationPort("flow.schedule", 1)

	log.Info("Flow scheduler started.")

	for {
		select {
		case <-schedule_chan:
		case <-time.After(5 * time.Second):
		}
		log.Info("Flow scheduler run.")

		flow_insts, err := model.GetFlowInstToRun()
		if err != nil {
			log.Error("GetFlowInstToRun failed:", err)
			continue
		}

		for _, flow_inst := range flow_insts {
			log.Info(flow_inst)
			flow_inst := flow_inst
			model.UpdateFlowInstState(flow_inst.PId, flow_inst.Key, flow_inst.FlowId, flow_inst.RunningDay, model.StateRunning)
			go func() {
				err := lua_helper.StartFlowInst(flow_inst)
				if err != nil {
					log.Error("Flow instance start failed:", err)
				}
			}()
		}
	}
}
