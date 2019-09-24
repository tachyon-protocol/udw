package udwNet

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwErr"
	"net"
	"sync"
)

func TcpNewListener(listenAddr string, processor func(conn net.Conn)) (closer func()) {
	tl := TcpNewListenerReturnListener(listenAddr, processor)
	return tl.Close
}

func TcpNewListenerReturnListener(listenAddr string, processor func(conn net.Conn)) (tl *TcpListener) {
	listener, err := getListen()("tcp", listenAddr)
	udwErr.PanicIfError(err)
	return tcpHandleConnFromExistListener(listener, processor)
}

func TcpNewListenerFromExistListener(listener net.Listener, processor func(conn net.Conn)) (closer func()) {
	tl := tcpHandleConnFromExistListener(listener, processor)
	return tl.Close
}

type TcpListener struct {
	listener        net.Listener
	closer          udwClose.Closer
	acceptCloseChan chan struct{}
	listenPort      int
	processor       func(conn net.Conn)
}

func tcpHandleConnFromExistListener(listener net.Listener, connHandler func(conn net.Conn)) *TcpListener {
	tl := &TcpListener{
		listener:        listener,
		acceptCloseChan: make(chan struct{}),
		processor:       connHandler,
	}
	tl.listenPort = MustGetPortFromNetAddrIgnoreNotFound(listener.Addr())
	go tl.readLoop()
	return tl
}

func (tl *TcpListener) GetListenPort() int {
	return tl.listenPort
}

func (tl *TcpListener) Close() {
	tl.closer.CloseWithCallback(func() {
		tl.listener.Close()
		<-tl.acceptCloseChan
	})
}
func (tl *TcpListener) readLoop() {
	for {
		if tl.closer.IsClose() {
			close(tl.acceptCloseChan)
			return
		}
		conn, err := tl.listener.Accept()
		if err != nil {
			close(tl.acceptCloseChan)
			if IsSocketCloseError(err) && tl.closer.IsClose() {
				return
			}
			tl.Close()
			panic(err)
		}
		go tl.processor(conn)
	}
}

func TcpNewDialer() (response NewDialerResponse) {
	return NewDialerResponse{
		Dialer: net.Dial,
		Closer: func() {},
	}
}

var gMockedListen func(network, address string) (net.Listener, error) = net.Listen
var gMockedListenLock sync.Mutex

func SetMockedListen(f func(network, address string) (net.Listener, error)) {
	gMockedListenLock.Lock()
	gMockedListen = f
	gMockedListenLock.Unlock()
}

func getListen() func(network, address string) (net.Listener, error) {
	gMockedListenLock.Lock()
	defer gMockedListenLock.Unlock()
	return gMockedListen
}
