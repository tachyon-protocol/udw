package udwNet

import (
	"io"
	"net"
)

type ConnHandler interface {
	ConnHandle(conn net.Conn)
}

type ConnHandlerFunc func(conn net.Conn)

func (f ConnHandlerFunc) ConnHandle(conn net.Conn) {
	f(conn)
}

type ConnServer struct {
	Listener net.Listener
	Handler  ConnHandler
	Closer   io.Closer
}

func (server *ConnServer) Close() (err error) {
	err = server.Listener.Close()
	var err1 error
	if server.Closer != nil {
		err1 = server.Closer.Close()
	}
	if err != nil {
		return err
	}
	return err1
}

func (server *ConnServer) Start() (err error) {
	go func() {
		defer server.Listener.Close()
		for {
			conn, err := server.Listener.Accept()
			if err != nil {
				if IsSocketCloseError(err) {
					return
				}
				panic(err)

			}
			go server.Handler.ConnHandle(conn)
		}
	}()
	return nil
}

func (server *ConnServer) Addr() (net.Addr, error) {
	return server.Listener.Addr(), nil
}

func NewTCPServer(listenAddr string, hander ConnHandler, closer io.Closer) (s *ConnServer, err error) {
	s = &ConnServer{}
	s.Listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	s.Handler = hander
	s.Closer = closer
	return s, nil
}

func RunTCPServerV2(Listener net.Listener, handle ConnHandlerFunc) (closer func() error) {
	go func() {
		for {
			conn, err := Listener.Accept()
			if err != nil {
				if IsSocketCloseError(err) {
					return
				}
				panic(err)

			}
			go handle(conn)
		}
	}()
	return Listener.Close
}

func RunTCPServerListenAddr(listenAddr string, handle ConnHandlerFunc) (closer func() error) {
	return RunTCPServerV2(MustListen("tcp", listenAddr), handle)
}

func MustListen(network string, address string) net.Listener {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	return listener
}
