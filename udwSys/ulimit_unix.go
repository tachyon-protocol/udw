// +build linux darwin

package udwSys

import (
	"fmt"
	"syscall"
)

func SetCurrentMaxFileNum(limit uint64) (err error) {
	rLimit := syscall.Rlimit{
		Max: limit,
		Cur: limit,
	}
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("syscall.SetrLimit: %s", err)
	}
	return
}
