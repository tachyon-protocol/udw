// +build darwin dragonfly freebsd netbsd openbsd

package udwSysAsyncFd

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwSync"
	"io"
	"syscall"
	"unsafe"
)

func FdToRwc(tunFd int) (tun io.ReadWriteCloser, err error) {
	err = syscall.SetNonblock(tunFd, true)
	if err != nil {
		return nil, err
	}
	refd, err := getKqueneFd()
	if err != nil {
		return nil, err
	}
	var evs [2]syscall.Kevent_t
	*(*int)(unsafe.Pointer(&evs[0].Ident)) = tunFd
	evs[0].Filter = syscall.EVFILT_READ
	evs[0].Flags = syscall.EV_ADD | syscall.EV_CLEAR
	evs[0].Fflags = 0
	evs[0].Data = 0
	evs[1] = evs[0]
	evs[1].Filter = syscall.EVFILT_WRITE
	_, err = syscall.Kevent(refd, evs[:], nil, nil)
	if err != nil {
		return nil, err
	}

	tuns := &tunKquene{
		kqueuefd: refd,

		tunfd:       tunFd,
		readRequest: make(chan *epollReadRequest),
	}

	go tuns.kqueueThread()

	return tuns, nil
}

type epollReadRequest struct {
	data       []byte
	n          int
	err        error
	finishChan chan struct{}
}

type tunKquene struct {
	kqueuefd    int
	tunfd       int
	closer      udwClose.Closer
	readRequest chan *epollReadRequest

	writeBuf  udwBytes.QueueByteSlice
	writeCond udwSync.Cond
	canRead   bool
	canWrite  bool
}

func (e *tunKquene) Read(p []byte) (n int, err error) {
	for {
		e.writeCond.WaitCheckAndDo(func() bool {
			return e.canRead
		}, func() {
			n, err = syscall.Read(e.tunfd, p)
			if err != nil {
				if err == syscall.EAGAIN {
					e.canRead = false
				}
			}
		})
		if e.closer.IsClose() {
			return 0, udwNet.GetSocketCloseError()
		}
		if err != syscall.EAGAIN {
			return n, err
		}
	}
}
func (e *tunKquene) Write(p []byte) (n int, err error) {
	for {
		e.writeCond.WaitCheckAndDo(func() bool {
			return e.canWrite
		}, func() {
			n, err = syscall.Write(e.tunfd, p)
			if err != nil {
				if err == syscall.EAGAIN {
					e.canWrite = false
				}
			}
		})
		if e.closer.IsClose() {
			return 0, udwNet.GetSocketCloseError()
		}
		if err != syscall.EAGAIN {
			return n, err
		}
	}
}
func (e *tunKquene) Close() error {
	e.closer.CloseWithCallback(func() {
		syscall.Close(e.tunfd)
		syscall.Close(e.kqueuefd)
		e.writeCond.Close()
	})
	return nil
}
func (e *tunKquene) kqueueThread() {
	var events [2]syscall.Kevent_t
	ts := syscall.Timespec{
		Sec: 1,
	}
	for {

		nevent, err := syscall.Kevent(e.kqueuefd, nil, events[:], &ts)

		if e.closer.IsClose() {
			return
		}
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			udwLog.Log("udwSysAsyncFd", "syscall.Kevent error", nevent, err.Error())
			return
		}

		if nevent == 0 {
			continue
		}
		canRead := false
		canWrite := false
		for i := 0; i < nevent; i++ {
			if events[i].Filter == syscall.EVFILT_READ {
				canRead = true
			}
			if events[i].Filter == syscall.EVFILT_WRITE {
				canWrite = true
			}
		}
		e.writeCond.InLock(func() {
			if canRead {
				e.canRead = true
			}
			if canWrite {
				e.canWrite = true
			}
		})
		e.writeCond.Broadcast()
	}
}

func getKqueneFd() (fd int, err error) {
	kqFd, err := syscall.Kqueue()
	if err != nil {
		return 0, err
	}
	syscall.CloseOnExec(kqFd)
	return kqFd, nil
}
