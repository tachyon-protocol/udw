package udwIo

import "io"

type LimitReader struct {
	r     io.Reader
	n     int64
	isEof bool
}

func NewLimitReader(r io.Reader, n int64) *LimitReader {
	return &LimitReader{r: r, n: n}
}

func (l *LimitReader) Read(p []byte) (n int, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.n {
		p = p[0:l.n]
	}
	n, err = l.r.Read(p)
	l.n -= int64(n)
	if err == io.EOF {
		l.isEof = true
	}
	return n, err
}

func (l *LimitReader) IsUnderReaderEof() bool {
	return l.isEof
}
