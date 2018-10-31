package subcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"model/rpc_data"
	"time"

	"github.com/gorilla/rpc/json"
)

const userHZ = 100

var (
	server_host   = os.Getenv("SERVER_HOST")
	server_port   = os.Getenv("SERVER_PORT")
	server_prefix = func() string {
		if p := os.Getenv("SERVER_PREFIX"); len(p) != 0 {
			return p
		} else {
			return "/data_flow"
		}
	}()
)

func GetServerUrl(path string) string {
	url := "http://" + server_host
	if len(server_port) != 0 {
		url = url + ":" + server_port
	}

	if len(server_prefix) != 0 {
		if (server_prefix)[0] != '/' {
			url = url + "/" + server_prefix
		} else {
			url = url + server_prefix
		}
	}

	if len(path) > 0 {
		if url[len(url)-1] != '/' {
			url = url + "/"
		}
		if path[0] != '/' {
			url = url + path
		} else {
			url = url + path[1:]
		}
	}
	return url
}

type ProcStat struct {
	// The process ID.
	PID int
	// The filename of the executable.
	Comm string
	// The process state.
	State string
	// The PID of the parent of this process.
	PPID int
	// The process group ID of the process.
	PGRP int
	// The session ID of the process.
	Session int
	// The controlling terminal of the process.
	TTY int
	// The ID of the foreground process group of the controlling terminal of
	// the process.
	TPGID int
	// The kernel flags word of the process.
	Flags uint
	// The number of minor faults the process has made which have not required
	// loading a memory page from disk.
	MinFlt uint
	// The number of minor faults that the process's waited-for children have
	// made.
	CMinFlt uint
	// The number of major faults the process has made which have required
	// loading a memory page from disk.
	MajFlt uint
	// The number of major faults that the process's waited-for children have
	// made.
	CMajFlt uint
	// Amount of time that this process has been scheduled in user mode,
	// measured in clock ticks.
	UTime uint
	// Amount of time that this process has been scheduled in kernel mode,
	// measured in clock ticks.
	STime uint
	// Amount of time that this process's waited-for children have been
	// scheduled in user mode, measured in clock ticks.
	CUTime uint
	// Amount of time that this process's waited-for children have been
	// scheduled in kernel mode, measured in clock ticks.
	CSTime uint
	// For processes running a real-time scheduling policy, this is the negated
	// scheduling priority, minus one.
	Priority int
	// The nice value, a value in the range 19 (low priority) to -20 (high
	// priority).
	Nice int
	// Number of threads in this process.
	NumThreads int
	// The time the process started after system boot, the value is expressed
	// in clock ticks.
	Starttime uint64
	// Virtual memory size in bytes.
	VSize int
	// Resident set size in pages.
	RSS int
}

// NewStat returns the current status information of the process.
func GetProcStat(pid int) (ProcStat, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return ProcStat{}, err
	}

	var (
		ignore int

		s = ProcStat{PID: pid}
		l = bytes.Index(data, []byte("("))
		r = bytes.LastIndex(data, []byte(")"))
	)

	if l < 0 || r < 0 {
		return ProcStat{}, fmt.Errorf(
			"unexpected format, couldn't extract comm: %s",
			data,
		)
	}

	s.Comm = string(data[l+1 : r])
	_, err = fmt.Fscan(
		bytes.NewBuffer(data[r+2:]),
		&s.State,
		&s.PPID,
		&s.PGRP,
		&s.Session,
		&s.TTY,
		&s.TPGID,
		&s.Flags,
		&s.MinFlt,
		&s.CMinFlt,
		&s.MajFlt,
		&s.CMajFlt,
		&s.UTime,
		&s.STime,
		&s.CUTime,
		&s.CSTime,
		&s.Priority,
		&s.Nice,
		&s.NumThreads,
		&ignore,
		&s.Starttime,
		&s.VSize,
		&s.RSS,
	)
	if err != nil {
		return ProcStat{}, err
	}

	return s, nil
}

func (s ProcStat) VirtualMemory() int {
	return s.VSize
}

// StartTime returns the unix timestamp of the process in seconds.
func (s ProcStat) StartTime() time.Time {
	return time.Unix(BootTime()+(int64(s.Starttime)/userHZ), 0)
}

// CPUTime returns the total CPU user and system time in seconds.
func (s ProcStat) CPUTime() float64 {
	return float64(s.UTime+s.STime) / userHZ
}

var (
	bootTime int64 = 0
)

func BootTime() int64 {
	if bootTime > 0 {
		return bootTime
	}
	inFile, _ := os.Open("/proc/stat")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		n, _ := fmt.Sscanf(scanner.Text(), "btime %d", &bootTime)
		if n == 1 {
			return bootTime
		}
	}
	return 0
}

func SendAlarm(flow_inst_id int, task_id, content string, retries int) error {
	request_args := rpc_data.TaskAlarmArgs{}
	request_args.FlowInstId = flow_inst_id
	request_args.TaskId = task_id
	request_args.Content = content

	fmt.Println("TaskAlarm", request_args, GetServerUrl("/api"))
	defer fmt.Println("TaskAlarm", request_args, GetServerUrl("/api"), " Finished")

	request, err := json.EncodeClientRequest("FlowService.TaskAlarm", request_args)
	if err != nil {
		return fmt.Errorf(" [Error] EncodeClientRequest failed:", err)
	}

	if retries < 0 {
		retries = 10
	}

	for i := 0; i <= retries; i++ {
		if i > 0 {
			time.Sleep(5)
		}

		os.Stdout.Sync()
		req, err := http.NewRequest("POST", GetServerUrl("/api"), bytes.NewReader(request))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		client.Timeout = time.Second * 5
		resp, err := client.Do(req)

		fmt.Println("Post fin")

		os.Stdout.Sync()
		if err != nil {
			fmt.Fprintln(os.Stderr, time.Now(), " [Error] Post failed:", err)
			err = fmt.Errorf("[Error] Post failed: %v", err)
			continue
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}

		if resp.StatusCode != 200 {
			fmt.Fprintln(os.Stderr, time.Now(), " [Error] Post failed StatusCode:", resp.StatusCode)
			err = fmt.Errorf("[Error] Post failed StatusCode: %v", resp.StatusCode)
			continue
		}

		err = nil
		break
	}
	return err
}
