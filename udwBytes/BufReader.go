package udwBytes

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BufReader struct {
	buf []byte
	pos int
}

func NewBufReader(buf []byte) *BufReader {
	return &BufReader{
		buf: buf,
	}
}

func (r *BufReader) ResetWithBuffer(buf []byte) {
	r.pos = 0
	r.buf = buf
}

func (r *BufReader) Read(inBuf []byte) (n int, err error) {
	remainSize := len(r.buf) - r.pos
	if remainSize > len(inBuf) {
		copy(inBuf, r.buf[r.pos:])
		r.pos += len(inBuf)
		return len(inBuf), nil
	} else {
		copy(inBuf, r.buf[r.pos:])
		r.pos += remainSize
		return remainSize, io.EOF
	}
}

func (r *BufReader) ReadByte() (b byte, err error) {
	if r.pos >= len(r.buf) {
		return 0, io.EOF
	}
	b = r.buf[r.pos]
	r.pos += 1
	return b, nil
}

func (r *BufReader) MustReadByte() (b byte) {
	b = r.buf[r.pos]
	r.pos += 1
	return b
}

func (r *BufReader) ReadByteOrEof() (b byte, isRead bool) {
	if r.pos >= len(r.buf) {
		return 0, false
	}
	b = r.buf[r.pos]
	r.pos += 1
	return b, true
}

func (r *BufReader) ReadAt(inBuf []byte, off int64) (n int, err error) {
	remainSize := len(r.buf) - int(off)
	if remainSize > len(inBuf) {
		copy(inBuf, r.buf[r.pos:])
		return len(inBuf), nil
	} else {
		copy(inBuf, r.buf[off:])
		return remainSize, io.EOF
	}
}

func (r *BufReader) IsEof() bool {
	return r.pos >= len(r.buf)
}

func (r *BufReader) GetRemainSize() int {
	return len(r.buf) - r.pos
}

func (r *BufReader) GetPos() int {
	return r.pos
}

func (r *BufReader) SetPos(pos int) {
	r.pos = pos
}

func (r *BufReader) AddPos(rel int) {
	r.pos += rel
}
func (r *BufReader) GetBuf() []byte {
	return r.buf
}

func (r *BufReader) ReadMaxByteNum(num int) []byte {
	startPos := r.pos
	if startPos >= len(r.buf) {
		return nil
	} else if r.pos+num < len(r.buf) {
		r.pos += num
	} else {
		r.pos = len(r.buf)
	}
	return r.buf[startPos:r.pos]
}

func (r *BufReader) ReadHttpHeaderWithCallback(fn func(name []byte, value []byte)) (errMsg string) {
	for {
		if r.IsEof() {

			return ""
		}
		thisLine := r.ReadToLineEnd()
		thisLine = bytes.TrimSpace(thisLine)
		if len(thisLine) == 0 {

			return ""
		}

		headerName, value := SplitTwoBetweenFirst(thisLine, []byte(": "))
		if len(headerName) == 0 || len(value) == 0 {

			return "7yyaavsecr"
		}
		fn(headerName, value)
	}
}

func (r *BufReader) ReadString255() (s string, errMsg string) {
	if r.pos >= len(r.buf) {
		return "", "emkdt38k48 unexpected EOF"
	}
	size := r.buf[r.pos]
	r.pos += 1
	if r.pos+int(size) > len(r.buf) {
		return "", "7srvg667mv"
	}
	s = string(r.buf[r.pos : r.pos+int(size)])
	r.pos += int(size)
	return s, ""
}

func (r *BufReader) ReadCString() (s string) {
	startPos := r.pos
	for {
		if r.pos >= len(r.buf) {
			return string(r.buf[startPos:])
		}
		b := r.buf[r.pos]
		r.pos++
		if b == 0 {
			break
		}
	}
	return string(r.buf[startPos : r.pos-1])
}

func (r *BufReader) Close() (err error) {
	return nil
}

func (r *BufReader) ReadBigEndUint16() (x uint16, isOk bool) {
	buf := r.ReadMaxByteNum(2)
	if len(buf) != 2 {
		return 0, false
	}
	x = binary.BigEndian.Uint16(buf)
	return x, true
}

func (r *BufReader) ReadBigEndUint32() (x uint32, isOk bool) {
	buf := r.ReadMaxByteNum(4)
	if len(buf) != 4 {
		return 0, false
	}
	x = binary.BigEndian.Uint32(buf)
	return x, true
}

func (r *BufReader) ReadBigEndUint64() (x uint64, isOk bool) {
	buf := r.ReadMaxByteNum(8)
	if len(buf) != 8 {
		return 0, false
	}
	x = binary.BigEndian.Uint64(buf)
	return x, true
}

func (r *BufReader) ReadLittleEndUint16() (x uint16, isOk bool) {
	buf := r.ReadMaxByteNum(2)
	if len(buf) != 2 {
		return 0, false
	}
	x = binary.LittleEndian.Uint16(buf)
	return x, true
}

func (r *BufReader) ReadLittleEndUint32() (x uint32, isOk bool) {
	buf := r.ReadMaxByteNum(4)
	if len(buf) != 4 {
		return 0, false
	}
	x = binary.LittleEndian.Uint32(buf)
	return x, true
}

func (r *BufReader) ReadLittleEndUint64() (x uint64, isOk bool) {
	buf := r.ReadMaxByteNum(8)
	if len(buf) != 8 {
		return 0, false
	}
	x = binary.LittleEndian.Uint64(buf)
	return x, true
}

func (r *BufReader) ReadStringLenUint32() (s string, isOk bool) {
	sLen, ok := r.ReadLittleEndUint32()
	if !ok {
		return s, false
	}
	if r.pos+int(sLen) > len(r.buf) {
		return s, false
	}
	s = string(r.buf[r.pos : r.pos+int(sLen)])
	r.pos += int(sLen)
	return s, true
}

func (r *BufReader) PeekByte() (b byte, ok bool) {
	if len(r.buf) <= r.pos {
		return 0, false
	}
	return r.buf[r.pos], true
}

func (r *BufReader) ReadByteSlice(num int) (s []byte, ok bool) {
	end := r.pos + num
	if num < 0 || end > len(r.buf) {
		return nil, false
	}
	s = r.buf[r.pos:end]
	r.pos = end
	return s, true
}

func (r *BufReader) MustReadByteSlice(num int) (s []byte) {
	end := r.pos + num
	if end > len(r.buf) {
		panic(`t2ga4b2n9x`)
	}
	s = r.buf[r.pos:end]
	r.pos = end
	return
}
func (r *BufReader) MustReadLittleEndUint32() (x uint32) {
	buf := r.ReadMaxByteNum(4)
	if len(buf) != 4 {
		panic("bezxbb5n3b")
	}
	x = binary.LittleEndian.Uint32(buf)
	return x
}
