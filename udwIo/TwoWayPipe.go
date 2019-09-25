package udwIo

import (
	"io"
)

func TwoWayPipe() (rwc1 io.ReadWriteCloser, rwc2 io.ReadWriteCloser) {
	r1, w2 := io.Pipe()
	r2, w1 := io.Pipe()
	closer := MultiCloser(r1, r2)
	return StructWriterReaderCloser{
			Reader: r1,
			Writer: w1,
			Closer: closer,
		}, StructWriterReaderCloser{
			Reader: r2,
			Writer: w2,
			Closer: closer,
		}
}
