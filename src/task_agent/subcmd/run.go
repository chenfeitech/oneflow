package subcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"config"
	"utils/helper"

	"github.com/codegangsta/cli"
)

var (
	flow_inst_id int
	pid	     string
	flow_id      string
	task_id      string
	key          string
	date         string
)

func init() {
	flags := []cli.Flag{
		cli.IntFlag{
			Name:   "flow_inst_id, i",
			Usage:  "flow instance id",
			EnvVar: "FLOW_INST_ID",
		},
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
		cli.StringFlag{
			Name:   "proc_path",
			Usage:  "task proc path",
			Value:  *config.ServerRoot + "/proc",
			EnvVar: "PROC_PATH",
		},
		cli.StringFlag{
			Name:   "reportor",
			Usage:  "task status reportor",
			Value:  helper.GetIPAddr(),
			EnvVar: "PROC_PATH",
		},
		cli.StringFlag{
			Name:  "uuid, u",
			Usage: "run uuid",
		},
		cli.StringFlag{
			Name:  "program, p",
			Usage: "program to run",
		},
		cli.StringSliceFlag{
			Name:  "args, a",
			Usage: "run arguments",
		},
		cli.IntFlag{
			Name:  "rpc_rfd, r",
			Usage: "Rpc read file descriptor",
		},
		cli.IntFlag{
			Name:  "rpc_wfd, w",
			Usage: "Rpc write file descriptor",
		},
		cli.BoolFlag{
			EnvVar: "MONITOR",
			Name:   "monitor, m",
			Usage:  "Monitor process",
		},
	}

	Register(cli.Command{
		Name:   "run",
		Usage:  "Run program on node",
		Action: run,

		Flags: flags,
	})
}

func run(c *cli.Context) {
	m := c.Bool("monitor")
	if m {
		monitor(c)
		return
	}
	uuid := c.String("uuid")
	proc_path := c.String("proc_path")
	program := c.String("program")
	args := c.StringSlice("args")
	flow_inst_id = c.Int("flow_inst_id")
	pid = c.String("pid")
	flow_id = c.String("flow_id")
	reportor := c.String("reportor")
	task_id = c.String("task_id")
	key = c.String("key")
	date = c.String("date")

	if len(uuid) < 10 {
		fmt.Fprintln(os.Stderr, "run: Argument 'uuid' length must be greater to 10.")
		os.Exit(1)
	}

	now := time.Now().Format("20060102")

	if len(program) == 0 {
		fmt.Fprintln(os.Stderr, "run: Argument provide program to run.")
		os.Exit(1)
	}

	//	for _, env := range os.Environ() {
	//fmt.Println(env)
	//	}

	killPreviousProcess(strings.Join(append([]string{program}, args...), " "), proc_path, now)

	log_path := proc_path + "/" + now + "/" + uuid

	if _, err := os.Stat(log_path + "/pid.log"); os.IsExist(err) {
		pid, err := ioutil.ReadFile(log_path + "/pid.log")
		if err != nil {
			fmt.Fprintln(os.Stderr, "run: Open pid file ", log_path+"/pid.log", " failed:", err)
		} else {
			fmt.Fprintln(os.Stderr, "run: Process is exists. Pid:", pid)
			os.Exit(1)
		}
	}

	require_files := strings.Split(os.Getenv("REQUIRE_FILES"), ":")
	for _, require_file := range require_files {
		err := RequireFile(require_file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "run: Require file ", require_file, ":", err, ".")
			os.Exit(1)
		}
	}

	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, "run: Program '", program, "' can not found.")
		os.Exit(1)
	}

	if err := os.MkdirAll(log_path, 0777); err != nil {
		fmt.Fprintln(os.Stderr, "run: Make dir for process info failed:", err, ".")
		os.Exit(1)
	}

	rs, ws, _ := os.Pipe()
	rc, wc, _ := os.Pipe()

	monitor_args := os.Args[1:]

	monitor_args = append(monitor_args, "-r=3", "-w=4", "-m=true")
	env := append(os.Environ(), "MONITOR=true")

	cmd := exec.Command(os.Args[0], monitor_args...)

	cmd.Env = env
	cmd.ExtraFiles = []*os.File{rc, ws}
	cmd.SysProcAttr = new(syscall.SysProcAttr)
	cmd.SysProcAttr.Setpgid = true

	report_ch := make(chan os.Signal, 1)
	signal.Notify(report_ch, syscall.SIGUSR2)

	bootstamp := new(Bootstamp)
	bootstamp.on_report_running = func(report_args *RunningArgs) int {
		if report_args.Pid != 0 {
			fmt.Fprintf(os.Stdout, "[%v] RUN:DATE[%s]:UUID[%s]:PID[%d]\n", helper.GetIPAddr(), now, uuid, report_args.Pid)
		}
		if len(pid) > 0 && len(flow_id) > 0 && len(task_id) > 0 && len(date) > 0 {
			extra_data := make(map[string]string)
			extra_data["__PID"] = fmt.Sprint(report_args.Pid)
			ReportState(pid, flow_id, task_id, key, date, reportor, report_args.State, extra_data, -1)
		}
		rs.Close()
		return os.Getpid()
	}
	rpc.Register(bootstamp)

	conn_s := &rwRpcConn{
		Reader: rs,
		Writer: wc,
		closec: make(chan bool, 1),
	}

	go rpc.ServeConn(conn_s)
	time.Sleep(10 * time.Microsecond)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Start monitor failed:", err)
		os.Exit(1)
	}

	monitor_stop_ch := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		monitor_stop_ch <- err
	}()

	select {
	case _ = <-report_ch:
		if !bootstamp.running_reported {
			fmt.Fprintln(os.Stderr, "Monitor not report status.")
			os.Exit(1)
		}
		return

	case err := <-monitor_stop_ch:
		fmt.Fprintln(os.Stderr, "Monitor except quit:", err)
		os.Exit(1)

	case <-time.After(10 * time.Second):
		fmt.Fprintln(os.Stderr, "Wait process status timeout.")
		os.Exit(1)
	}
}

