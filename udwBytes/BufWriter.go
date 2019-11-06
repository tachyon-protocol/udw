package udwBytes

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type BufWriter struct {
	buf []byte
}

func NewBufWriter(buf []byte) *BufWriter {
	return &BufWriter{
		buf: buf,
	}
}

func NewBufWriterString(s string) *BufWriter {
	return &BufWriter{
		buf: []byte(s),
	}
}

func ResetBufWriter(inW *BufWriter) *BufWriter {
	if inW == nil {
		return &BufWriter{}
	}
	inW.Reset()
	return inW
}

func (w *BufWriter) Write_(p []byte) {
	w.buf = append(w.buf, p...)
}

func (w *BufWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)

	return len(p), nil
}

func (w *BufWriter) WriteString_(s string) {
	w.buf = append(w.buf, s...)
}

func (w *BufWriter) WriteString(p string) (n int, err error) {
	w.buf = append(w.buf, p...)

	return len(p), nil
}

const minRead = 512

func (w *BufWriter) ReadFrom(r io.Reader) (n int64, err error) {
	for {

		w.TryGrow(minRead)
		m, e := r.Read(w.buf[len(w.buf):cap(w.buf)])

		w.buf = w.buf[:len(w.buf)+m]
		n += int64(m)
		if e == io.EOF {
			break
		}
		if e != nil {
			return n, e
		}
	}
	return n, nil
}

func (w *BufWriter) WriteByte_(b uint8) {
	w.buf = append(w.buf, b)
}

func (w *BufWriter) WriteByte(b uint8) error {
	w.buf = append(w.buf, b)
	return nil
}

func (w *BufWriter) WriteBigEndUint16(v uint16) {
	w.TryGrow(2)
	binary.BigEndian.PutUint16(w.buf[len(w.buf):len(w.buf)+2], v)
	w.buf = w.buf[:len(w.buf)+2]
}
func (w *BufWriter) WriteBigEndUint32(v uint32) {
	w.TryGrow(4)
	binary.BigEndian.PutUint32(w.buf[len(w.buf):len(w.buf)+4], v)
	w.buf = w.buf[:len(w.buf)+4]
}
func (w *BufWriter) WriteBigEndUint64(v uint64) {
	w.TryGrow(8)
	binary.BigEndian.PutUint64(w.buf[len(w.buf):len(w.buf)+8], v)
	w.buf = w.buf[:len(w.buf)+8]
}
func (w *BufWriter) WriteLittleEndUint16(v uint16) {
	w.TryGrow(2)
	binary.LittleEndian.PutUint16(w.buf[len(w.buf):len(w.buf)+2], v)
	w.buf = w.buf[:len(w.buf)+2]
}
func (w *BufWriter) WriteLittleEndUint32(v uint32) {
	w.TryGrow(4)
	binary.LittleEndian.PutUint32(w.buf[len(w.buf):len(w.buf)+4], v)
	w.buf = w.buf[:len(w.buf)+4]
}
func (w *BufWriter) WriteLittleEndUint64(v uint64) {
	w.TryGrow(8)
	binary.LittleEndian.PutUint64(w.buf[len(w.buf):len(w.buf)+8], v)
	w.buf = w.buf[:len(w.buf)+8]
}

func (w *BufWriter) WriteLittleEndFloat64(f float64) {
	w.WriteLittleEndUint64(math.Float64bits(f))
}
func (w *BufWriter) WriteLittleEndFloat32(f float32) {
	w.WriteLittleEndUint32(math.Float32bits(f))
}
func (w *BufWriter) MustWriteString255(s string) {
	if len(s) > 255 {
		panic(fmt.Errorf("[MustWriteString255] len(s)[%d]>255", len(s)))
	}

	w.buf = append(w.buf, byte(len(s)))
	w.buf = append(w.buf, s...)

}

func (w *BufWriter) WriteBool(b bool) {
	b1 := byte(0)
	if b == true {
		b1 = 1
	}
	w.WriteByte_(b1)
}

func (w *BufWriter) WriteByteBySize(size int, b byte) {
	w.TryGrow(size)
	startPos := len(w.buf)
	w.buf = w.buf[:len(w.buf)+size]
	for i := startPos; i < startPos+size; i++ {
		w.buf[i] = b
	}
}

func (w *BufWriter) WriteZeroByteBySize(size int) {
	w.WriteByteBySize(size, 0)
}

func (w *BufWriter) ResetWithBuffer(buf []byte) {
	if buf != nil {
		w.buf = buf[:0]
	} else {
		w.buf = nil
	}
}
func (w *BufWriter) Reset() {
	w.buf = w.buf[:0]
}

func (w *BufWriter) GetBytesClone() []byte {
	return Clone(w.GetBytes())
}

func (w *BufWriter) GetBytes() []byte {
	return w.buf
}

func (w *BufWriter) GetString() string {
	return string(w.buf)
}

func (w *BufWriter) GetHeadBuffer(size int) []byte {
	w.tryGrowWithSetPos(size)
	return w.buf[len(w.buf) : len(w.buf)+size]
}

func (w *BufWriter) AddPos(offset int) {
	w.tryGrowWithSetPos(offset)

	w.buf = w.buf[:len(w.buf)+offset]
}

func (w *BufWriter) SetPos(pos int) {
	w.tryGrowWithSetPos(pos - len(w.buf))
	w.buf = w.buf[:pos]
}

func (w *BufWriter) GetPos() int {
	return len(w.buf)
}
func (w *BufWriter) GetLen() int {
	return len(w.buf)
}
func (w *BufWriter) GetCap() int {
	return cap(w.buf)
}

func (w *BufWriter) GetAllocNoUseBuf() []byte {
	return w.buf[len(w.buf):cap(w.buf)]
}

func (w BufWriter) GetLastByte() (b byte, ok bool) {
	if len(w.buf) == 0 {
		return 0, false
	}
	return w.buf[len(w.buf)-1], true
}

func (w *BufWriter) TryGrow(toWrite int) {
	needSize := len(w.buf) + toWrite
	thisCap := cap(w.buf)
	if needSize > thisCap {
		targetSize := thisCap*2 + toWrite
		if targetSize < 32 {
			targetSize = 32
		}
		oldPos := len(w.buf)
		newBuf := make([]byte, oldPos, targetSize)
		copy(newBuf, w.buf)
		w.buf = newBuf
	}
}

func (w *BufWriter) tryGrowWithSetPos(toWrite int) {
	needSize := len(w.buf) + toWrite
	thisCap := cap(w.buf)
	if needSize > thisCap {
		targetSize := thisCap*2 + toWrite
		if targetSize < 32 {
			targetSize = 32
		}
		oldPos := len(w.buf)

		newBuf := make([]byte, targetSize, targetSize)
		copy(newBuf, w.buf[:cap(w.buf)])
		w.buf = newBuf[:oldPos]
	}
}
