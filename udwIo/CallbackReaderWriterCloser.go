package udwIo

type CallbackReaderWriterCloser struct {
	Reader func(b []byte) (n int, err error)
	Writer func(b []byte) (n int, err error)
	Closer func() (err error)
}

func (crwc CallbackReaderWriterCloser) Read(b []byte) (n int, err error) {
	return crwc.Reader(b)
}

func (crwc CallbackReaderWriterCloser) Write(b []byte) (n int, err error) {
	return crwc.Writer(b)
}

func (crwc CallbackReaderWriterCloser) Close() (err error) {
	return crwc.Closer()
}
