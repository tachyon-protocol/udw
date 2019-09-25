package udwCmd

import (
	"os/exec"
	"strings"
)

func Run(cmd string) (err error) {
	return CmdString(cmd).Run()
}

func StdioSliceRun(args []string) (err error) {
	return CmdSlice(args).StdioRun()
}

func MustRun(cmd string) {
	CmdString(cmd).MustRun()
}

func ProxyRun(cmd string) {
	CmdString(cmd).ProxyRun()
}

func MustRunNotExistStatusCheck(cmd string) {
	CmdString(cmd).MustRunAndNotExitStatusCheck()
}

func MustRunAndReturnOutput(cmd string) []byte {
	return CmdString(cmd).MustRunAndReturnOutput()
}

func MustCombinedOutput(cmd string) []byte {
	return CmdString(cmd).MustCombinedOutput()
}

func MustCombinedOutputWithErrorPrintln(cmd string) []byte {
	return CmdString(cmd).MustCombinedOutputWithErrorPrintln()
}

func MustCombinedOutputAndNotExitStatusCheck(cmd string) []byte {
	return CmdString(cmd).MustCombinedOutputAndNotExitStatusCheck()
}

func MustRunInBash(cmd string) {
	CmdBash(cmd).MustRun()
}

func MustRunInBashAndReturn(cmd string) []byte {
	b := CmdBash(cmd).MustRunAndReturn()
	return b
}

func IsErrorExitStatus(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*exec.ExitError)
	if ok {
		return true
	}
	_errStr := strings.Trim(err.Error(), "\n")
	if strings.HasPrefix(_errStr, "exit status") {
		return true
	}
	return false
}
