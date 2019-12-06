package udwBufio

import (
	"errors"
	"io"
)

type BufioReader struct {
	r           io.Reader
	maxNeedSize int
	bufW        []byte

	noReadSp int

	noReadEp int
}

func NewBufioReader(r io.Reader, maxNeedSize int) *BufioReader {
	if maxNeedSize == 0 {
		maxNeedSize = 32 * 1024 * 1024
	}
	return &BufioReader{
		r:           r,
		maxNeedSize: maxNeedSize,
		noReadSp:    0,
		noReadEp:    0,
	}
}

func (crb *BufioReader) ReadBySize(need int) (b []byte, errMsg string) {
	err := crb.fill(need)
	if err != nil {
		return nil, err.Error()
	}
	b = crb.bufW[crb.noReadSp : crb.noReadSp+need]
	crb.noReadSp += need
	return b, ""
}

func (crb *BufioReader) Read(buf []byte) (nr int, err error) {
	n := crb.noReadEp - crb.noReadSp
	if n > 0 {
		bufSize := len(buf)
		if n >= bufSize {
			copy(buf, crb.bufW[crb.noReadSp:crb.noReadSp+bufSize])
			crb.noReadSp += bufSize
			return bufSize, nil
		} else {
			copy(buf, crb.bufW[crb.noReadSp:crb.noReadEp])
			crb.bufW = crb.bufW[:0]
			crb.noReadSp = 0
			crb.noReadEp = 0
			return n, nil
		}
	}
	return crb.r.Read(buf)
}

func (crb *BufioReader) AddPos(pos int) {
	afterValue := crb.noReadSp + pos
	if afterValue < 0 {

		panic("z8xgzab2qr")
	}
	crb.noReadSp = afterValue
}

func (crb *BufioReader) fill(need int) (err error) {
	if need < 0 {
		return errors.New("xtf3vsxsxx")
	}
	n := crb.noReadEp - crb.noReadSp
	if n >= need {

		return nil
	}
	if need > crb.maxNeedSize {
		return errors.New("jbfy8v9fzm")
	}

	if n == 0 {
		crb.bufW = crb.bufW[:0]
		crb.noReadSp = 0
		crb.noReadEp = 0
	} else if crb.noReadSp != 0 {

		copy(crb.bufW[0:n], crb.bufW[crb.noReadSp:crb.noReadEp])
		crb.bufW = crb.bufW[0:n]
		crb.noReadSp = 0
		crb.noReadEp = n
	}
	needSize := len(crb.bufW) + need
	thisCap := cap(crb.bufW)
	if needSize > thisCap {
		targetSize := thisCap*2 + need
		if targetSize < 4*1024 {
			targetSize = 4 * 1024
		}
		newBuf := make([]byte, len(crb.bufW), targetSize)
		copy(newBuf, crb.bufW)
		crb.bufW = newBuf
	}
	for {
		buf := crb.bufW[len(crb.bufW):cap(crb.bufW)]
		nn, err := crb.r.Read(buf)
		n += nn
		crb.bufW = crb.bufW[:len(crb.bufW)+nn]
		crb.noReadEp += nn
		if err == nil {
			if n < need {
				continue
			}
			return nil
		} else if err == io.EOF {
			if n >= need {
				return nil
			}
			return io.ErrUnexpectedEOF
		} else {
			return err
		}
	}
}

func (crb *BufioReader) ReadByteErrMsg() (b byte, errMsg string) {
	bs, errMsg := crb.ReadBySize(1)
	if errMsg != "" {
		return 0, errMsg
	}
	return bs[0], ""
}

func (crb *BufioReader) ReadUvarint() (x uint64, errMsg string) {
	var s uint
	i := 0
	for {
		b, errMsg := crb.ReadByteErrMsg()
		if errMsg != "" {
			return 0, errMsg
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, "2vg5peydjt overflow"
			}
			return x | uint64(b)<<s, ""
		}
		x |= uint64(b&0x7f) << s
		s += 7
		i++
	}
}

func (crb *BufioReader) PeekByte() (b byte, errMsg string) {
	bs, errMsg := crb.ReadBySize(1)
	if errMsg != "" {
		return 0, errMsg
	}
	crb.AddPos(-1)
	return bs[0], ""
}