func monitor(c *cli.Context) {
	uuid := c.String("uuid")
	proc_path := c.String("proc_path")
	program := c.String("program")
	args := c.StringSlice("args")

	flow_inst_id = c.Int("flow_inst_id")
	pid = c.String("pid")
	flow_id = c.String("flow_id")
	task_id = c.String("task_id")
	date = c.String("date")

	now := time.Now().Format("20060102")
	log_path := proc_path + "/" + now + "/" + uuid

	rpc_rfd := c.Int("rpc_rfd")
	rpc_wfd := c.Int("rpc_wfd")

	rc := os.NewFile(uintptr(rpc_rfd), "")
	ws := os.NewFile(uintptr(rpc_wfd), "")

	exe_path, err := exec.LookPath(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, "run: Program '", program, "' can not found.")
		os.Exit(1)
	}

	date := c.String("date")
	flow_id := c.String("flow_id")
	task_id := c.String("task_id")
	reportor := c.String("reportor")
	key := c.String("key")
	extra_data := make(map[string]string)

	log, alarm, err := exec_program(exe_path, program, args, log_path, func(pid int) {
		conn_c := &rwRpcConn{
			Reader: rc,
			Writer: ws,
			closec: make(chan bool, 1),
		}
		client := rpc.NewClient(conn_c)
		args := RunningArgs{}
		args.Pid = pid
		args.State = 1
		var reply int

		r := lastReport.Add(1 * time.Second)
		now := time.Now()
		if r.After(now) {
			time.Sleep(r.Sub(now))
		}
		lastReport = time.Now()

		err := client.Call("Bootstamp.ReportRunning", &args, &reply)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ReportRunning error:", err)
			os.Exit(1)
		}
		client.Close()
		parent, err := os.FindProcess(reply)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Not found parent process.")
			return
		}
		if parent != nil {
			parent.Signal(syscall.SIGUSR2)
		}
	})

	extra_data["__LOG"] = log
	extra_data["__ALARM"] = alarm

	state := 2
	if err != nil {
		state = 3
		extra_data["__LOG"] = log + fmt.Sprintln("\n[ERROR]", err)
	}

	{
		r := lastReport.Add(1 * time.Second)
		now := time.Now()
		if r.After(now) {
			time.Sleep(r.Sub(now))
		}
		lastReport = time.Now()
	}

	if len(pid) > 0 && len(flow_id) > 0 && len(task_id) > 0 {
		err := ReportState(pid, flow_id, task_id, key, date, reportor, state, extra_data, -1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Report running failed:", err)
		}
	}
}

