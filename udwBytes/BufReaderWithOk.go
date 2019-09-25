package udwBytes

type BufReaderWithOk struct {
	br       BufReader
	hasError bool
}

func NewBufReaderWithOk(buf []byte) *BufReaderWithOk {
	return &BufReaderWithOk{
		br: BufReader{
			buf: buf,
		},
	}
}

func (buf *BufReaderWithOk) ReadBigEndUint16() (x uint16) {
	x, ok := buf.br.ReadBigEndUint16()
	if ok == false {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) ReadLittleEndUint16() (x uint16) {
	x, ok := buf.br.ReadLittleEndUint16()
	if ok == false {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) ReadLittleEndUint32() (x uint32) {
	x, ok := buf.br.ReadLittleEndUint32()
	if ok == false {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) ReadLittleEndUint64() (x uint64) {
	x, ok := buf.br.ReadLittleEndUint64()
	if ok == false {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) ReadUvarint() (x uint64) {
	x, ok := buf.br.ReadUvarint()
	if ok == false {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) GetRemainSize() int {
	return buf.br.GetRemainSize()
}

func (buf *BufReaderWithOk) ReadStringLenUvarint() string {
	x, ok := buf.br.ReadStringLenUvarint()
	if ok == false {
		buf.hasError = true
	}
	return x
}
func (buf *BufReaderWithOk) ReadSliceBySize(size int) []byte {
	x := buf.br.ReadMaxByteNum(size)
	if len(x) != size {
		buf.hasError = true
	}
	return x
}

func (buf *BufReaderWithOk) GetPos() int {
	return buf.br.GetPos()
}

func (buf *BufReaderWithOk) IsOk() bool {
	return buf.hasError == false
}

func (buf *BufReaderWithOk) IsEof() bool {
	return buf.br.IsEof()
}
