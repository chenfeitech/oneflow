package model

import (
	"config"

	"strings"
	"time"

	log "github.com/cihub/seelog"
)

type Flow struct {
	Id            string     `sql:"id"`
	Name          string     `sql:"name"`
	Description   string     `sql:"description"`
	CreateTime    time.Time  `sql:"create_time"`
	Creator       string     `sql:"creator"`
	StartTimer    string     `sql:"start_timer"`
	NextRunTime   *time.Time `sql:"next_run_time"`
	StartupScript string     `sql:"startup_script"`
	LastRunTime   *time.Time `sql:"last_run_time"`
	LastRunLog    string     `sql:"last_run_log"`
}

func FindFlow(condition string, args ...interface{}) ([]*Flow, error) {
	sqlStr := "SELECT `id`, `name`, `description`, `create_time`, `creator`, `start_timer`, `next_run_time`, `startup_script`, `last_run_time`, `last_run_log` FROM `tbFlow`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*Flow, 0)

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
			model := Flow{}
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
			model.Id = (string)((*(values[0].(*interface{}))).([]uint8))
			model.Name = (string)((*(values[1].(*interface{}))).([]uint8))
			model.Description = (string)((*(values[2].(*interface{}))).([]uint8))
			model.CreateTime = (*(values[3].(*interface{}))).(time.Time)
			model.Creator = (string)((*(values[4].(*interface{}))).([]uint8))
			model.StartTimer = (string)((*(values[5].(*interface{}))).([]uint8))
			if *(values[6].(*interface{})) != nil {
				t_NextRunTime := (*(values[6].(*interface{}))).(time.Time)
				model.NextRunTime = &t_NextRunTime
			}
			model.StartupScript = (string)((*(values[7].(*interface{}))).([]uint8))
			if *(values[8].(*interface{})) != nil {
				LastRunTime := (*(values[8].(*interface{}))).(time.Time)
				model.LastRunTime = &LastRunTime
			}
			model.LastRunLog = (string)((*(values[9].(*interface{}))).([]uint8))
			results = append(results, &model)
		}
	}
	return results, nil
}

func GetFlow(condition string, args ...interface{}) (*Flow, error) {
	results, err := FindFlow(condition, args...)

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

func GetFlowByKey(id string) (*Flow, error) {
	return GetFlow("`id`=?", id)
}

func AddFlow(flow *Flow, tasks []*Task) error {
	tx, err := config.GetDBConnect().Begin()
	if err != nil {
		log.Error("id: ", flow.Id, " err: ", err)
		return err
	}
	sqlStr := "INSERT INTO `tbFlow` (`id`, `name`, `description`, `creator`, `start_timer`, `startup_script`, last_run_log) VALUES(?,?,?,?,'',?,'')"
	_, err = tx.Exec(sqlStr, flow.Id, flow.Name, flow.Description, flow.Creator, flow.StartupScript)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, task := range tasks {
		sqlStr = "INSERT INTO `tbTask` (`id`, `flow_id`, `name`, `max_retries`, `description`, `order_id`, `parent_id`, `script`) VALUES(?,?,?,?,?,?,?,?)"
		_, err = tx.Exec(sqlStr, task.Id, flow.Id, task.Name, task.MaxRetries, task.Description, task.OrderId, task.ParentId, task.Script)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func UpdateFlow(flow *Flow, tasks []*Task, deleteIds []string) error {
	tx, err := config.GetDBConnect().Begin()
	if err != nil {
		log.Error("id: ", flow.Id, " err: ", err)
		return err
	}
	_, err = tx.Exec("UPDATE `tbFlow` SET `name`=?, `description`=?, `startup_script`=? WHERE `id`=?", flow.Name, flow.Description, flow.StartupScript, flow.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(deleteIds) > 0 {
		ids := make([]string, len(deleteIds), len(deleteIds))
		args := make([]interface{}, len(deleteIds)+1, len(deleteIds)+1)
		args[0] = flow.Id
		for i := 0; i < len(deleteIds); i++ {
			ids[i] = "?"
			args[i+1] = deleteIds[i]
		}

		_, err = tx.Exec("DELETE FROM `tbTask` WHERE `flow_id`=? AND `id` in ("+strings.Join(ids, ",")+")", args...)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	insert_stmt, err := tx.Prepare("REPLACE INTO `tbTask` (`id`, `flow_id`, `name`, `max_retries`, `description`, `order_id`, `parent_id`, `script`) VALUES(?,?,?,?,?,?,?,?)")
	defer insert_stmt.Close()
	if err != nil {
		tx.Rollback()
		return err
	}
	for i, task := range tasks {
		_, err := insert_stmt.Exec(task.Id, flow.Id, task.Name, task.MaxRetries, task.Description, i, task.ParentId, task.Script)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func FindFlowToSchedule() ([]*Flow, error) {
	return FindFlow("`start_timer` <> '' AND `startup_script` <> '' AND  `next_run_time` IS NULL")
}

func SetFlowNextRunTime(id string, t time.Time) error {
	sqlStr := "UPDATE `tbFlow` SET `next_run_time`=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, id)
	return err
}

func UpdateFlowLastRunLog(id string, t time.Time, log string) error {
	sqlStr := "UPDATE `tbFlow` SET `last_run_time`=?, `last_run_log`=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, log, id)
	return err
}

func FindFlowToRun(t time.Time) ([]*Flow, error) {
	return FindFlow("`next_run_time` <= ?", t)
}
