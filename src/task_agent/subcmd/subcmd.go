package subcmd

import "github.com/codegangsta/cli"

var (
	cmds = make([]cli.Command, 0)
)

func Register(cmd cli.Command) {
	cmds = append(cmds, cmd)
}

func Commands() []cli.Command {
	return cmds
}
