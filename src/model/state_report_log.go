package model

import (
	"config"
	"time"
)

type StateReportLog struct {
	Id         int       `sql:"id"`
	PId        string    `sql:"pid"`
	Key        string    `sql:"key"`
	FlowId     string    `sql:"flow_id"`
	TaskId     string    `sql:"task_id"`
	State      int       `sql:"state"`
	Creator    string    `sql:"creator"`
	Date       time.Time `sql:"date"`
	ReportTime time.Time `sql:"report_time"`
	ExtraData  string    `sql:"extra_data"`
	FlowInstId int
}

func AddStateReportLog(model *StateReportLog) (int64, error) {
	sqlStr := "INSERT INTO `tbStateReportLog` (`pid`, `key`, `flow_id`, `task_id`, `state`, `creator`, `date`, `extra_data`) VALUES(?,?,?,?,?,?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.PId, model.Key, model.FlowId, model.TaskId, model.State, model.Creator, model.Date, model.ExtraData)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindStateReportLog(condition string, args ...interface{}) ([]*StateReportLog, error) {
	sqlStr := "SELECT `id`, `pid`, `key`, `flow_id`, `task_id`, `state`, `creator`, `date`, `report_time`, `extra_data` FROM `tbStateReportLog`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*StateReportLog, 0)

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
			model := StateReportLog{}
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
			}
			rows.Scan(values...)
			model.Id = (int)((*(values[0].(*interface{}))).(int64))
			model.PId = (string)((*(values[1].(*interface{}))).([]uint8))
			model.Key = (string)((*(values[2].(*interface{}))).([]uint8))
			model.FlowId = (string)((*(values[3].(*interface{}))).([]uint8))
			model.TaskId = (string)((*(values[4].(*interface{}))).([]uint8))
			model.State = (int)((*(values[5].(*interface{}))).(int64))
			model.Creator = (string)((*(values[6].(*interface{}))).([]uint8))
			model.Date = (*(values[7].(*interface{}))).(time.Time)
			model.ReportTime = (*(values[8].(*interface{}))).(time.Time)
			model.ExtraData = (string)((*(values[9].(*interface{}))).([]uint8))

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetStateReportLog(condition string, args ...interface{}) (*StateReportLog, error) {
	results, err := FindStateReportLog(condition, args...)

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
