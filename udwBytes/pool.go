package udwBytes

import (
	"sync"
)

type BufWriterPool struct {
	pool sync.Pool
}

func (pool *BufWriterPool) Get() *BufWriter {
	w := pool.pool.Get()
	if w == nil {

		return &BufWriter{}
	}
	w2, ok := w.(*BufWriter)
	if ok && w2 != nil {
		w2.Reset()
		return w2
	}

	return &BufWriter{}
}

func (pool *BufWriterPool) Put(w *BufWriter) {
	if w != nil {
		pool.pool.Put(w)
	}
}

func (pool *BufWriterPool) GetAndCloneFromByteSlice(buf []byte) *BufWriter {
	bw := pool.Get()
	bw.Write(buf)
	return bw
}
