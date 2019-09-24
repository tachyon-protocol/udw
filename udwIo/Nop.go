package udwIo

import (
	"io"
	"io/ioutil"
)

var NopCloser io.Closer = Nop

var Nop = tnop{}

type tnop struct{}

func (c tnop) Close() (err error) {
	return nil
}
func (c tnop) Write(b []byte) (n int, err error) {
	return len(b), nil
}
func (c tnop) Read(b []byte) (n int, err error) {
	return len(b), nil
}

func NewReaderWithNopCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func NewRwcFromReader(r io.Reader) io.ReadWriteCloser {
	return &StructWriterReaderCloser{
		Reader: r,
		Writer: Nop,
		Closer: Nop,
	}
}
