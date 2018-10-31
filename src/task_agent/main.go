package main

import (
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"task_agent/subcmd"
	"time"

	"github.com/codegangsta/cli"
)

func main() {
	debug.SetTraceback("crash")
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "task_agent"
	app.Usage = "Task agent command-line interface"
	app.Version = "0.2.0"
	app.Commands = subcmd.Commands()
	app.Run(os.Args)
}
