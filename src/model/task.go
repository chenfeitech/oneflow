package model

import (
	"config"
)

type Task struct {
	Id          string  `sql:"id"`
	FlowId      string  `sql:"flow_id"`
	Name        string  `sql:"name"`
	Description string  `sql:"description"`
	OrderId     int     `sql:"order_id"`
	ParentId    int     `sql:"parent_id"`
	Script      *string `sql:"script"`
	MaxRetries  int     `sql:"max_retries"`
}

func AddTask(model *Task) (int64, error) {
	sqlStr := "INSERT INTO `tbTask` (`id`, `flow_id`, `name`, `description`, `order_id`, `parent_id`, `script`, `max_retries`) VALUES(?,?,?,?,?,?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.Id, model.FlowId, model.Name, model.Description, model.OrderId, model.ParentId, model.Script, model.MaxRetries)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindTask(condition string, args ...interface{}) ([]*Task, error) {
	sqlStr := "SELECT `id`, `flow_id`, `name`, `description`, `order_id`, `parent_id`, `script`, `max_retries` FROM `tbTask`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*Task, 0)

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
			model := Task{}
			values := []interface{}{
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
				model.Script = nil
			} else {
				t_Script := (string)((*(values[6].(*interface{}))).([]uint8))
				model.Script = &t_Script
			}
			model.MaxRetries = (int)((*(values[7].(*interface{}))).(int64))

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetTask(condition string, args ...interface{}) (*Task, error) {
	results, err := FindTask(condition, args...)

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

func FindTaskByFlowId(flow_id string) ([]*Task, error) {
	return FindTask("`flow_id`=?  ORDER BY `order_id`", flow_id)
}

func GetTaskById(flow_id string, task_id string) (*Task, error) {
	return GetTask("`flow_id`=? AND `id`=?",
		flow_id, task_id)
}

func GetNextTask(flow_id string, task_id string) (*Task, error) {
	return GetTask("`flow_id`=? AND `order_id`>(SELECT `order_id` FROM `tbTask` t WHERE t.`flow_id`=? AND `id`=?) ORDER BY `order_id` limit 1",
		flow_id, flow_id, task_id)
}
