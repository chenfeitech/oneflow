package subcmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"config"

	"github.com/codegangsta/cli"
)

func init() {
	Register(cli.Command{
		Name:    "kill",
		Aliases: []string{"k"},
		Usage:   "Kill task",
		Action:  kill,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "proc_path",
				Usage:  "task proc path",
				Value:  *config.ServerRoot + "/proc",
				EnvVar: "PROC_PATH",
			},
			cli.StringFlag{
				Name:  "uuid, u",
				Usage: "run uuid",
			},
			cli.StringFlag{
				Name:  "date, d",
				Usage: "run day",
			},
		},
	})
}

func kill(c *cli.Context) {
	proc_path := c.String("proc_path")
	uuid := c.String("uuid")
	date := strings.Replace(c.String("date"), "-", "", -1)
	dir_path := proc_path + "/" + date + "/"
	pid_file := dir_path + uuid + "/pid"
	agent_pid_file := dir_path + uuid + "/agent.pid"

	process_start_day, err := time.ParseInLocation("20060102", date, time.Local)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill: ", pid_file, " date parse failed:", err)
		return
	}

	kill2(pid_file, process_start_day)
	kill2(agent_pid_file, process_start_day)
}

func kill2(pid_file string, process_start_day time.Time) {
	_, err := os.Stat(pid_file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill: Stat file ", pid_file, " failed:", err)
		return
	}

	pid_bytes, err := ioutil.ReadFile(pid_file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill: Read file ", pid_file, " failed:", err)
		return
	}
	pid, err := strconv.Atoi(strings.Trim(string(pid_bytes), " \r\n\t"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill: Parse pid ", string(pid_bytes), " failed:", err)
		return
	}
	process, err := os.FindProcess(pid)
	if err == nil {
		ps, err := GetProcStat(pid)
		if err != nil {
			fmt.Fprintln(os.Stderr, "kill: GetProcStat ", pid, " failed:", err)
			return
		}
		if math.Abs(ps.StartTime().Sub(process_start_day).Hours()) >= 25 {
			fmt.Fprintln(os.Stderr, "kill: Process ", pid, " startup time not match.")
			fmt.Println("Process start day:", process_start_day)
			fmt.Println("Process start time:", ps.StartTime())
			return
		}

		process.Signal(syscall.SIGTERM)
		exited := make(chan struct{}, 1)
		go func() {
			process.Wait()
			exited <- struct{}{}
		}()

		select {
		case <-exited:
		case <-time.After(3 * time.Second):
			err = process.Kill()
			if err != nil {
				fmt.Fprintln(os.Stderr, "kill: Kill pid ", pid, " failed:", err)
				return
			}
		}

	}
}
