package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/rpc/json"
	"net/http"
	"os"
	"model/rpc_data"
	"strings"
	"time"
)

var (
	pid         = flag.String("pid", "", "product id.")
	flow_id     = flag.String("flow_id", "", "flow id.")
	task_id     = flag.String("task_id", "", "task id.")
	date        = flag.String("date", "", "task state date. format: yyyy-mm-dd")
	state       = flag.Int("state", 0, "task state: 1=Ready, 2=Running, 3=Succeed, 4=Failed")
	reportor    = flag.String("reportor", "", "task state reportor")
	service_url = flag.String("service_url", "http://localhost:3001/api", "flow service api url.")
)

func main() {
	flag.Parse()

	args := rpc_data.StateDataArgs{}
	args.PId = *pid
	args.FlowId = *flow_id
	args.TaskId = *task_id
	args.State = *state
	args.Creator = *reportor
	args.ExtraData = make(map[string]string)

	if len(*pid) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [PId].")
		os.Exit(1)
	}
	if len(*flow_id) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [flow_id].")
		os.Exit(1)
	}
	if len(*task_id) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [task_id].")
		os.Exit(1)
	}
	if *state <= 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Argument [state] must gra.")
		os.Exit(1)
	}
	// d, err := time.Parse("2006-01-02", *date)

	for _, arg := range flag.Args() {
		pair := strings.SplitN(arg, "=", 2)
		if len(pair) == 2 {
			args.ExtraData[pair[0]] = pair[1]
		}
	}

	request, err := json.EncodeClientRequest("FlowService.ReportState", args)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(request))

	resp, err := http.Post("http://localhost:3001/api", "application/json; charset=utf-8", bytes.NewReader(request))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
	if resp.StatusCode != 200 {
		fmt.Println("ERROR")
		return
	}

	reply := rpc_data.StateReply{}
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply)
}