func exec_program(exe_path string, program string, args []string, log_path string, running_handler func(pid int)) (log string, alarm string, err error) {
	/* Change the file mode mask */
	_ = syscall.Umask(0022)

	os.Chdir("/")

	if f, e := os.OpenFile("/dev/null", os.O_RDWR, 0); e == nil {
		syscall.Dup2(int(f.Fd()), int(os.Stdin.Fd()))
	} else {
		return "", "", fmt.Errorf("Error: Redirect stdin failed: %v", e)
	}

	if f, e := os.OpenFile(log_path+"/out.log", os.O_RDWR|os.O_CREATE, 0666); e == nil {
		syscall.Dup2(int(f.Fd()), int(os.Stdout.Fd()))
	} else {
		return "", "", fmt.Errorf("Error: Redirect stdout failed: %v", e)
	}

	if f, e := os.OpenFile(log_path+"/err.log", os.O_RDWR|os.O_CREATE, 0666); e == nil {
		syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	} else {
		return "", "", fmt.Errorf("Error: Redirect stderr failed: %v", e)
	}

	if f, e := os.OpenFile(log_path+"/cmd", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); e == nil {
		fmt.Fprintln(f, strings.Join(append([]string{program}, args...), " "))
		f.Close()
	}

	if f, e := os.OpenFile(log_path+"/agent.pid", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); e == nil {
		fmt.Fprintln(f, os.Getpid())
		f.Close()
	}

	// // create a new SID for the child process
	// s_ret, s_errno := syscall.Setsid()
	// if s_errno != nil {
	// 	return "", "", fmt.Errorf("Error: syscall.Setsid errno: %v", s_errno)
	// }
	// if s_ret < 0 {
	// 	return "", "", fmt.Errorf("Error: syscall.Setsid failed: %v", s_ret)
	// }

	os.Chdir(path.Dir(exe_path))
	env := os.Environ()

	//	for _, env := range os.Environ() {
	//fmt.Println(env)
	//	}

	os.Stdout.Sync()

	log_rd, log_wr, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("Error: Create log pipe failed: %v", err)
	}
	defer log_rd.Close()
	defer log_wr.Close()

	alarm_rd, alarm_wr, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("Error: Create alarm pipe failed: %v", err)
	}
	defer alarm_rd.Close()
	defer alarm_wr.Close()

RUN:
	cmd := exec.Command(exe_path, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env = append(env, "__LOG_FD=3", "__ALARM_FD=4")
	cmd.Env = env

	cmd.ExtraFiles = []*os.File{log_wr, alarm_wr}

	execErr := cmd.Start()
	if execErr != nil {
		if execErr == syscall.ETXTBSY {
			time.Sleep(1 * time.Second)
			goto RUN
		}
		return "", "", fmt.Errorf("Error: Execute command failed: %v", execErr)
	}

	if f, e := os.OpenFile(log_path+"/pid", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); e == nil {
		fmt.Fprintln(f, cmd.Process.Pid)
		f.Close()
	}

	log_ch := make(chan string, 1)
	go readPipe(log_rd, log_ch)

	alarm_ch := make(chan string, 1)
	go processAlarmPipe(alarm_rd, alarm_ch)

	//running_report_ch := make(chan interface{}, 1)
	//go running_handler(running_report_ch)
	running_handler(cmd.Process.Pid)

	fmt.Println("Begin Wait")
	exec_err := cmd.Wait()
	fmt.Println("End Wait")

	alarm_wr.Close()
	log_wr.Close()

	fmt.Println("Get alarm_content")
	var alarm_content string
	select {
	case alarm_content = <-alarm_ch:
	case <-time.After(5 * time.Microsecond):
		fmt.Fprintln(os.Stderr, "Error: Read alarm pipe timedout")
	}

	fmt.Println("Get log")
	var log_content string
	select {
	case log_content = <-log_ch:
	case <-time.After(5 * time.Microsecond):
		fmt.Fprintln(os.Stderr, "Error: Read log pipe timedout")
	}

	fmt.Println("Get exit status")
	err = nil
	if exec_err != nil {
		if exiterr, ok := exec_err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			err = fmt.Errorf("cmd.Wait: %v", err)
		}
	}

	fmt.Println("Wait running report")
	//<-running_report_ch
	fmt.Println("Command finished ", err)
	return log_content, alarm_content, err
}

