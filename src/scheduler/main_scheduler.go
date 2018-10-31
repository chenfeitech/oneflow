package scheduler

import (
	"fmt"
	"sync"
	"time"

	"lua_helper"

	log "github.com/cihub/seelog"
)

type scheduler struct {
	name string

	tableScheduler      string // Scheduler 配置表名
	columnId            string // ScheduleID 列
	columnJobName       string // Job名称列
	columnPattern       string // 时间模式列
	columnScript        string // 执行脚本列
	columnLastRunTime   string // 最后一次执行时间列
	columnNextRunTime   string // 下一次执行时间列
	columnLastRunResult string // 最后一次执行结果列
	columnLastRunError  string // 最后一次执行错误信息列
	conditionEnable     string // 是否启用条件
	tableLog            string // 日志表

	waitGroup sync.WaitGroup
	jobChan   chan *jobSchedule
}

func New(name,
	tableScheduler,
	columnId,
	columnJobName,
	columnPattern,
	columnScript,
	columnLastRunTime,
	columnNextRunTime,
	columnLastRunResult,
	columnLastRunError,
	conditionEnable,
	tableLog string) *scheduler {
	return &scheduler{
		name,
		tableScheduler,
		columnId,
		columnJobName,
		columnPattern,
		columnScript,
		columnLastRunTime,
		columnNextRunTime,
		columnLastRunResult,
		columnLastRunError,
		conditionEnable,
		tableLog,

		sync.WaitGroup{},
		make(chan *jobSchedule),
	}
}

func (s *scheduler) RunLoop() {
	next_run_time := time.Now()
	next_run_time = next_run_time.Add(time.Duration(-next_run_time.Second())*time.Second + time.Duration(-next_run_time.Nanosecond())*time.Nanosecond)

	for i := 0; i < 8; i++ {
		go s.jobRunner()
	}
	for {
		if time.Now().After(next_run_time) {
			log.Info("Begin schedule ", s.name)

			s.scheduleNewJobs(next_run_time)
			s.runJobs(next_run_time)

			next_run_time = next_run_time.Add(time.Minute)
			log.Info("End schedule ", s.name)
		} else {
			time.Sleep(time.Second)
		}
	}
}

// 处理新添加的计划任务
func (s *scheduler) scheduleNewJobs(t time.Time) {
	new_jobs, err := s.FindNew()

	if err != nil {
		log.Error("FindNewJobSchedule error:", err)
		return
	}

	for _, job := range new_jobs {
		actual, err := Parse(job.Pattern)
		if err != nil {
			log.Error("Parse job ", job.Id, ":", job.JobName, " pattern error:", err)
			continue
		}

		// 计算计划任务执行时间
		last_run_time := t
		if !actual.IsTimeMatches(t) {
			last_run_time = actual.Next(t)
		}
		if err := s.SetNextRunTime(job.Id, last_run_time); err != nil {
			log.Error("Schedule job ", job.Id, ":", job.JobName, " run at ", last_run_time, " error:", err)
		} else {
			log.Debug("Schedule job ", job.Id, ":", job.JobName, " run at ", last_run_time)
		}
	}
}

// 加载计划任务插入运行Channel
func (s *scheduler) runJobs(t time.Time) {
	to_run_jobs, err := s.FindToRun(t)

	if err != nil {
		log.Error("FindJobScheduleToRun error:", err)
		return
	}

	for _, job := range to_run_jobs {
		s.waitGroup.Add(1)
		s.jobChan <- job
	}
	s.waitGroup.Wait()
}

func (s *scheduler) jobRunner() {
	for {
		job, flag := <-s.jobChan
		if !flag {
			log.Error("flag: ", flag)
			return
		}
		func(j *jobSchedule) {
			defer s.waitGroup.Done()
			begin_time := time.Now()
			err := s.runJob(j, begin_time)
			if err != nil {
				log.Error("Run schedule job ", job.Id, ":", job.JobName, " failed:", err)
			}
			s.rescheduleJob(j, begin_time)
		}(job)
	}
}

func (s *scheduler) runJob(job *jobSchedule, begin_time time.Time) (err error) {
	log.Info("Begin run schedule job ", job.Id, ":", job.JobName)
	defer log.Info("End run schedule job ", job.Id, ":", job.JobName)

	L := lua_helper.GetState()
	defer lua_helper.RevokeState(L)
	err = L.DoString(job.Script)

	end_time := time.Now()

	job_log := JobRunLog{}
	job_log.JobId = job.Id
	job_log.BeginTime = begin_time
	job_log.EndTime = end_time
	job_log.ScheduleTime = *job.NextRunTime

	if err != nil {
		log.Info("run schedule job ", job.Id, ":", job.JobName, " err: ", err)
		job_log.Result = -1
		job_log.Errors = err.Error()
	} else {
		job_log.Result = 0
	}

	if dberr := s.SetLastRunLog(job.Id, begin_time, job_log.Result, job_log.Errors); dberr != nil {
		log.Error("Save schedule job ", job.Id, ":", job.JobName, " last run time error:", err)
		if err == nil {
			err = dberr
		}
	}

	if _, dberr := s.AddRunLog(&job_log); dberr != nil {
		log.Error("Save schedule job ", job.Id, ":", job.JobName, " run log failed:", err, fmt.Sprintf("\n%+v", job_log))
		if err == nil {
			err = dberr
		}
	}
	return err
}

func (s *scheduler) rescheduleJob(job *jobSchedule, begin_time time.Time) error {
	actual, err := Parse(job.Pattern)
	if err != nil {
		log.Error("Parse job ", job.Id, ":", job.JobName, " pattern error:", err)
		return err
	}

	next_run_time := actual.Next(begin_time)

	if err := s.SetNextRunTime(job.Id, next_run_time); err != nil {
		log.Error("Reschedule job ", job.Id, ":", job.JobName, " run after ", next_run_time, " error:", err)
		return err
	} else {
		log.Debug("Reschedule job ", job.Id, ":", job.JobName, " run after ", next_run_time)
	}
	return nil
}
