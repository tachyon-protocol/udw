package udwNetTestV2

import (
	"fmt"
	"net"
	"sync"
)

type InmemoryListener struct {
	lock                      sync.Mutex
	closed                    bool
	conns                     chan net.Conn
	connInnerBufferPacketSize int
}

func NewInmemoryListener() *InmemoryListener {
	return &InmemoryListener{
		conns:                     make(chan net.Conn, 1024),
		connInnerBufferPacketSize: 4,
	}
}

func (ln *InmemoryListener) SetConnInnerBufferPacketSize(size int) {
	ln.connInnerBufferPacketSize = size
}

func (ln *InmemoryListener) Accept() (net.Conn, error) {
	c, ok := <-ln.conns
	if !ok {
		return nil, fmt.Errorf("InmemoryListener is already closed: use of closed network connection")
	}
	return c, nil
}

func (ln *InmemoryListener) Close() error {
	var err error

	ln.lock.Lock()
	if !ln.closed {
		close(ln.conns)
		ln.closed = true
	} else {
		err = fmt.Errorf("InmemoryListener is already closed")
	}
	ln.lock.Unlock()
	return err
}

func (ln *InmemoryListener) Addr() net.Addr {
	return &net.UnixAddr{
		Name: "InmemoryListener",
		Net:  "memory",
	}
}

func (ln *InmemoryListener) Dial() (net.Conn, error) {
	pc := NewPipeConns(ln.connInnerBufferPacketSize)
	cConn := pc.Conn1()
	sConn := pc.Conn2()
	ln.lock.Lock()
	if !ln.closed {
		ln.conns <- sConn
	} else {
		sConn.Close()
		cConn.Close()
		cConn = nil
	}
	ln.lock.Unlock()

	if cConn == nil {
		return nil, fmt.Errorf("InmemoryListener is already closed")
	}
	return cConn, nil
}
