package udwBufio

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
)

type BufioWriter struct {
	tmpBuf         udwBytes.BufWriter
	w              io.Writer
	softMaxBufSize int
}

func NewBufioWriter(w io.Writer, softMaxBufSize int) *BufioWriter {
	if softMaxBufSize == 0 {
		softMaxBufSize = 32 * 1024
	}
	return &BufioWriter{
		w:              w,
		softMaxBufSize: softMaxBufSize,
	}
}

func (w *BufioWriter) SoftFlush() (errMsg string) {
	if w.tmpBuf.GetLen() >= w.softMaxBufSize {
		return w.Flush()
	}
	return
}

func (w *BufioWriter) Flush() (errMsg string) {
	if w.tmpBuf.GetLen() == 0 {
		return ""
	}
	_, err := w.w.Write(w.tmpBuf.GetBytes())
	if err != nil {
		return err.Error()
	}
	w.tmpBuf.Reset()
	return ""
}

func (w *BufioWriter) WriteUvarint(x uint64) {
	w.tmpBuf.WriteUvarint(x)
}

func (w *BufioWriter) Write_(buf []byte) {
	w.tmpBuf.Write_(buf)
}

func (w *BufioWriter) WriteByte_(buf byte) {
	w.tmpBuf.WriteByte_(buf)
}

func (w *BufioWriter) WriteLittleEndFloat64(x float64) {
	w.tmpBuf.WriteLittleEndFloat64(x)
}

func (w *BufioWriter) WriteLittleEndFloat32(x float32) {
	w.tmpBuf.WriteLittleEndFloat32(x)
}

func (w *BufioWriter) AddPos(offset int) {
	w.tmpBuf.AddPos(offset)
}

func (w *BufioWriter) GetPos() int {
	return w.tmpBuf.GetPos()
}
func (w *BufioWriter) SetPos(pos int) {
	w.tmpBuf.SetPos(pos)
}
