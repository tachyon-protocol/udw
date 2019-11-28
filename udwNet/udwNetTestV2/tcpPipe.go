package udwNetTestV2

import (
	"net"
	"sync"
)

func TcpPipe() (c1 *net.TCPConn, c2 *net.TCPConn, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var c2Err error
	go func() {
		var c2_1 net.Conn
		c2_1, c2Err = ln.Accept()
		if c2Err != nil {
			wg.Done()
			return
		}
		c2 = c2_1.(*net.TCPConn)
		c2Err = c2.SetLinger(0)
		if c2Err != nil {
			c2_1.Close()
			wg.Done()
			return
		}
		c2Err = ln.Close()
		if c2Err != nil {
			c2_1.Close()
			wg.Done()
			return
		}
		wg.Done()
	}()
	var c1_1 net.Conn
	c1_1, err = net.Dial("tcp", ln.Addr().String())
	if err != nil {
		ln.Close()
		return nil, nil, err
	}
	wg.Wait()
	if c2Err != nil {
		c1_1.Close()
		return nil, nil, c2Err
	}
	c1 = c1_1.(*net.TCPConn)

	err = c1.SetLinger(0)
	if err != nil {
		c1.Close()
		return nil, nil, err
	}
	return c1, c2, nil
}

func MustTcpPipe() (c1 net.Conn, c2 net.Conn) {
	c1, c2, err := TcpPipe()
	if err != nil {
		panic(err)
	}
	return c1, c2
}
