package model

import (
	"config"
	"fmt"
	"time"
)

type TaskInst struct {
	FlowInstId     int        `sql:"flow_inst_id"`
	Id             string     `sql:"id"`
	FlowId         string     `sql:"flow_id"`
	Name           string     `sql:"name"`
	Description    string     `sql:"description"`
	OrderId        int        `sql:"order_id"`
	ParentId       int        `sql:"parent_id"`
	State          *int       `sql:"state"`
	ReadyTime      *time.Time `sql:"ready_time"`
	RunningTime    *time.Time `sql:"running_time"`
	SucceedTime    *time.Time `sql:"succeed_time"`
	FailedTime     *time.Time `sql:"failed_time"`
	LastUpdateTime *time.Time `sql:"last_update_time"`
	ScriptOutput   *string    `sql:"script_output"`
	Retries        int        `sql:"retries"`
	RemoteExecHost string     `sql:"remote_exec_host"`
	RemoteExecUuid string     `sql:"remote_exec_uuid"`
}

func FindTaskInst(condition string, args ...interface{}) ([]*TaskInst, error) {
	sqlStr := "SELECT t.`id`, t.`flow_id`, t.`name`, t.`description`, t.`order_id`, t.`parent_id`, ti.`state`, ti.ready_time, ti.running_time, ti.succeed_time, ti.failed_time, ti.last_update_time, ti.script_output, ti.retries, fi.id, ti.remote_exec_host, ti.remote_exec_uuid FROM `tbTask` t " +
		" inner join `tbFlowInst` fi ON t.`flow_id`=fi.flow_id AND fi.id=? left join `tbTaskInst` ti ON t.id = ti.`task_id` AND fi.id = ti.flow_inst_id"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}

	results := make([]*TaskInst, 0)

	stmt, err := config.GetDBConnect().Prepare(sqlStr)
	if err != nil {
		return results, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return results, err
	} else {
		defer rows.Close()
		for rows.Next() {
			model := TaskInst{}
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
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)
			model.Id = (string)((*(values[0].(*interface{}))).([]uint8))
			model.FlowId = (string)((*(values[1].(*interface{}))).([]uint8))
			model.Name = (string)((*(values[2].(*interface{}))).([]uint8))
			model.Description = (string)((*(values[3].(*interface{}))).([]uint8))
			model.OrderId = (int)((*(values[4].(*interface{}))).(int64))
			model.ParentId = (int)((*(values[5].(*interface{}))).(int64))
			if *(values[6].(*interface{})) == nil {
				model.State = nil
			} else {
				t_State := (int)((*(values[6].(*interface{}))).(int64))
				model.State = &t_State
			}
			if *(values[7].(*interface{})) == nil {
				model.ReadyTime = nil
			} else {
				t_ReadyTime := (*(values[7].(*interface{}))).(time.Time)
				model.ReadyTime = &t_ReadyTime
			}
			if *(values[8].(*interface{})) == nil {
				model.RunningTime = nil
			} else {
				t_RunningTime := (*(values[8].(*interface{}))).(time.Time)
				model.RunningTime = &t_RunningTime
			}
			if *(values[9].(*interface{})) == nil {
				model.SucceedTime = nil
			} else {
				t_SucceedTime := (*(values[9].(*interface{}))).(time.Time)
				model.SucceedTime = &t_SucceedTime
			}
			if *(values[10].(*interface{})) == nil {
				model.FailedTime = nil
			} else {
				t_FailedTime := (*(values[10].(*interface{}))).(time.Time)
				model.FailedTime = &t_FailedTime
			}
			if *(values[11].(*interface{})) == nil {
				model.LastUpdateTime = nil
			} else {
				t_LastUpdateTime := (*(values[11].(*interface{}))).(time.Time)
				model.LastUpdateTime = &t_LastUpdateTime
			}
			if *(values[12].(*interface{})) == nil {
				model.ScriptOutput = nil
			} else {
				t_ScriptOutput := (string)((*(values[12].(*interface{}))).([]uint8))
				model.ScriptOutput = &t_ScriptOutput
			}
			if *(values[13].(*interface{})) == nil {
			} else {
				model.Retries = (int)((*(values[13].(*interface{}))).(int64))
			}
			model.FlowInstId = (int)((*(values[14].(*interface{}))).(int64))
			if *(values[15].(*interface{})) == nil {
			} else {
				model.RemoteExecHost = (string)((*(values[15].(*interface{}))).([]uint8))
			}
			if *(values[16].(*interface{})) == nil {
			} else {
				model.RemoteExecUuid = (string)((*(values[16].(*interface{}))).([]uint8))
			}
			results = append(results, &model)
		}
	}
	return results, nil
}

