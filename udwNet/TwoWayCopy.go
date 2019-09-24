package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwDebug"
	"net"
	"sync"
	"time"
)

type ConnTwoWayCopyRequest struct {
	FromHopConn  net.Conn
	NextHopConn  net.Conn
	AddTimeout   time.Duration
	IsPrintError bool
}

var gTwoWayCopyBufPool udwBytes.BufWriterPool

func ConnTwoWayCopy(req ConnTwoWayCopyRequest) {
	var closeOnce sync.Once
	timeOutDiff := time.Second
	var lastSetTimeoutTime time.Time
	var lastSetTimeoutTimeSecondLocker sync.Mutex
	if req.AddTimeout > 0 {
		timeOutDiff = time.Duration(float64(req.AddTimeout) / float64(10))
		if timeOutDiff > time.Second {
			timeOutDiff = time.Second
		}
	}
	setDeadLineFn := func() (errMsg string) {
		thisSetTimeout := time.Now().Add(req.AddTimeout)
		thisSetTimeout2 := thisSetTimeout.Add(timeOutDiff)
		lastSetTimeoutTimeSecondLocker.Lock()
		if lastSetTimeoutTime.After(thisSetTimeout) {
			lastSetTimeoutTimeSecondLocker.Unlock()
			return ""
		}
		lastSetTimeoutTime = thisSetTimeout2
		lastSetTimeoutTimeSecondLocker.Unlock()
		err := req.FromHopConn.SetDeadline(thisSetTimeout2)
		if err != nil {
			return "gd9dvqw7zf " + err.Error()
		}
		err = req.NextHopConn.SetDeadline(thisSetTimeout2)
		if err != nil {
			return "6m6veb2g76 " + err.Error()
		}
		return ""
	}
	copyFn := func(conn1 net.Conn, conn2 net.Conn) (isCloseError bool, errMsg string) {
		var err error
		var n int

		pw := gTwoWayCopyBufPool.Get()
		buf := pw.GetHeadBuffer(32 * 1024)
		for {
			if req.AddTimeout > 0 {
				errMsg = setDeadLineFn()
				if errMsg != "" {
					break
				}
			}
			n, err = conn1.Read(buf)
			if err != nil {
				errMsg = "vtveq5a52z " + err.Error()
				break
			}
			if n == 0 {
				continue
			}
			_, err = conn2.Write(buf[:n])
			if err != nil {
				errMsg = "e5jnj8jhup " + err.Error()
				break
			}
		}
		gTwoWayCopyBufPool.Put(pw)
		return IsSocketCloseError(err), errMsg
	}

	f := func(conn1 net.Conn, conn2 net.Conn) {
		iscloseError, errMsg := copyFn(conn1, conn2)
		closeOnce.Do(func() {
			conn1.Close()
			conn2.Close()
		})
		if req.IsPrintError && iscloseError == false && errMsg != "" {
			fmt.Println("dtwxg8dsua", errMsg, udwDebug.GetCurrentAllStackString(1))
			return
		}
	}
	go f(req.FromHopConn, req.NextHopConn)
	f(req.NextHopConn, req.FromHopConn)
}
