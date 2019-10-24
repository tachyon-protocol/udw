// +build ios macAppStore

package udwSysLog

import (
	"bufio"
	"os"
	"sync"
	"syscall"
)

func lineLog(f *os.File) {
	const logSize = 1024
	r := bufio.NewReaderSize(f, logSize)
	for {
		line, _, err := r.ReadLine()
		str := string(line)
		if err != nil {
			str += " " + err.Error()
		}
		log(str)
		if err != nil {
			break
		}
	}
}

var (
	gPipeWriter *os.File
)
var gStdOutRedirectIosOnce sync.Once

func StdOutRedirectIos() {
	gStdOutRedirectIosOnce.Do(func() {
		r, w, err := os.Pipe()
		if err != nil {
			log("[stdOutErrRedirectL1]1 os.Pipe fail " + err.Error())
			return
		}
		gPipeWriter = w
		err = syscall.Dup2(int(w.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			log("[stdOutErrRedirectL1]2 syscall.Dup3 fail " + err.Error())
			return
		}
		go lineLog(r)
	})

}
