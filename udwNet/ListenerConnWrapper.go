package udwNet

import "net"

type ConnWrapper func(conn net.Conn) (net.Conn, error)

func ListenerConnWrapper(l net.Listener, connWrapper ConnWrapper) net.Listener {
	return &listenerConnWrapper{
		Listener:    l,
		connWrapper: connWrapper,
	}
}

type listenerConnWrapper struct {
	net.Listener
	connWrapper ConnWrapper
}

func (l listenerConnWrapper) Accept() (c net.Conn, err error) {
	c1, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	c2, err := l.connWrapper(c1)
	if err != nil {
		c1.Close()
		return nil, err
	}
	return c2, nil
}
