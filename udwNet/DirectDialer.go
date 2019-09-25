package udwNet

import (
	"net"
)

type DirectDialer interface {
	DirectDial() (net.Conn, error)
}

type DirectDialerFunc func() (net.Conn, error)

func (f DirectDialerFunc) Dial() (net.Conn, error) {
	return f()
}

type Dialer func(network, address string) (net.Conn, error)

func NewFixedDialerV2(parent Dialer, network string, address string) func(network, address string) (net.Conn, error) {
	return func(network1, address1 string) (net.Conn, error) {
		return parent(network, address)
	}
}

type FixedAddressDialer func() (net.Conn, error)

func NewFixedAddressDialer(parent Dialer, network string, address string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return parent(network, address)
	}
}
