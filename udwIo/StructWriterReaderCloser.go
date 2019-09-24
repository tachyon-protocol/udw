package udwIo

import (
	"io"
)

type StructWriterReaderCloser struct {
	io.Writer
	io.Reader
	io.Closer
}
