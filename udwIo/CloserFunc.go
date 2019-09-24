package udwIo

type CloserFunc func() (err error)

func (f CloserFunc) Close() (err error) {
	return f()
}
