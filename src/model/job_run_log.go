package model

import (
	"config"
	"time"
)

type JobRunLog struct {
	Id           int
	JobId        int
	Output       string
	Errors       string
	ScheduleTime time.Time
	BeginTime    time.Time
	EndTime      time.Time
	Result       int
}

func AddJobRunLog(model *JobRunLog) (int64, error) {
	sqlStr := "INSERT INTO `tbJobRunLog` (`job_id`, `output`, `errors`, `schedule_time`, `begin_time`, `end_time`, `result`) VALUES(?,?,?,?,?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.JobId, model.Output, model.Errors, model.ScheduleTime, model.BeginTime, model.EndTime, model.Result)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindJobRunLog(condition string, args ...interface{}) ([]*JobRunLog, error) {
	sqlStr := "SELECT `id`, `job_id`, `output`, ifnull(`errors`, ''), `schedule_time`, `begin_time`, `end_time`, `result` FROM `tbJobRunLog`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*JobRunLog, 0)

	rows, err := config.GetDBConnect().Query(sqlStr, args...)
	if err != nil {
		return results, err
	} else {
		defer rows.Close()
		for rows.Next() {
			model := JobRunLog{}
			rows.Scan(&model.Id, &model.JobId, &model.Output, &model.Errors, &model.ScheduleTime, &model.BeginTime, &model.EndTime, &model.Result)

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetJobRunLog(condition string, args ...interface{}) (*JobRunLog, error) {
	results, err := FindJobRunLog(condition, args...)

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
