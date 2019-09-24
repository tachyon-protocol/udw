package udwIo

import (
	"io"
	"sync"
)

type LockWriter struct {
	W      io.Writer
	locker sync.Mutex
}

func (w *LockWriter) Write(p []byte) (n int, err error) {
	w.locker.Lock()
	n, err = w.W.Write(p)
	w.locker.Unlock()
	return n, err
}
func NewLockWriter(w io.Writer) io.Writer {
	return &LockWriter{
		W: w,
	}
}
