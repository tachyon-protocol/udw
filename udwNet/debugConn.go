package udwNet

import (
	"fmt"
	"io"
	"net"

	"github.com/tachyon-protocol/udw/udwIo"
	"github.com/tachyon-protocol/udw/udwRandNewId"
)

type connRwcer struct {
	net.Conn
	rwc io.ReadWriteCloser
}

func (conn connRwcer) Read(p []byte) (n int, err error) {
	return conn.rwc.Read(p)
}
func (conn connRwcer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {

		panic("[Write] len(p)==0")
	}
	return conn.rwc.Write(p)
}
func (conn connRwcer) Close() (err error) {
	return conn.rwc.Close()
}
func (conn connRwcer) GetUnderlyingConn() net.Conn {
	return conn.Conn
}

func NewDebugConn(conn net.Conn, name string) net.Conn {
	return connRwcer{
		Conn: conn,
		rwc:  udwIo.NewDebugRwc(conn, name+"["+udwRandNewId.NextUint64IdString()+"]["+conn.LocalAddr().String()+"-"+conn.RemoteAddr().String()+"]"),
	}
}

func NewDebugConnNoData(conn net.Conn) net.Conn {
	return connRwcer{
		Conn: conn,
		rwc:  udwIo.NewDebugRwcNoData(conn, "["+udwRandNewId.NextUint64IdString()+"]["+conn.LocalAddr().String()+"-"+conn.RemoteAddr().String()+"]"),
	}
}

func NewDebugConnNoDataWithName(conn net.Conn, name string) net.Conn {
	return connRwcer{
		Conn: conn,
		rwc:  udwIo.NewDebugRwcNoData(conn, name+"["+udwRandNewId.NextUint64IdString()+"]"),
	}
}

func NewDebugDialerNoData(parent Dialer) Dialer {
	return func(network, address string) (net.Conn, error) {
		conn, err := parent(network, address)
		if err != nil {
			return nil, err
		}
		return NewDebugConnNoData(conn), nil
	}
}

type sizeConnRwcer struct {
	net.Conn
	name string
}

func (c sizeConnRwcer) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Write len:", n, "err:", err)
	} else {
		fmt.Println("[debugConn]", c.name, "Write len:", n)
	}
	return n, err
}
func (c sizeConnRwcer) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b), "err:", err)
	} else {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b))
	}
	return n, err
}
func (c sizeConnRwcer) Close() (err error) {
	err = c.Conn.Close()
	fmt.Println("[debugConn]", c.name, "Close err:", err)
	return err
}

func NewSizeDebugConn(conn net.Conn, name string) net.Conn {
	return sizeConnRwcer{
		Conn: conn,
		name: name,
	}
}

type stringConnRwcer struct {
	net.Conn
	name     string
	readNum  int
	writeNum int
}

func (c *stringConnRwcer) Write(b []byte) (n int, err error) {
	c.writeNum++
	n, err = c.Conn.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Write len:", n, "err:", err)
	} else {
		fmt.Printf("[debugConn] [%s] Write %d len: %d content: %q<EndOfContent>\n", c.name, c.writeNum, n, string(b[:n]))
	}
	return n, err
}
func (c *stringConnRwcer) Read(b []byte) (n int, err error) {
	c.readNum++
	n, err = c.Conn.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.name, "Read len:", n, "iLen", len(b), "err:", err)
	} else {
		fmt.Printf("[debugConn] [%s] Read %d len: %d iLen: %d content: %q<EndOfContent>\n", c.name, c.readNum, n, len(b), string(b[:n]))
	}
	return n, err
}
func (c *stringConnRwcer) Close() (err error) {
	err = c.Conn.Close()
	fmt.Println("[debugConn]", c.name, "Close err:", err)
	return err
}

func NewStringDebugConn(conn net.Conn, name string) net.Conn {
	return &stringConnRwcer{
		Conn: conn,
		name: name,
	}
}
