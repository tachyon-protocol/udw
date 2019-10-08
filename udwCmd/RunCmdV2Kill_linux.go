package udwCmd

import (
	"os"
	"syscall"
)

func (resp *CmdV2Response) init() {
	resp.execCmd = ExecCmd(resp.req.CmdSlice)

	resp.execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	resp.execCmd.Stdin = os.Stdin
	resp.execCmd.Stdout = os.Stdout
	resp.execCmd.Stderr = os.Stderr
}

func (resp *CmdV2Response) MustKill() {
	execCmd := resp.execCmd
	if execCmd == nil {
		panic("[CmdV2Response.MustKill] resp.execCmd==nil")
	}

	err := syscall.Kill(-execCmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		panic(err)
	}
}
