package udwNet

import (
	"errors"
	"io"
	"net"
	"time"
)

func NewRwcOverConn(rwc io.ReadWriteCloser, conn net.Conn) net.Conn {
	return &RwcOverConn{
		Rwc:    rwc,
		Reader: rwc,
		Writer: rwc,
		Closer: rwc,
		Conn:   conn,
	}
}

type RwcOverConn struct {
	Rwc io.ReadWriteCloser
	io.Reader
	io.Writer
	io.Closer
	net.Conn
}

func (c *RwcOverConn) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

func (c *RwcOverConn) Write(p []byte) (n int, err error) {
	return c.Writer.Write(p)
}
func (c *RwcOverConn) Close() (err error) {
	return c.Closer.Close()
}
func (c *RwcOverConn) GetUnderlyingConn() net.Conn {
	return c.Conn
}
func RwcConnWithSetDeadline(rwc io.ReadWriteCloser, setDeadline func(t time.Time) error) net.Conn {
	return rwcConn{
		ReadWriteCloser:   rwc,
		SetDeadlineMethod: setDeadline,
	}
}

func RwcConn(rwc io.ReadWriteCloser) net.Conn {
	return rwcConn{
		ReadWriteCloser: rwc,
		SetDeadlineMethod: func(t time.Time) error {
			return errors.New("udwNet.rwcConn does not support deadlines")
		},
	}
}

type rwcConn struct {
	io.ReadWriteCloser
	SetDeadlineMethod func(t time.Time) error
}

func (c rwcConn) LocalAddr() net.Addr {
	return FakeAddr
}

func (c rwcConn) RemoteAddr() net.Addr {
	return FakeAddr
}

func (c rwcConn) SetDeadline(t time.Time) error {
	if c.SetDeadlineMethod == nil {
		return errors.New("udwNet.rwcConn does not support deadlines")
	} else {
		return c.SetDeadlineMethod(t)
	}
}

func (c rwcConn) SetReadDeadline(t time.Time) error {
	return errors.New("udwNet.rwcConn does not support SetReadDeadline")
}

func (c rwcConn) SetWriteDeadline(t time.Time) error {
	return errors.New("udwNet.rwcConn does not support deadlines")
}
