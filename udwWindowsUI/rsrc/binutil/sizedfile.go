package binutil

import (
	"bytes"
	"io"
)

type SizedReader interface {
	io.Reader
	Size() int64
}

type SizedFile struct {
	buf  *bytes.Buffer
	size int64
}

func (r *SizedFile) Read(p []byte) (n int, err error) { return r.buf.Read(p) }
func (r *SizedFile) Size() int64                      { return r.size }
func (r *SizedFile) Close() error                     { return nil }

func SizeFileFromBuffer(buf []byte) *SizedFile {
	return &SizedFile{
		buf:  bytes.NewBuffer(buf),
		size: int64(len(buf)),
	}
}
