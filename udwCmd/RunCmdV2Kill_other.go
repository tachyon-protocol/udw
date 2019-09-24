package udwCmd

import (
	"fmt"
	"os"
)

func (resp *CmdV2Response) init() {
	resp.execCmd = ExecCmd(resp.req.CmdSlice)
	resp.execCmd.Stdin = os.Stdin
	resp.execCmd.Stdout = os.Stdout
	resp.execCmd.Stderr = os.Stderr
}

func (resp *CmdV2Response) MustKill() {
	execCmd := resp.execCmd
	if execCmd == nil {
		panic("[CmdV2Response.MustKill] resp.execCmd==nil")
	}
	fmt.Println("[WARNING] current platform not support kill all child processes")
	err := execCmd.Process.Kill()
	if err != nil {
		panic(err)
	}
}
