package udwNet

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"net"
	"strconv"
	"sync"
	"testing"
)

func TestTcpListener_Close(t *testing.T) {
	listenAddr := ":1234"
	msg := "Winter is here"
	repeat := 100
	equalCount := 0
	equalCountLock := sync.Mutex{}
	for i := 0; i < repeat; i++ {
		readDoneChan := make(chan struct{})
		_msg := []byte(msg + strconv.Itoa(i))
		connCloser := udwClose.NewCloser()
		listener := TcpNewListenerReturnListener(listenAddr, func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, len(_msg))
			for {
				if connCloser.IsClose() {
					return
				}
				_, err := io.ReadFull(conn, buf)
				if err != nil {
					if connCloser.IsClose() {
						return
					}
					udwErr.PanicIfError(err)
				}
				udwTest.Equal(buf, _msg)
				close(readDoneChan)
				equalCountLock.Lock()
				equalCount++
				equalCountLock.Unlock()
			}
		})
		conn, err := net.Dial("tcp", "127.0.0.1"+listenAddr)
		udwErr.PanicIfError(err)
		connCloser.AddOnClose(func() {
			conn.Close()
		})
		_, err = conn.Write(_msg)
		udwErr.PanicIfError(err)
		<-readDoneChan
		listener.Close()
		connCloser.Close()
	}
	equalCountLock.Lock()
	udwTest.Equal(equalCount, repeat)
	equalCountLock.Unlock()
}
