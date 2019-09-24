package udwNet

import (
	"bytes"
	"errors"
	"net"
	"time"
)

func BufioConn() net.Conn {
	return &bufioConn{
		Buffer: &bytes.Buffer{},
	}
}

type bufioConn struct {
	*bytes.Buffer
	hasClose bool
}

func (c *bufioConn) Read(p []byte) (n int, err error) {
	if c.hasClose {
		return 0, GetSocketCloseError()
	}
	return c.Buffer.Read(p)
}

func (c *bufioConn) Write(p []byte) (n int, err error) {
	if c.hasClose {
		return 0, GetSocketCloseError()
	}
	return c.Buffer.Write(p)
}

func (c *bufioConn) Close() error {
	if c.hasClose {
		return GetSocketCloseError()
	}
	c.hasClose = true
	return nil
}

func (c bufioConn) LocalAddr() net.Addr {
	return FakeAddr
}

func (c bufioConn) RemoteAddr() net.Addr {
	return FakeAddr
}

func (c bufioConn) SetDeadline(t time.Time) error {
	return errors.New("udwNet.BufioConn does not support deadlines")
}

func (c bufioConn) SetReadDeadline(t time.Time) error {
	return errors.New("udwNet.BufioConn does not support deadlines")
}

func (c bufioConn) SetWriteDeadline(t time.Time) error {
	return errors.New("udwNet.BufioConn does not support deadlines")
}

var FakeAddr = fakeAddr{}

type fakeAddr struct{}

func (a fakeAddr) Network() string {
	return "fakeAddr"
}

func (a fakeAddr) String() string {
	return "fakeAddr"
}
