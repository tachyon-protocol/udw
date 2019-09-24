package udwNet

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwLog"
	"net"
)

const debugMustRunUdpListener = false

type UdpReadContext struct {
	Conn *net.UDPConn
	Data []byte
	Addr *net.UDPAddr
}

func MustRunUdpListener(listenAddr string, processor func(ctx UdpReadContext)) *UdpListener {
	listenAddrI, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		panic(err)
	}
	listenConn, err := net.ListenUDP("udp", listenAddrI)
	if err != nil {
		panic(err)
	}
	listenPort := MustGetPortFromNetAddr(listenConn.LocalAddr())
	obj := &UdpListener{
		conn:           listenConn,
		innerCloseChan: make(chan struct{}),
		processor:      processor,
		listenPort:     listenPort,
	}
	go obj.readLoop()
	return obj
}

type UdpListener struct {
	conn           *net.UDPConn
	innerCloseChan chan struct{}
	processor      func(ctx UdpReadContext)
	listenPort     int
	closer         udwClose.Closer
}

func (l *UdpListener) GetListenConn() *net.UDPConn {
	return l.conn
}
func (l *UdpListener) GetListenPort() int {
	return l.listenPort
}
func (l *UdpListener) Close() {
	l.closer.CloseWithCallback(func() {
		l.conn.Close()
		<-l.innerCloseChan
	})
}
func (l *UdpListener) IsClose() bool {
	return l.closer.IsClose()
}
func (l *UdpListener) readLoop() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := l.conn.ReadFromUDP(buf)
		if err != nil {
			if IsSocketCloseError(err) {
				close(l.innerCloseChan)
				if debugMustRunUdpListener && l.closer.IsClose() == false {
					udwLog.Log("debug", "[udwNet.UdpListener.readLoop] IsSocketCloseError", err)
					l.Close()
				}
				return
			}

			panic("[udwNet.UdpListener.readLoop] ReadFromUDP " + err.Error())

		}
		l.processor(UdpReadContext{
			Conn: l.conn,
			Data: buf[:n],
			Addr: addr,
		})
	}
}
