package udwNet

import (
	"net"
)

type NewDialerResponse struct {
	Dialer Dialer
	Closer func()
}

type NewListener func(listenAddr string, processor func(conn net.Conn)) (Closer func())
type NewDialer func() (response NewDialerResponse)

type TunImp struct {
	NewDialer   func() (response NewDialerResponse)
	NewListener func(listenAddr string, processor func(conn net.Conn)) (Closer func())
}
