package model

import (
	"config"
	"utils/helper"

	"encoding/json"
	"time"

	log "github.com/cihub/seelog"
)

var (
	StateTimeField = map[int]string{
		StateReady:   "ready_time",
		StateRunning: "running_time",
		StateSucceed: "succeed_time",
		StateFailed:  "failed_time",
	}
)

func SaveState(pid string, key string, flowId string, taskId string, state int,
	createor string, runningDay time.Time, extraData map[string]string, timestamp int64) (int64, error) {
	defer func() {
		go func() {
			j, _ := json.Marshal(map[string]interface{}{"pid": pid, "key": key, "flow_id": flowId, "task_id": taskId, "running_day": runningDay.Format("2006-01-02"), "state": state})
			helper.WsHub.Broadcast(j)
		}()
	}()

	tx, err := config.GetDBConnect().Begin()
	if err != nil {
		return 0, err
	}
	flow_state := state
	if state == StateReady {
		flow_state = StateRunning
	}
	if state == StateSucceed {
		flow_state = StateRunning
	}
	sqlStr := "INSERT INTO `tbFlowInst` (`flow_id`, `pid`, `key`, `running_day`, `creator`, `last_task_id`, `last_task_state`, `last_update_time`, `state`, `startup_script`) VALUES(?,?,?,?,?,?,?, now(), ?, '') ON DUPLICATE KEY UPDATE `last_task_id`=?, `last_task_state`=?, `last_update_time`=now(), `state`=?"
	result, err := tx.Exec(sqlStr, flowId, pid, key, runningDay, createor, taskId, state, flow_state, taskId, state, flow_state)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()

	if id == 0 || err != nil {
		err = tx.QueryRow("SELECT `id` FROM `tbFlowInst` WHERE `flow_id`=? AND `pid`=? AND `key`=? AND `running_day`=date(?)",
			flowId, pid, key, runningDay).Scan(&id)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	state_time_field := StateTimeField[state]

	if len(state_time_field) > 0 {
		sqlStr = "INSERT INTO `tbTaskInst` (`flow_inst_id`, `task_id`, `state`, `last_update_time`, " +
			state_time_field + ", `status_timestamp`) VALUES(?,?,?,now(),now(),?) " +
			" ON DUPLICATE KEY UPDATE `state`=if(`status_timestamp`>?, `state`, ?), `last_update_time`=now(), " +
			state_time_field + "=now(), `status_timestamp`=if(`status_timestamp`>?, `status_timestamp`, ?)"
	} else {
		sqlStr = "INSERT INTO `tbTaskInst` (`flow_inst_id`, `task_id`, `state`,`last_update_time`,`status_timestamp`)" +
			" VALUES(?,?,?,now(),?) ON DUPLICATE KEY UPDATE `state`=if(`status_timestamp`>?, `state`, ?),`last_update_time`=now()," +
			" `status_timestamp`=if(`status_timestamp`>?, `status_timestamp`, ?)"
	}

	_, err = tx.Exec(sqlStr, id, taskId, state, timestamp, timestamp, state, timestamp, timestamp)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

type FlowInst struct {
	Id             int       `sql:"id"`
	FlowId         string    `sql:"flow_id"`
	Name           string    `sql:"name"`
	Description    string    `sql:"description"`
	PId            string    `sql:"pid"`
	Key            string    `sql:"key"`
	RunningDay     time.Time `sql:"running_day"`
	CreateTime     time.Time `sql:"create_time"`
	State          int       `sql:"state"`
	StartupScript  string    `sql:"startup_script"`
	LastTask       string    `sql:"last_task"`
	LastTaskState  int       `sql:"last_task_state"`
	LastUpdateTime time.Time `sql:"last_update_time"`
	Creator        string    `sql:"creator"`
	BeginTask      string    `sql:"begin_task"`
	EndTask        string    `sql:"end_task"`
}

func FindFlowInst(condition string, args ...interface{}) ([]*FlowInst, error) {
	sqlStr := "SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`, fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time`, ifnull(t.name, fi.`last_task_id`), fi.`last_task_state`, fi.last_update_time, fi.`state`, fi.`startup_script`, fi.`begin_task`, fi.`end_task` from `tbFlowInst` fi inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id  and t.`flow_id` = fi.`flow_id`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*FlowInst, 0)

	stmt, err := config.GetDBConnect().Prepare(sqlStr)
	if err != nil {
		log.Error("sql: ", sqlStr, " err: ", err)
		return results, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Error("sql: ", sqlStr, " err: ", err)
		return results, err
	} else {
		defer rows.Close()
		for rows.Next() {
			model := FlowInst{}
			values := []interface{}{
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)
			model.Id = (int)((*(values[0].(*interface{}))).(int64))
			model.FlowId = (string)((*(values[1].(*interface{}))).([]uint8))
			model.Name = (string)((*(values[2].(*interface{}))).([]uint8))
			model.Description = (string)((*(values[3].(*interface{}))).([]uint8))
			model.PId = (string)((*(values[4].(*interface{}))).([]uint8))
			model.Key = (string)((*(values[5].(*interface{}))).([]uint8))
			model.RunningDay = (*(values[6].(*interface{}))).(time.Time)
			model.CreateTime = (*(values[7].(*interface{}))).(time.Time)
			model.LastTask = (string)((*(values[8].(*interface{}))).([]uint8))
			model.LastTaskState = (int)((*(values[9].(*interface{}))).(int64))
			model.LastUpdateTime = (*(values[10].(*interface{}))).(time.Time)
			model.State = (int)((*(values[11].(*interface{}))).(int64))
			model.StartupScript = (string)((*(values[12].(*interface{}))).([]uint8))
			model.BeginTask = (string)((*(values[13].(*interface{}))).([]uint8))
			model.EndTask = (string)((*(values[14].(*interface{}))).([]uint8))

			results = append(results, &model)
		}
	}
	return results, nil
}

func FindFlowInstByPage(page, page_size int, condition string, args ...interface{}) ([]*FlowInst, error) {
	return FindFlowInst(condition, args...)
}

func GetFlowInst(condition string, args ...interface{}) (*FlowInst, error) {
	results, err := FindFlowInst(condition, args...)

	if err != nil {
		log.Error("condition: ", condition, " err: ", err)
		return nil, err
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return nil, nil
		}
	}
}

func AddFlowInst(inst *FlowInst) (int, error) {
	sqlStr := "INSERT INTO `tbFlowInst` (`flow_id`, `pid`, `key`, `running_day`, `creator`, `last_task_id`, `last_task_state`, `last_update_time`, `state`, `startup_script`, `begin_task`, `end_task`) VALUES(?,?,?,?,?,?,?, now(),?,?,?,?) ON DUPLICATE KEY UPDATE `state`=?, `startup_script`=?, `creator`=?, `begin_task`=?, `end_task`=?"
	result, err := config.GetDBConnect().Exec(sqlStr, inst.FlowId, inst.PId, inst.Key,
		inst.RunningDay, inst.Creator, "", StateReady, StateReady, inst.StartupScript,
		inst.BeginTask, inst.EndTask,
		StateReady, inst.StartupScript, inst.Creator,
		inst.BeginTask, inst.EndTask)
	if err != nil {
		log.Error("Pid: ", inst.PId, "flowId: ", inst.FlowId, " err: ", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if id == 0 || err != nil {
		err = config.GetDBConnect().QueryRow("SELECT `id` FROM `tbFlowInst` WHERE `flow_id`=? AND `pid`=? AND `key`=? AND `running_day`=date(?)",
			inst.FlowId, inst.PId, inst.Key, inst.RunningDay).Scan(&id)
		if err != nil {
			log.Error("Pid: ", inst.PId, "flowId: ", inst.FlowId, " err: ", err)
			return 0, err
		}
	}
	return int(id), nil
}

func GetFlowInstById(id int) (*FlowInst, error) {
	return GetFlowInst("fi.`id`=?", id)
}

func GetFlowInst4(pid string, key string, flowId string, runningDay time.Time) (*FlowInst, error) {
	return GetFlowInst("fi.`flow_id`=? AND fi.`pid`=? AND fi.`key`=? AND fi.`running_day`=date(?)", flowId, pid, key, runningDay)
}

func (fi *FlowInst) GetTasksInsts() []*TaskInst {
	tasks, _ := FindTaskInstByFlow(fi.Id, fi.FlowId)
	return tasks
}

func GetFlowInstState(pid string, key string, flowId string, runningDay time.Time) (state int, err error) {
	flow_inst, err := GetFlowInst("fi.`flow_id`=? AND fi.`pid`=? AND fi.`key`=? AND fi.`running_day`=?", flowId, pid, key, runningDay)
	if err != nil {
		log.Error("Pid: ", pid, "flowId: ", flowId, " err: ", err)
		return 0, err
	}

	if flow_inst != nil {
		state = flow_inst.State
	}
	return
}

func UpdateFlowInstState(pid string, key string, flowId string, runningDay time.Time, state int) (err error) {
	log.Info(pid, key, flowId, state)
	_, err = config.GetDBConnect().Exec("UPDATE tbFlowInst SET `state`=? WHERE `flow_id`=? AND `pid`=? AND `key`=? AND `running_day`=date(?)", state, flowId, pid, key, runningDay)
	if err != nil {
		log.Error("Pid: ", pid, "flowId: ", flowId, " err: ", err)
		return err
	}

	return
}

func UpdateFlowInstStateById(id int, state int) (err error) {
	_, err = config.GetDBConnect().Exec("UPDATE tbFlowInst SET `state`=? WHERE `id`=?", state, id)
	if err != nil {
		log.Error("id: ", id, " err: ", err)
		return err
	}

	return
}

func GetFlowInstToRun() ([]*FlowInst, error) {
	return FindFlowInst("state = 0 AND NOT EXISTS (SELECT * FROM tbFlowInst ff WHERE ff.`flow_id`=fi.`flow_id` and ff.`pid`=fi.`pid` and ff.`key`=fi.`key` and ((ff.`running_day`<fi.`running_day` AND ff.`state`=0) OR  ff.`state`=3 OR  ff.`state`=1))")
}

func SetFlowInstStateSucceed(id int) (err error) {
	_, err = config.GetDBConnect().Exec("UPDATE tbTaskInst SET `state`=? WHERE `flow_inst_id`=? AND `state`=?", StateSucceed, id, StateFailed)
	if err != nil {
		log.Error("id: ", id, " err: ", err)
		return err
	}

	_, err = config.GetDBConnect().Exec("UPDATE tbFlowInst SET `state`=? WHERE `id`=?", StateSucceed, id)
	if err != nil {
		log.Error("id: ", id, " err: ", err)
		return err
	}

	return
}

func GetFlowInstByName(flow_id string, pid string, running_day string) (*FlowInst, error) {
	return GetFlowInst("fi.`flow_id`=? AND fi.`pid`=? AND fi.`running_day`=?", flow_id, pid, running_day)
}
