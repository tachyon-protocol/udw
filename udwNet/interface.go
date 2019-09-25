package udwNet

import (
	"io"
)

type RwcDialer interface {
	RwcDial(addr string) (rwc io.ReadWriteCloser, err error)
}

type RwcDialerFunc func(addr string) (rwc io.ReadWriteCloser, err error)

func (f RwcDialerFunc) RwcDial(addr string) (rwc io.ReadWriteCloser, err error) {
	return f(addr)
}
