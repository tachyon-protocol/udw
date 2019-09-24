package udwCmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Cmd struct {
	cmd *exec.Cmd
}

func CmdBash(cmd string) *Cmd {
	return CmdSlice([]string{"bash", "-c", cmd})
}

func (c *Cmd) SetDir(path string) *Cmd {
	c.cmd.Dir = path
	return c
}

func (c *Cmd) PrintCmdLine() {
	c.FprintCmdLine(os.Stdout)
}

func (c *Cmd) FprintCmdLine(w io.Writer) {
	fmt.Fprintln(w, ">", strings.Join(c.cmd.Args, " "))
}

func (c *Cmd) Run() error {
	c.PrintCmdLine()
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return c.cmd.Run()
}

func (c *Cmd) ProxyRun() {
	err := c.Run()
	if err != nil {

		fmt.Println(err)
		os.Exit(2)
		return
	}
}

func (c *Cmd) GetExecCmd() *exec.Cmd {
	return c.cmd
}

func (c *Cmd) RunAndReturnOutput() (b []byte, err error) {
	c.PrintCmdLine()
	buf := &bytes.Buffer{}
	w := io.MultiWriter(buf, os.Stdout)
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	err = c.cmd.Run()
	return buf.Bytes(), err
}

func (c *Cmd) CombinedOutput() (b []byte, err error) {
	return c.cmd.CombinedOutput()
}

func (c *Cmd) RunAndTeeOutputToFile(path string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0777))
	if err != nil {
		return err
	}
	w := io.MultiWriter(f, os.Stdout)
	c.FprintCmdLine(w)
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	c.cmd.Stdin = os.Stdin
	return c.cmd.Run()
}

func (c *Cmd) StdioRun() error {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return c.cmd.Run()
}

func (c *Cmd) RunAndNotExitStatusCheck() error {
	err := c.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return nil
		}
		return err
	}
	return nil
}

func (c *Cmd) MustStdioRun() {
	err := c.StdioRun()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRunAndReturnOutput() (b []byte) {
	b, err := c.RunAndReturnOutput()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Cmd) MustRunAndReturnOutputAndNotExitStatusCheck() (b []byte) {
	b, err := c.RunAndReturnOutput()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return b
		}
		panic(err)
	}
	return b
}

func (c *Cmd) MustCombinedOutput() (b []byte) {
	b, err := c.cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Cmd) MustCombinedOutputWithErrorPrintln() (b []byte) {
	b, err := c.cmd.CombinedOutput()
	if err != nil {
		fmt.Println(">", strings.Join(c.cmd.Args, " "))
		os.Stdout.Write(b)
		panic(err)
	}
	return b
}

func (c *Cmd) MustCombinedOutputAndNotExitStatusCheck() (b []byte) {
	b, err := c.cmd.CombinedOutput()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return b
		}
		panic(err)
	}
	return b
}

func (c *Cmd) CombinedOutputAndNotExitStatusCheck() (b []byte, err error) {
	b, err = c.cmd.CombinedOutput()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return b, nil
		}
		return b, err
	}
	return b, nil
}

func (c *Cmd) MustCombinedOutputAndNotExitStatusCheckWithStdin(stdin []byte) (b []byte) {
	c.cmd.Stdin = bytes.NewBuffer(stdin)
	b, err := c.cmd.CombinedOutput()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return b
		}
		panic(err)
	}
	return b
}

func (c *Cmd) MustRunAndNotExitStatusCheck() {
	err := c.RunAndNotExitStatusCheck()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRun() {
	err := c.Run()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRunAndReturn() []byte {
	b, err := c.RunAndReturnOutput()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Cmd) MustRunWithWriter(w io.Writer) {
	c.FprintCmdLine(w)
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	err := c.cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRunWithWriterAndNotExitStatusCheck(w io.Writer) {
	c.FprintCmdLine(w)
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	err := c.cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			panic(err)
		}
	}
}

func (c *Cmd) MustOutputToWriter(w io.Writer) {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	err := c.cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) MustRunWithFullError() {
	b, err := c.RunAndReturnOutput()
	if err != nil {
		panic(fmt.Sprintln(string(b), err))
	}
}

func (c *Cmd) MustHiddenRunAndGetExitStatus() int {
	err := c.cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			panic(err)
		}
	}
	return GetExecCmdExitStatus(c.cmd)
}

func (c *Cmd) MustHiddenRunAndIsSuccess() bool {
	err := c.cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			panic(err)
		}
	}
	return c.cmd.ProcessState.Success()
}

func (c *Cmd) MustRunWithStdin(stdin []byte) {
	c.PrintCmdLine()
	c.cmd.Stdin = bytes.NewBuffer(stdin)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	err := c.cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (c *Cmd) RunAndReturnOutputWithStdin(stdin []byte) (b []byte, err error) {
	c.PrintCmdLine()
	buf := &bytes.Buffer{}
	w := io.MultiWriter(buf, os.Stdout)
	c.cmd.Stdin = bytes.NewBuffer(stdin)
	c.cmd.Stdout = w
	c.cmd.Stderr = w
	err = c.cmd.Run()
	return buf.Bytes(), err
}

type exitStatuser interface {
	ExitStatus() int
}

func GetExecCmdExitStatus(cmd *exec.Cmd) int {
	return cmd.ProcessState.Sys().(exitStatuser).ExitStatus()
}

func Exist(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil

}

func (c *Cmd) RunAndReturnOutputWithTimeout(timeout time.Duration) ([]byte, error) {
	c.PrintCmdLine()
	buf := &bytes.Buffer{}
	c.cmd.Stdout = buf
	c.cmd.Stderr = buf
	if err := c.cmd.Start(); err != nil {
		return buf.Bytes(), err
	}

	done := make(chan error)
	go func() {
		done <- c.cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):

		if err = c.cmd.Process.Kill(); err != nil && !strings.Contains(err.Error(), "process already finished") {
			panic(fmt.Errorf("failed to kill: %s, error: %s", c.cmd.Path, err))
		}
		go func() {
			<-done
		}()
		return buf.Bytes(), fmt.Errorf("Timeout. process:%s %v killed", c.cmd.Path, c.cmd.Args)
	case err = <-done:
		return buf.Bytes(), nil
	}
}

func (c *Cmd) MustStdoutRun() {
	c.PrintCmdLine()
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stdout
	err := c.cmd.Run()
	if err != nil {
		panic(err)
	}
}
