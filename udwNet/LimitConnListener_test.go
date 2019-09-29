package udwNet

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestNewLimitConnListener(t *testing.T) {
	tcpLn, err := net.Listen("tcp", "127.0.0.1:30004")
	udwErr.PanicIfError(err)
	ln := NewLimitConnTcpListener(LimitConnListenerRequest{
		OriginLn:      tcpLn.(*net.TCPListener),
		MaxConnNumber: 2,
	})
	inNumber := &udwSync.Int{}
	closer := TcpNewListenerFromExistListener(ln, func(conn net.Conn) {
		inNumber.Add(1)
		buf := make([]byte, 1)
		conn.Read(buf)
		conn.Close()
	})
	defer closer()
	newCConnFn := func() net.Conn {
		cc, err := net.Dial("tcp", "127.0.0.1:30004")
		udwErr.PanicIfError(err)
		return cc
	}
	c1 := newCConnFn()
	defer c1.Close()
	c2 := newCConnFn()
	defer c2.Close()
	c3 := newCConnFn()
	defer c3.Close()
	udwTest.Equal(inNumber.Get(), 2)
	c1.Close()
	c2.Close()
	c4 := newCConnFn()
	defer c4.Close()
	c5 := newCConnFn()
	defer c5.Close()
	udwTest.Equal(inNumber.Get(), 4)
}
