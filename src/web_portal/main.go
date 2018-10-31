package main

import (
	//"github.com/gorilla/mux"
	"config"
	"flag"
	"flow"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"schedule_job"
	"time"
	"web_portal/server"

	log "github.com/cihub/seelog"
)

var (
	logtoconsole = flag.Bool("logtoconsole", false, "Log ouptut to console.")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	//This won't work
	//http.Handle("/static/", http.FileServer(http.Dir("./static/")))

	//These will work
	//first alternative
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//second alternative
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

	go flow.RunLoop()
	go schedule_job.RunLoop()

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
		    <format id="log" format="%Date %Time [%Level] %File %Func %Msg%n"/>
		</formats>
	</seelog>
	`
	logtoconsoleconf = `
	<seelog>
		<outputs>
			<console formatid="out"/>
		</outputs>
		<formats>
		    <format id="out" format="[%Level] %File %Func %Msg%n"/>
		</formats>
	</seelog>
	`
)
