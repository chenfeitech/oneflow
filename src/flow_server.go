package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"flag"
	"time"

	"config"
	"scheduler"
	"web_portal/server"

	//"github.com/gorilla/mux"
	log "github.com/cihub/seelog"
)

var (
	logtoconsole = flag.Bool("logtoconsole", false, "Log ouptut to console.")
)

func main() {
	//debug.SetTraceback("crash")
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// Initialize log
	var logconf []byte
	if *logtoconsole {
		logconf = []byte(logtoconsoleconf)
	} else {
		logconf = []byte(logtofileconf)
	}
	logger, _ := log.LoggerFromConfigAsBytes(logconf)
	log.ReplaceLogger(logger)
	defer log.Flush()

	log.Info("Args:", os.Args)
	log.Info("Envs:", os.Environ())

	flow_job_scheduler := scheduler.New(
		"FlowJobScheduler", // name,
		"tbFlowConf",       // tableScheduler,
		"`id`",             // columnId,
		"concat(pid, '_', process_id)", // columnJobName,
		"`start_time`",                     // columnPattern,
		`concat('StartFlow("', process_id,'", "', pid,'", "", time.Now().Add(' ,data_delay, '*-24*time.Hour), "idata", nil)')`, // columnScript,
		`last_run_time`,                      // columnLastRunTime,
		`next_run_time`,                      // columnNextRunTime,
		"last_result",                        // columnLastRunResult,
		"last_error",                         // columnLastRunError,
		"isactive=1 AND now() > active_date", // conditionEnabled,
		"tbFlowSchdRunLog",               // tableLog
	)
	go flow_job_scheduler.RunLoop()
	go scheduler.ScheduleLoop()

	http.Handle("/", server.Router)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *config.ServerPort), nil)
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
	}
}

const (
	logtofileconf = `
	<seelog>
		<outputs>
			<rollingfile formatid="log" type="size" filename="../log/data_flow.log" maxsize="100000000" maxrolls="5"/>
		</outputs>
		<formats>
		    <format id="log" format="%Date %Time [%Level] %File %Line %Func %Msg%n"/>
		</formats>
	</seelog>
	`
	logtoconsoleconf = `
	<seelog>
		<outputs>
			<console formatid="out"/>
		</outputs>
		<formats>
		    <format id="out" format="%Time [%Level] %File %Line %Func %Msg%n"/>
		</formats>
	</seelog>
	`
)
