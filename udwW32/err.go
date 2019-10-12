// +build windows

package udwW32

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"syscall"
)

func SyscallErrorToMsg(err error) string {
	if err == nil {
		return ""
	}
	switch errO := err.(type) {
	case syscall.Errno:
		if errO == 0 {
			return ""
		}
		return "W32: " + udwStrconv.FormatUint64Hex(uint64(uintptr(errO))) + " " + errO.Error()
	default:
		return err.Error()
	}
}

func IsSyscallErrorHappen(err error) bool {
	if err == nil {
		return false
	}
	errNo, ok := err.(syscall.Errno)
	if ok && errNo == 0 {
		return false
	}
	return true
}
