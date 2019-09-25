package udwCmd

import "os/exec"

func ExecCmd(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	return cmd
}
