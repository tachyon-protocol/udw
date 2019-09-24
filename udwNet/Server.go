package udwNet

import (
	"net"
	"strconv"
)

type Server interface {
	Start() error

	Close() error

	Addr() (net.Addr, error)
}

func MustGetServerAddrString(s Server) string {
	addr, err := s.Addr()
	if err != nil {
		panic(err)
	}
	return addr.String()
}

func MustGetServerLocalAddrString(s Server) string {
	addr, err := s.Addr()
	if err != nil {
		panic(err)
	}
	port, err := PortFromNetAddr(addr)
	if err != nil {
		panic(err)
	}
	return "127.0.0.1:" + strconv.Itoa(port)
}

func MustServerStart(s Server) {
	err := s.Start()
	if err != nil {
		panic(err)
	}
}

type FuncServer struct {
	StartFunc func() error
	CloseFunc func() error
	AddrFunc  func() (net.Addr, error)
	ExistAddr net.Addr
}

func (s *FuncServer) Start() error {
	return s.StartFunc()
}

func (s *FuncServer) Close() error {
	return s.CloseFunc()
}

func (s *FuncServer) Addr() (net.Addr, error) {
	if s.ExistAddr != nil {
		return s.ExistAddr, nil
	}
	return s.AddrFunc()
}
