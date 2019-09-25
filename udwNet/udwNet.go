package udwNet

import (
	"fmt"
	"net"
	"strconv"
)

func MustGetLocalAddrFromListener(listener net.Listener) string {
	return MustGetLocalAddrFromAddr(listener.Addr())
}

func MustGetLocalAddrFromAddr(addr net.Addr) string {
	tcpAddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		panic(err)
	}
	return "127.0.0.1:" + strconv.Itoa(tcpAddr.Port)
}

func MustTcpRandomListen() net.Listener {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return l
}

func CloseRead(conn net.Conn) error {
	tcpC := mustGetTcpConnFromConn(conn)
	return tcpC.CloseRead()
}

func CloseWrite(conn net.Conn) error {
	tcpC := mustGetTcpConnFromConn(conn)
	return tcpC.CloseWrite()
}

func mustGetTcpConnFromConn(conn net.Conn) *net.TCPConn {
	tcpC, ok := conn.(*net.TCPConn)
	if ok {
		return tcpC
	}
	conner, ok := conn.(GetUnderlyingConner)
	if ok {
		return mustGetTcpConnFromConn(conner.GetUnderlyingConn())
	}
	panic(fmt.Errorf("not support conn type %T", conn))
}

type GetUnderlyingConner interface {
	GetUnderlyingConn() net.Conn
}
