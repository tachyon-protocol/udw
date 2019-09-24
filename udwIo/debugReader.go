package udwIo

import (
	"fmt"
	"io"
)

type debugReader struct {
	r    io.Reader
	name string
}

func NewDebugReader(r io.Reader, name string) *debugReader {
	return &debugReader{
		r:    r,
		name: name,
	}
}

func (r *debugReader) Read(b []byte) (n int, err error) {
	fmt.Println("start read", r.name, len(b))
	n, err = r.r.Read(b)
	fmt.Println("finish write", r.name, len(b), n, err)
	return n, err
}
