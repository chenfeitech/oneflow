package model

import (
	"config"
	"time"
)

type TaskInstAlarm struct {
	Id         int       `sql:"id"`
	FlowInstId int       `sql:"flow_inst_id"`
	TaskId     string    `sql:"task_id"`
	Content    string    `sql:"content"`
	Timestamp  time.Time `sql:"timestamp"`
}

func AddTaskInstAlarm(model *TaskInstAlarm) (int64, error) {
	sqlStr := "INSERT INTO `tbTaskInstAlarm` (`flow_inst_id`, `task_id`, `content`) VALUES(?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.FlowInstId, model.TaskId, model.Content)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindTaskInstAlarm(condition string, args ...interface{}) ([]*TaskInstAlarm, error) {
	sqlStr := "SELECT `id`, `flow_inst_id`, `task_id`, `content`, `timestamp` FROM `tbTaskInstAlarm`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*TaskInstAlarm, 0)

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
			model := TaskInstAlarm{}
			values := []interface{}{
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)
			model.Id = (int)((*(values[0].(*interface{}))).(int64))
			model.FlowInstId = (int)((*(values[1].(*interface{}))).(int64))
			model.TaskId = (string)((*(values[2].(*interface{}))).([]uint8))
			model.Content = (string)((*(values[3].(*interface{}))).([]uint8))
			model.Timestamp = (*(values[4].(*interface{}))).(time.Time)

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetTaskInstAlarm(condition string, args ...interface{}) (*TaskInstAlarm, error) {
	results, err := FindTaskInstAlarm(condition, args...)

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
