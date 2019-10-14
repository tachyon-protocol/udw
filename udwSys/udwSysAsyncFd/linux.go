// +build linux

package udwSysAsyncFd

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwSync"
	"io"
	"syscall"
)

func FdToRwc(tunFd int) (tun io.ReadWriteCloser, err error) {
	err = syscall.SetNonblock(tunFd, true)
	if err != nil {
		return nil, err
	}
	refd, err := getEpollFdAndAddTunFd(tunFd, syscall.EPOLLIN)
	if err != nil {
		return nil, err
	}
	wefd, err := getEpollFdAndAddTunFd(tunFd, syscall.EPOLLOUT)
	if err != nil {
		return nil, err
	}

	tuns := &tunEpoll{
		refd:        refd,
		wefd:        wefd,
		tunfd:       tunFd,
		readRequest: make(chan *epollReadRequest),
	}
	go tuns.readThread()
	go tuns.writeThread()
	return tuns, nil
}

type epollReadRequest struct {
	data       []byte
	n          int
	err        error
	finishChan chan struct{}
}

type tunEpoll struct {
	refd        int
	wefd        int
	tunfd       int
	closer      udwClose.Closer
	readRequest chan *epollReadRequest

	writeBuf  udwBytes.QueueByteSlice
	writeCond udwSync.Cond
}

func (e *tunEpoll) Read(p []byte) (n int, err error) {
	if e.closer.IsClose() {
		return 0, udwNet.GetSocketCloseError()
	}

	req := &epollReadRequest{
		data:       p,
		finishChan: make(chan struct{}),
	}
	select {
	case e.readRequest <- req:
		<-req.finishChan
		return req.n, req.err
	case <-e.closer.GetCloseChan():
		return 0, udwNet.GetSocketCloseError()
	}
}
func (e *tunEpoll) Write(p []byte) (n int, err error) {
	if e.closer.IsClose() {
		return 0, udwNet.GetSocketCloseError()
	}

	e.writeCond.InLock(func() {
		e.writeBuf.AddOne(p)

	})
	e.writeCond.Signal()
	return len(p), nil
}
func (e *tunEpoll) Close() error {
	e.closer.CloseWithCallback(func() {
		syscall.Close(e.tunfd)
		syscall.Close(e.refd)
		syscall.Close(e.wefd)
		e.writeCond.Close()
	})
	return nil
}
func (e *tunEpoll) readThread() {
	var events [1]syscall.EpollEvent
	for {
		nevents, err := syscall.EpollWait(e.refd, events[:], 1000)
		if e.closer.IsClose() {
			return
		}
		if err != nil {
			if udwNet.IsInterruptedSystemCall(err) {

				continue
			}
			udwLog.Log("error", "[establishTunEpoll.EpollWait]", err.Error())
			e.Close()
			return
		}

		if nevents > 0 && events[0].Events&syscall.EPOLLIN != 0 {

			select {
			case rr := <-e.readRequest:
				n, err := syscall.Read(e.tunfd, rr.data)
				rr.n = n
				rr.err = err
				rr.finishChan <- struct{}{}
			case <-e.closer.GetCloseChan():
				return
			}
		}
	}
}
func (e *tunEpoll) writeThread() {
	var events [1]syscall.EpollEvent
	for {
		nevents, err := syscall.EpollWait(e.wefd, events[:], 1000)
		if e.closer.IsClose() {
			return
		}
		if err != nil {
			udwLog.Log("error", "[tunEpoll.writeThread.EpollWait]", err.Error())
			e.Close()
			return
		}
		if nevents > 0 && events[0].Events&syscall.EPOLLOUT != 0 {
			thisData := []byte(nil)
			e.writeCond.WaitCheckAndDo(func() bool {
				return e.writeBuf.HasData()
			}, func() {
				thisData = e.writeBuf.GetOne()
			})
			if len(thisData) == 0 {
				continue
			}
			_, err := syscall.Write(e.tunfd, thisData)
			if err != nil {
				udwLog.Log("error", "[tunEpoll.writeThread.syscall.Write]", err.Error())
				e.Close()
				return
			}
			e.writeCond.InLock(func() {
				e.writeBuf.RemoveOne()
			})
		}

	}
}

func getEpollFd() (fd int, err error) {
	fd, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err == nil {
		return fd, nil
	}
	fd, err = syscall.EpollCreate(1024)
	if err == nil {
		return fd, nil
	}
	return 0, fmt.Errorf("[getEpollFd] EpollCreate fail %s", err)
}

func getEpollFdAndAddTunFd(tunfd int, event uint32) (fd int, err error) {
	refd, err := getEpollFd()
	if err != nil {
		return 0, err
	}
	var ev syscall.EpollEvent
	ev.Events = event
	ev.Fd = int32(tunfd)
	err = syscall.EpollCtl(refd, syscall.EPOLL_CTL_ADD, tunfd, &ev)
	if err != nil {
		return 0, fmt.Errorf("[getEpollFdAndAddTunFd.EpollCtl] %s", err)
	}
	return refd, nil
}
