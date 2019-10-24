// +build android

package udwSysLog

import (
	"bufio"
	"os"
	"sync"
	"syscall"
)

func StdOutErrRedirectAndroid() {
	gStdOutErrRedirectOnce.Do(func() {
		stdOutErrRedirectAndroidL1()
	})
}

var gStdOutErrRedirectOnce sync.Once

var (
	gPipeWriter *os.File
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

func stdOutErrRedirectAndroidL1() {
	r, w, err := os.Pipe()
	if err != nil {
		log("[stdOutErrRedirectL1]1 os.Pipe fail " + err.Error())
		return
	}
	gPipeWriter = w
	if err := syscall.Dup3(int(w.Fd()), int(os.Stderr.Fd()), 0); err != nil {
		log("[stdOutErrRedirectL1]2 syscall.Dup3 fail " + err.Error())
		return
	}
	if err := syscall.Dup3(int(w.Fd()), int(os.Stdout.Fd()), 0); err != nil {
		log("[stdOutErrRedirectL1]4 syscall.Dup3 fail " + err.Error())
		return
	}
	go lineLog(r)
}
