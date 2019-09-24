package udwNet

import (
	"net"
	"time"
)

func NewTimeoutDialer(timeout time.Duration) func(network, addr string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}
		return &timeoutConn{
			Conn:    conn,
			Timeout: timeout,
		}, nil
	}
}

func TimeoutConn(conn net.Conn, timeout time.Duration) net.Conn {
	return &timeoutConn{
		Conn:    conn,
		Timeout: timeout,
	}
}

type timeoutConn struct {
	net.Conn
	Timeout time.Duration
}

func (c *timeoutConn) Read(b []byte) (n int, err error) {
	err = c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *timeoutConn) Write(b []byte) (n int, err error) {
	err = c.Conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}
