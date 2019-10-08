package udwNet

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwIo"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

func (conn *fasterTcpConn) Read(b []byte) (nr int, err error) {
	err = conn.TCPConn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		return 0, err
	}
	nr, err = conn.TCPConn.Read(b)
	if err != nil {
		return nr, err
	}
	conn.closeLock.Lock()
	defer conn.closeLock.Unlock()
	fdNum := int(conn.fd.Fd())
	if fdNum == -1 {
		return nr, err
	}
	err = syscall.SetsockoptInt(fdNum, syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1)
	if err != nil {
		udwErr.LogErrorWithStack(err)
		return
	}
	return
}
func (conn *fasterTcpConn) Write(b []byte) (nr int, err error) {
	err = conn.TCPConn.SetWriteDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		return 0, err
	}
	return conn.TCPConn.Write(b)
}

func (conn *fasterTcpConn) Close() (err error) {
	conn.closeLock.Lock()
	defer conn.closeLock.Unlock()

	return udwIo.MultiErrorHandle(conn.TCPConn.CloseRead, conn.TCPConn.Close, conn.fd.Close)
}

func (conn *fasterTcpConn) GetUnderlyingConn() net.Conn {
	return conn.TCPConn
}

type fasterTcpConn struct {
	*net.TCPConn
	fd        *os.File
	closeLock sync.Mutex
}

func LessDelayTcpConn(conn *net.TCPConn) (connOut net.Conn, err error) {

	fd, err := conn.File()
	if err != nil {
		udwErr.LogErrorWithStack(err)
		return
	}
	conn1, err := net.FileConn(fd)
	if err != nil {
		fd.Close()
		udwErr.LogErrorWithStack(err)
		return
	}
	conn.Close()

	return &fasterTcpConn{TCPConn: conn1.(*net.TCPConn), fd: fd}, nil
}

func LessDelayDial(network string, address string) (conn net.Conn, err error) {
	conn, err = net.Dial(network, address)
	if err != nil {
		return conn, err
	}
	return LessDelayTcpConn(conn.(*net.TCPConn))
}

type fasterTcpListener struct {
	*net.TCPListener
}

func MustLessDelayListen(network string, address string) net.Listener {
	return fasterTcpListener{MustListen(network, address).(*net.TCPListener)}
}

func (l fasterTcpListener) Accept() (outC net.Conn, err error) {
	c, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return c, err
	}
	return LessDelayTcpConn(c)
}
