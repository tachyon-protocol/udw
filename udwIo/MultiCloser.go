package udwIo

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"io"
)

func MultiCloser(closers ...io.Closer) io.Closer {
	c := make([]io.Closer, len(closers))
	copy(c, closers)
	return multiCloser(c)
}

type multiCloser []io.Closer

func (c multiCloser) Close() (err error) {
	for _, closer := range c {
		err1 := closer.Close()
		if err1 != nil {
			err = err1
		}
	}
	return err
}

func MultiErrorHandle(fs ...func() error) error {
	return NewMultiErrorHandler(fs...)()
}

func NewMultiErrorHandler(fs ...func() error) func() error {
	return func() error {
		var errS string
		for _, f := range fs {
			err1 := f()
			if err1 != nil {
				errS += "[" + err1.Error() + "] "
			}
		}
		if errS == "" {
			return nil
		}
		return udwErr.New(errS)
	}
}

func ErrorHandlerNil() error {
	return nil
}
