package scheduler

import (
	"fmt"
	"sync"
	"time"

	"lua_helper"
	"model"

	log "github.com/cihub/seelog"
)

var (
	wg       = sync.WaitGroup{}
	job_chan = make(chan *model.JobSchedule)
)

func RunLoop() {
	next_run_time := time.Now()
	next_run_time = next_run_time.Add(time.Duration(-next_run_time.Second())*time.Second + time.Duration(-next_run_time.Nanosecond())*time.Nanosecond)

	for i := 0; i < 8; i++ {
		go jobRunner()
	}
	for {
		if time.Now().After(next_run_time) {
			log.Info("Begin schedule")

			scheduleNewJobs(next_run_time)
			runJobs(next_run_time)

			next_run_time = next_run_time.Add(time.Minute)
			log.Info("End schedule")
		} else {
			time.Sleep(time.Second)
		}
	}
}

// 处理新添加的计划任务
func scheduleNewJobs(t time.Time) {
	new_jobs, err := model.FindNewJobSchedule()

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
		if err := model.SetJobScheduleNextRunTime(job.Id, last_run_time); err != nil {
			log.Error("Schedule job ", job.Id, ":", job.JobName, " run at ", last_run_time, " error:", err)
		} else {
			log.Debug("Schedule job ", job.Id, ":", job.JobName, " run at ", last_run_time)
		}
	}
}

// 加载计划任务插入运行Channel
func runJobs(t time.Time) {
	to_run_jobs, err := model.FindJobScheduleToRun(t)

	if err != nil {
		log.Error("FindJobScheduleToRun error:", err)
		return
	}

	for _, job := range to_run_jobs {
		wg.Add(1)
		job_chan <- job
	}
	wg.Wait()
}

func jobRunner() {
	for {
		job, flag := <-job_chan
		if !flag {
			log.Error("flag: ", flag)
			return
		}
		func(j *model.JobSchedule) {
			defer wg.Done()
			begin_time := time.Now()
			err := runJob(j, begin_time)
			if err != nil {
				log.Error("Run schedule job ", job.Id, ":", job.JobName, " failed:", err)
			}
			rescheduleJob(j, begin_time)
		}(job)
	}
}

func runJob(job *model.JobSchedule, begin_time time.Time) (err error) {
	log.Info("Begin run schedule job ", job.Id, ":", job.JobName)
	defer log.Info("End run schedule job ", job.Id, ":", job.JobName)

	L := lua_helper.GetState()
	defer lua_helper.RevokeState(L)
	err = L.DoString(job.Script)

	end_time := time.Now()

	job_log := model.JobRunLog{}
	job_log.JobId = job.Id
	job_log.BeginTime = begin_time
	job_log.EndTime = end_time
	job_log.ScheduleTime = *job.NextRunTime

	if err != nil {
		log.Error("run schedule job ", job.Id, ":", job.JobName, " err: ", err)
		job_log.Result = -1
		job_log.Errors = err.Error()
	} else {
		job_log.Result = 0
	}

	if dberr := model.SetJobScheduleLastRunTime(job.Id, begin_time); dberr != nil {
		log.Error("Save schedule job ", job.Id, ":", job.JobName, " last run time error:", err)
		if err == nil {
			err = dberr
		}
	}

	if _, dberr := model.AddJobRunLog(&job_log); dberr != nil {
		log.Error("Save schedule job ", job.Id, ":", job.JobName, " run log failed:", err, fmt.Sprintf("\n%+v", job_log))
		if err == nil {
			err = dberr
		}
	}
	return err
}

func rescheduleJob(job *model.JobSchedule, begin_time time.Time) error {
	actual, err := Parse(job.Pattern)
	if err != nil {
		log.Error("Parse job ", job.Id, ":", job.JobName, " pattern error:", err)
		return err
	}

	next_run_time := actual.Next(begin_time)

	if err := model.SetJobScheduleNextRunTime(job.Id, next_run_time); err != nil {
		log.Error("Reschedule job ", job.Id, ":", job.JobName, " run after ", next_run_time, " error:", err)
		return err
	} else {
		log.Debug("Reschedule job ", job.Id, ":", job.JobName, " run after ", next_run_time)
	}
	return nil
}
