package udwIo

import (
	"io"
	"io/ioutil"
)

func MustReadAll(r io.Reader) (b []byte) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

func ReadAll(r io.Reader) (b []byte, err error) {
	return ioutil.ReadAll(r)
}

var DiscardReadFrom = ioutil.Discard.(io.ReaderFrom).ReadFrom

func NewWriteAddCloser(w io.Writer, closer io.Closer) io.WriteCloser {
	return &writeCloser{
		Writer: w,
		Closer: closer,
	}
}

type writeCloser struct {
	io.Writer
	io.Closer
}

func NewReadAddCloser(r io.Reader, closer io.Closer) io.ReadCloser {
	return &readCloser{
		Reader: r,
		Closer: closer,
	}
}

type readCloser struct {
	io.Reader
	io.Closer
}
