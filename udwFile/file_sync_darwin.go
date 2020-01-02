// +build darwin,!go1.12

package udwFile

import (
	"os"
	"syscall"
)

func FileSync(f *os.File) error {
	_, _, err := syscall.Syscall(syscall.SYS_FCNTL, f.Fd(), syscall.F_FULLFSYNC, 0)
	if err == 0 {
		return nil
	}
	return err
}
