package subcmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"config"
	"middleware/jsonrpc_client"
	"middleware/nsqhelper"
	"model/rpc_data"

	"github.com/codegangsta/cli"
)

var lastReport = time.Now()

func init() {
	Register(cli.Command{
		Name:    "report",
		Aliases: []string{"r"},
		Usage:   "Report job run status",
		Action:  report,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "pid, g",
				Usage:  "pid id",
				EnvVar: "PID",
			},
			cli.StringFlag{
				Name:   "flow_id, f",
				Usage:  "flow id",
				EnvVar: "FLOW_ID",
			},
			cli.StringFlag{
				Name:   "task_id, t",
				Usage:  "task id",
				EnvVar: "TASK_ID",
			},
			cli.StringFlag{
				Name:   "key, k",
				Usage:  "flow instance key",
				EnvVar: "KEY",
			},
			cli.StringFlag{
				Name:   "date, d",
				Usage:  "task state date. format: yyyy-mm-dd",
				EnvVar: "DATE",
			},
			cli.IntFlag{
				Name:  "state, s",
				Usage: "task state: 1=Ready, 2=Running, 3=Succeed, 4=Failed",
			},
			cli.StringSliceFlag{
				Name:  "extra_data, e",
				Usage: "task extra data. format: key=value",
			},
			cli.StringFlag{
				Name:  "reportor, r",
				Usage: "task state reportor",
			},
			cli.IntFlag{
				Name:  "retries",
				Usage: "task state reportor",
				Value: -1,
			},
		},
	})
}

func report(c *cli.Context) {
	pid := c.String("pid")
	flow_id := c.String("flow_id")
	task_id := c.String("task_id")
	key := c.String("key")
	date := c.String("date")
	state := c.Int("state")
	reportor := c.String("reportor")

	retries := c.Int("retries")

	if len(pid) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [pid].")
		os.Exit(1)
	}
	if len(flow_id) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [flow_id].")
		os.Exit(1)
	}
	if len(task_id) == 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Requires argument [task_id].")
		os.Exit(1)
	}
	if state <= 0 {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Argument [state] must be positive numbers.")
		os.Exit(1)
	}
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Fprintln(os.Stderr, time.Now(), " [Error] Argument [date] parse failed:", err)
		os.Exit(1)
	}

	extra_data := make(map[string]string)
	for _, pair := range c.StringSlice("extra_data") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			fmt.Fprintln(os.Stderr, time.Now(), " [Error] Argument [extra_data] illegal:", pair)
			os.Exit(1)
		}
		extra_data[kv[0]] = kv[1]
	}

	err = ReportState(pid, flow_id, task_id, key, date, reportor, state, extra_data, retries)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func ReportState(pid, flow_id, task_id, key, date, reportor string, state int, extra_data map[string]string, retries int) error {

	request_args := rpc_data.StateDataArgs{}
	request_args.PId = pid
	request_args.FlowId = flow_id
	request_args.TaskId = task_id
	request_args.Key = key
	request_args.State = state
	request_args.Creator = reportor
	request_args.Date = date
	request_args.ExtraData = extra_data
	fmt.Println("ReportState", request_args, GetServerUrl("/api"))
	defer fmt.Println("ReportState", request_args, GetServerUrl("/api"), " Finished")

	request, err := jsonrpc_client.EncodeClientRequest("FlowService.ReportState", request_args)
	if err != nil {
		return fmt.Errorf(" [Error] EncodeClientRequest failed:%v", err)
	}

	if retries < 0 {
		retries = 100
	}

	for i := 0; i <= retries; i++ {
		if i > 0 {
			time.Sleep(5 * time.Second)
		}

		os.Stdout.Sync()
		err = nsqhelper.PublishMessage(config.NSQFlowTaskStatusTopic, request)
		if err != nil {
			fmt.Fprintln(os.Stderr, time.Now(), " [Error] Publish message failed:", err)
			err = fmt.Errorf("[Error] Publish message: %v", err)
			continue
		}
		os.Stdout.Sync()
		break

		// req, err := http.NewRequest("POST", GetServerUrl("/api"), bytes.NewReader(request))
		// req.Header.Set("Content-Type", "application/json")
		// client := &http.Client{}
		// client.Timeout = time.Second * 5
		// resp, err := client.Do(req)

		// fmt.Println("Post fin")

		// os.Stdout.Sync()
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, time.Now(), " [Error] Post failed:", err)
		// 	err = fmt.Errorf("[Error] Post failed: %v", err)
		// 	continue
		// }
		// if resp.Body != nil {
		// 	defer resp.Body.Close()
		// }

		// if resp.StatusCode != 200 {
		// 	fmt.Fprintln(os.Stderr, time.Now(), " [Error] Post failed StatusCode:", resp.StatusCode)
		// 	err = fmt.Errorf("[Error] Post failed StatusCode: %v", resp.StatusCode)
		// 	continue
		// }

		// reply := rpc_data.StateReply{}
		// err = json.DecodeClientResponse(resp.Body, &reply)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, time.Now(), " [Error] DecodeClientResponse:", err)
		// 	err = fmt.Errorf("[Error] DecodeClientResponse: %v", resp.StatusCode)
		// 	continue
		// }

		// fmt.Println(reply)
		// err = nil
	}
	return err
}
