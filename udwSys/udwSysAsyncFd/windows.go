// +build windows

package udwSysAsyncFd

import (
	"io"
	"syscall"
)

type tTun struct {
	fd syscall.Handle
}

func FdToRwc(tunFd int) (tun io.ReadWriteCloser, err error) {
	return &tTun{
		fd: syscall.Handle(tunFd),
	}, nil
}

func (t *tTun) Read(buf []byte) (n int, err error) {
	var done uint32
	event := syscallCreateEvent()
	defer syscall.CloseHandle(event)
	over := &syscall.Overlapped{
		HEvent: event,
	}
	err = syscall.ReadFile(syscall.Handle(t.fd), buf, &done, over)
	if err == nil {
		return int(over.InternalHigh), nil
	}
	if err != nil && err != syscall.ERROR_IO_PENDING {
		return 0, err
	}
	_, err = syscall.WaitForSingleObject(event, iNFINITE)
	if err != nil {
		return 0, err
	}
	return int(over.InternalHigh), nil
}

func (t *tTun) Write(buf []byte) (n int, err error) {
	var done uint32
	event := syscallCreateEvent()
	defer syscall.CloseHandle(event)
	over := &syscall.Overlapped{
		HEvent: event,
	}
	err = syscall.WriteFile(syscall.Handle(t.fd), buf, &done, over)
	if err == nil {
		return int(over.InternalHigh), nil
	}
	if err != nil && err != syscall.ERROR_IO_PENDING {
		return 0, err
	}
	_, err = syscall.WaitForSingleObject(event, iNFINITE)
	if err != nil {
		return 0, err
	}
	return int(over.InternalHigh), nil
}

func (t *tTun) Close() (err error) {
	return syscall.CloseHandle(t.fd)
}

var (
	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procCreateEvent = modkernel32.NewProc("CreateEventW")
)

const iNFINITE = 0xffffffff

func syscallCreateEvent() (handler syscall.Handle) {
	handleru, _, err := syscall.Syscall6(procCreateEvent.Addr(), 4, 0, 1, 1, 0, 0, 0)
	if err != 0 {
		panic(err)
	}
	return syscall.Handle(handleru)
}
