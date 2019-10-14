// +build !js

package udwConsole

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit() {

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}

func WaitForExitOrCloser(closer *udwClose.Closer) {

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-ch:
		return
	case <-closer.GetCloseChan():
		return
	}
}
