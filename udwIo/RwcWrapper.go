package udwIo

import (
	"io"
)

type RwcWrapper interface {
	RwcWrap(in io.ReadWriteCloser) (out io.ReadWriteCloser, err error)
}

type RwcWrapperFunc func(in io.ReadWriteCloser) (out io.ReadWriteCloser, err error)

func (f RwcWrapperFunc) RwcWrap(in io.ReadWriteCloser) (out io.ReadWriteCloser, err error) {
	return f(in)
}
