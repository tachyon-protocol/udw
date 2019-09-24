package udwIo

import (
	"fmt"
	"io"
)

type debugWriter struct {
	w    io.Writer
	name string
}

func NewDebugWriter(w io.Writer, name string) *debugWriter {
	return &debugWriter{
		w:    w,
		name: name,
	}
}

func (r *debugWriter) Write(b []byte) (n int, err error) {
	fmt.Println("start write", r.name, len(b))
	n, err = r.w.Write(b)
	fmt.Println("finish write", r.name, len(b), n, err)
	return n, err
}
