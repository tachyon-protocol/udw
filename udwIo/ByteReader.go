package udwIo

import (
	"io"
)

type SingleByteReader interface {
	io.Reader
	io.ByteReader
}

const maxConsecutiveEmptyReads = 100

func NewSingleByteReader(r io.Reader) SingleByteReader {
	return singleByteReader{Reader: r}
}

type singleByteReader struct {
	buf [1]byte
	io.Reader
}

func (r singleByteReader) ReadByte() (c byte, err error) {

	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := r.Reader.Read(r.buf[:])
		if err != nil {
			return 0, err
		}
		if n > 0 {
			return r.buf[0], nil
		}
	}
	return 0, io.ErrNoProgress
}