func killPreviousProcess(cmdline, proc_path, date string) error {
	dir_path := proc_path + "/" + date + "/"

	dirs, err := ioutil.ReadDir(dir_path)
	if err != nil {
		return err
	}

	for _, proc_dir := range dirs {
		if proc_dir.IsDir() {
			if exs_cmdline, err := ioutil.ReadFile(dir_path + proc_dir.Name() + "/cmd"); err == nil {
				if strings.Trim(string(exs_cmdline), "\r\n") == cmdline {
					pid_file := dir_path + proc_dir.Name() + "/pid"
					pid_file_stat, err := os.Stat(pid_file)
					if err != nil {
						fmt.Fprintln(os.Stderr, "run: Stat file ", pid_file, " failed:", err)
						continue
					}
					if pid_file_stat.ModTime().Before(time.Now().Add(-12 * time.Hour)) {
						continue
					}
					pid_bytes, err := ioutil.ReadFile(pid_file)
					if err != nil {
						fmt.Fprintln(os.Stderr, "run: Read file ", pid_file, " failed:", err)
						continue
					}
					pid, err := strconv.Atoi(strings.Trim(string(pid_bytes), " \r\n\t"))
					if err != nil {
						fmt.Fprintln(os.Stderr, "run: Parse pid ", string(pid_bytes), " failed:", err)
						continue
					}
					process, err := os.FindProcess(pid)
					if err == nil {
						err = process.Kill()
						if err != nil {
							fmt.Fprintln(os.Stderr, "run: Kill pid ", pid, " failed:", err)
							continue
						}
					}
				}
			}
		}
	}
	return nil
}

func readPipe(r *os.File, ch chan string) {
	d, err := ioutil.ReadAll(r)
	if err != nil {
		ch <- "ioutil.ReadAll failed:" + err.Error()
	} else {
		ch <- string(d)
	}
	close(ch)
}

func processAlarmPipe(r *os.File, ch chan string) {
	rd := bufio.NewReader(r)
	buf := make([]byte, 4*1024, 4*1024)
	var buffer bytes.Buffer

	for {
		n, err := rd.Read(buf)
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			fmt.Println("Read alarm pipe failed:", err)
			break
		} else {
			buffer.Write(buf[:n])
			content := string(buf[:n])
			SendAlarm(flow_inst_id, task_id, content, -1)
		}
	}
	ch <- buffer.String()
	close(ch)
}

type Bootstamp struct {
	running_reported  bool
	on_report_running func(args *RunningArgs) int
}

type RunningArgs struct {
	State int
	Pid   int
}

func (t *Bootstamp) ReportRunning(args *RunningArgs, reply *int) error {
	fmt.Println("Greet ", *args)
	*reply = 0
	if t.on_report_running != nil {
		*reply = t.on_report_running(args)
	}
	t.running_reported = true
	return nil
}

type dummyAddr string

func (a dummyAddr) Network() string {
	return string(a)
}

func (a dummyAddr) String() string {
	return string(a)
}

type noopConn struct{}

func (noopConn) LocalAddr() net.Addr                { return dummyAddr("local-addr") }
func (noopConn) RemoteAddr() net.Addr               { return dummyAddr("remote-addr") }
func (noopConn) SetDeadline(t time.Time) error      { return nil }
func (noopConn) SetReadDeadline(t time.Time) error  { return nil }
func (noopConn) SetWriteDeadline(t time.Time) error { return nil }

type rwRpcConn struct {
	io.Reader
	io.Writer
	noopConn

	closeFunc func() error // called if non-nil
	closec    chan bool    // else, if non-nil, send value to it on close
}

func (c *rwRpcConn) Close() error {
	if c.closeFunc != nil {
		return c.closeFunc()
	}
	select {
	case c.closec <- true:
	default:
	}
	return nil
}
