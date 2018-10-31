package model

import (
	"config"
	"fmt"
	"time"
)

type JobSchedule struct {
	Id          int        `sql:"id"`
	JobName     string     `sql:"job_name"`
	PId         string     `sql:"pid"`
	Pattern     string     `sql:"pattern"`
	Script      string     `sql:"script"`
	LastRunTime *time.Time `sql:"last_run_time"`
	NextRunTime *time.Time `sql:"next_run_time"`
	Creator     string     `sql:"creator"`
	CreateAt    time.Time  `sql:"create_at"`
	Enabled     int        `sql:"enabled"`
	LastResult  int        `sql:"last_result"`
	LastError   *string    `sql:"last_error"`
}

func AddJobSchedule(model *JobSchedule) (int64, error) {
	sqlStr := "INSERT INTO `tbJobSchedule` (`job_name`, `pattern`, `script`, `last_run_time`, `next_run_time`, `creator`, `enabled`, `last_result`, `last_error`) VALUES(?,?,?,?,?,?,?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.JobName, model.Pattern, model.Script, model.LastRunTime, model.NextRunTime, model.Creator, model.Enabled, model.LastResult, model.LastError)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindJobSchedule(condition string, args ...interface{}) ([]*JobSchedule, error) {
	sqlStr := "SELECT `id`, `job_name`, `pid`, `pattern`, `script`, `last_run_time`, `next_run_time`, `creator`, `create_at`, `enabled`, `last_result`, `last_error` FROM `tbJobSchedule`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*JobSchedule, 0)

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
			model := JobSchedule{}
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
			}
			rows.Scan(values...)
			model.Id = (int)((*(values[0].(*interface{}))).(int64))
			model.JobName = (string)((*(values[1].(*interface{}))).([]uint8))
			model.PId = (string)((*(values[2].(*interface{}))).([]uint8))
			model.Pattern = (string)((*(values[3].(*interface{}))).([]uint8))
			model.Script = (string)((*(values[4].(*interface{}))).([]uint8))
			if *(values[5].(*interface{})) == nil {
				model.LastRunTime = nil
			} else {
				t_LastRunTime := (*(values[5].(*interface{}))).(time.Time)
				model.LastRunTime = &t_LastRunTime
			}
			if *(values[6].(*interface{})) == nil {
				model.NextRunTime = nil
			} else {
				t_NextRunTime := (*(values[6].(*interface{}))).(time.Time)
				model.NextRunTime = &t_NextRunTime
			}
			model.Creator = (string)((*(values[7].(*interface{}))).([]uint8))
			model.CreateAt = (*(values[8].(*interface{}))).(time.Time)
			model.Enabled = (int)((*(values[9].(*interface{}))).(int64))
			model.LastResult = (int)((*(values[10].(*interface{}))).(int64))
			if *(values[11].(*interface{})) == nil {
				model.LastError = nil
			} else {
				t_LastError := (string)((*(values[11].(*interface{}))).([]uint8))
				model.LastError = &t_LastError
			}

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetJobSchedule(condition string, args ...interface{}) (*JobSchedule, error) {
	results, err := FindJobSchedule(condition, args...)

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

func GetJobScheduleByKey(id int) (*JobSchedule, error) {
	return GetJobSchedule("`Id`=?", id)
}

func SetJobScheduleLastRunTime(id int, t time.Time) error {
	sqlStr := "UPDATE `tbJobSchedule` SET `last_run_time`=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, id)
	return err
}

func SetJobScheduleNextRunTime(id int, t time.Time) error {
	sqlStr := "UPDATE `tbJobSchedule` SET `next_run_time`=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, id)
	return err
}

func FindNewJobSchedule() ([]*JobSchedule, error) {
	return FindJobSchedule("`next_run_time` IS NULL AND `enabled`=1")
}

func FindJobScheduleToRun(t time.Time) ([]*JobSchedule, error) {
	return FindJobSchedule("`next_run_time` <= ? AND `enabled`=1", t)
}

func FindJobScheduleByPage(page, page_size int, condition string, args ...interface{}) ([]*JobSchedule, error) {
	return FindJobSchedule(condition+fmt.Sprint(" 1=1 LIMIT ", (page-1)*page_size, ", ", page_size), args...)
}
