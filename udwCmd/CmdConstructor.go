package udwCmd

import (
	"strings"
)

func CmdSlice(args []string) *Cmd {
	if len(args) == 0 {
		panic("[CmdSlice] need the path of the command")
	}
	return &Cmd{
		cmd: ExecCmd(args),
	}
}

func CmdString(cmd string) *Cmd {
	if cmd == "" {
		panic("[CmdString] need the path of the command")
	}
	args := strings.Split(cmd, " ")
	return &Cmd{
		cmd: ExecCmd(args),
	}
}
