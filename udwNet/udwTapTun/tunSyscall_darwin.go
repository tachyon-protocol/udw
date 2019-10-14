package udwTapTun

import (
	"os"
	"syscall"
	"unsafe"
)

func SyscallIoctl(fd int, request, argp uintptr) error {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
	if errorp == 0 {
		return nil
	}
	return os.NewSyscallError("ioctl", errorp)
}

func SyscallConnect(s int, addr uintptr, addrlen uintptr) (err error) {
	_, _, e1 := syscall.Syscall(syscall.SYS_CONNECT, uintptr(s), uintptr(addr), uintptr(addrlen))
	if e1 != 0 {
		return os.NewSyscallError("connect", e1)
	}
	return nil
}

func SyscallFcntl(fd int, cmd int, arg int) (val int, err error) {
	r0, _, e1 := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	if e1 != 0 {
		return 0, os.NewSyscallError("fcntl", e1)
	}
	val = int(r0)
	return val, nil
}

func SyscallGetSockopt(s int, level int, name int, val uintptr, vallen *uint32) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
	if e1 != 0 {
		return os.NewSyscallError("getsockopt", e1)
	}
	return nil
}
