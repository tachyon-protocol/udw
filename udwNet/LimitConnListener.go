package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwClose"
	"net"
	"sync"
	"time"
)

const debugLimitConnTcpListener = false

type LimitConnListenerRequest struct {
	OriginLn      *net.TCPListener
	MaxConnNumber int
}

func NewLimitConnTcpListener(req LimitConnListenerRequest) (nextLn net.Listener) {
	return &limitConnListener{
		req: req,
	}
}

type limitConnListener struct {
	req                  LimitConnListenerRequest
	currentConnNum       int
	currentConnNumLocker sync.Mutex
}

func (ln *limitConnListener) Accept() (c net.Conn, err error) {
	for {
		tc, err := ln.req.OriginLn.AcceptTCP()
		if err != nil {
			return nil, err
		}
		ln.currentConnNumLocker.Lock()
		if ln.currentConnNum > ln.req.MaxConnNumber {
			ln.currentConnNumLocker.Unlock()
			tc.Close()
			if debugLimitConnTcpListener {
				fmt.Println("LimitConnTcpListener too many connection.", ln.req.MaxConnNumber)
			}
			continue
		}
		ln.currentConnNum++
		ln.currentConnNumLocker.Unlock()
		tc.SetKeepAlive(true)

		tc.SetKeepAlivePeriod(10 * time.Minute)
		conn2 := &limitConnListenerConn{
			Conn: tc,
		}
		conn2.closer.AddOnClose(func() {
			ln.currentConnNumLocker.Lock()
			ln.currentConnNum--
			ln.currentConnNumLocker.Unlock()
		})
		return conn2, nil
	}

}

func (ln *limitConnListener) Close() error {
	return ln.req.OriginLn.Close()
}

func (ln *limitConnListener) Addr() net.Addr {
	return ln.req.OriginLn.Addr()
}

type limitConnListenerConn struct {
	net.Conn
	closer udwClose.Closer
}

func (conn *limitConnListenerConn) Close() error {
	conn.closer.Close()
	return conn.Conn.Close()
}