func GetTaskInst(condition string, args ...interface{}) (*TaskInst, error) {
	results, err := FindTaskInst(condition, args...)

	if err != nil {
		return nil, err
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return nil, nil
		}
	}
}

func FindTaskInstByFlow(flow_inst_id int, flow_id string) ([]*TaskInst, error) {
	return FindTaskInst("t.`flow_id`=? ORDER BY t.order_id", flow_inst_id, flow_id)
}

func GetTaskInstById(flow_inst_id int, task_id string) (*TaskInst, error) {
	return GetTaskInst("t.`id`=? ORDER BY t.order_id", flow_inst_id, task_id)
}

func AddTaskInstScriptLog(flow_inst_id int, taskId string, script_log string) error {
	sqlStr := "UPDATE `tbTaskInst` SET `script_output`=concat(ifnull(`script_output`,''),?)  WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, script_log, flow_inst_id, taskId)
	return err
}

func UpdateTaskInstState(flow_inst_id int, taskId string, state int) error {
	sqlStr := "UPDATE `tbTaskInst` SET `State`=?  WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, state, flow_inst_id, taskId)
	return err
}

func UpdateTaskInstRetries(flow_inst_id int, taskId string, retries int) error {
	sqlStr := "UPDATE `tbTaskInst` SET `Retries`=?  WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, retries, flow_inst_id, taskId)
	return err
}

func UpdateTaskInstStateSinceOrderId(flow_inst_id int, flow_id string, order_id int, state int) error {
	sqlStr := "UPDATE `tbTaskInst` ti SET `State`=? WHERE `flow_inst_id`=? AND `task_id` IN (SELECT Id FROM `tbTask` WHERE `flow_id`=? and `order_id` >= ?)"
	_, err := config.GetDBConnect().Exec(sqlStr, state, flow_inst_id, flow_id, order_id)
	return err
}

func ResetTaskInstState(flow_inst_id int, taskId string) error {
	sqlStr := "UPDATE `tbTaskInst` SET `State`=-1, `retries`=0, remote_exec_host='', remote_exec_uuid='' WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, flow_inst_id, taskId)
	return err
}

func ResetTaskInstStateSinceOrderId(flow_inst_id int, flow_id string, order_id int) error {
	sqlStr := "UPDATE `tbTaskInst` ti SET `State`=-1, `retries`=0, remote_exec_host='', remote_exec_uuid=''  WHERE `flow_inst_id`=? AND `task_id` IN (SELECT Id FROM `tbTask` WHERE `flow_id`=? and `order_id` >= ?)"
	_, err := config.GetDBConnect().Exec(sqlStr, flow_inst_id, flow_id, order_id)
	fmt.Println(sqlStr, flow_inst_id, flow_id, order_id)
	return err
}
func IncreaseTaskInstRetries(flow_inst_id int, taskId string) error {
	sqlStr := "UPDATE `tbTaskInst` SET `retries`=`retries`+1 WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, flow_inst_id, taskId)
	return err
}

func UpdateTaskInstRemoteExec(flow_inst_id int, taskId string, host string, uuid string) error {
	sqlStr := "UPDATE `tbTaskInst` SET remote_exec_host=?, remote_exec_uuid=? WHERE `flow_inst_id`=? AND `task_id`=?"
	_, err := config.GetDBConnect().Exec(sqlStr, host, uuid, flow_inst_id, taskId)
	return err
}
