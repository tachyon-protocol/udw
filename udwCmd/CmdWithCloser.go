package udwCmd

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"os"
	"sync"
)

func (c *Cmd) MustAsyncRunWithCloser() *udwClose.Closer {
	c.PrintCmdLine()
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	err := c.cmd.Start()
	if err != nil {
		panic(err)
	}
	closer := udwClose.NewCloser()
	thisProcess := c.cmd.Process
	thisProcessLocker := sync.Mutex{}

	go func() {
		c.cmd.Wait()
		thisProcessLocker.Lock()
		thisProcess = nil
		thisProcessLocker.Unlock()
		closer.Close()
	}()
	closer.AddOnClose(func() {
		thisProcessLocker.Lock()
		if thisProcess == nil {
			thisProcessLocker.Unlock()
			return
		}
		thisProcess2 := thisProcess
		thisProcess = nil
		thisProcessLocker.Unlock()

		thisProcess2.Kill()
	})
	return closer
}

func (c *Cmd) MustAsyncRun() {
	c.PrintCmdLine()
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	err := c.cmd.Start()
	if err != nil {
		panic(err)
	}
	return
}
