package scheduler

import (
	"config"

	"time"

    log "github.com/cihub/seelog"
)

type jobSchedule struct {
	Id          int    `sql:"id"`
	JobName     string `sql:"job_name"`
	Pattern     string `sql:"pattern"`
	Script      string `sql:"script"`
	NextRunTime *time.Time
}

func (s *scheduler) findJobSchedule(condition string, args ...interface{}) ([]*jobSchedule, error) {
	sqlStr := "SELECT " + s.columnId + ", " + s.columnJobName + ", " + s.columnPattern + ", " + s.columnScript + ", " + s.columnNextRunTime + " FROM " + s.tableScheduler +
		" WHERE " + s.conditionEnable

	if len(condition) > 0 {
		sqlStr = sqlStr + " AND " + condition
	}
    log.Info("sql: ", sqlStr, " args: ", args)

	results := make([]*jobSchedule, 0)

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
			model := jobSchedule{}
			var nextRunTime interface{} = new(interface{})
			rows.Scan(&model.Id, &model.JobName, &model.Pattern, &model.Script, nextRunTime)
			if *(nextRunTime.(*interface{})) == nil {
				model.NextRunTime = nil
			} else {
				t_NextRunTime := (*(nextRunTime.(*interface{}))).(time.Time)
				model.NextRunTime = &t_NextRunTime
			}
			results = append(results, &model)
		}
	}
	return results, nil
}

func (s *scheduler) SetLastRunLog(id int, t time.Time, result int, errors string) error {
	sqlStr := "UPDATE " + s.tableScheduler + " SET " + s.columnLastRunTime + "=?, " + s.columnLastRunResult + "=?, " + s.columnLastRunError + "=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, result, errors, id)
	return err
}

func (s *scheduler) SetNextRunTime(id int, t time.Time) error {
	sqlStr := "UPDATE " + s.tableScheduler + " SET " + s.columnNextRunTime + "=? WHERE `id`=?"

	_, err := config.GetDBConnect().Exec(sqlStr, t, id)
	return err
}

func (s *scheduler) FindNew() ([]*jobSchedule, error) {
	return s.findJobSchedule(s.columnNextRunTime + " IS NULL")
}

func (s *scheduler) FindToRun(t time.Time) ([]*jobSchedule, error) {
	return s.findJobSchedule(s.columnNextRunTime+" <= ?", t)
}

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

func (s *scheduler) AddRunLog(model *JobRunLog) (int64, error) {
	sqlStr := "INSERT INTO " + s.tableLog + " (`job_id`, `output`, `errors`, `schedule_time`, `begin_time`, `end_time`, `result`) VALUES(?,?,?,?,?,?,?)"
	result, err := config.GetDBConnect().Exec(sqlStr, model.JobId, model.Output, model.Errors, model.ScheduleTime, model.BeginTime, model.EndTime, model.Result)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}
