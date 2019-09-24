package udwBytes

import (
	"encoding/binary"
	"io"
)

const readBufferWrapInitBufferSize = 4 * 1024

type ReadBufferWrap struct {
	R    io.Reader
	bufW []byte

	noReadSp int

	noReadEp int
}

func (crb *ReadBufferWrap) ReadBySize(need int) (b []byte, errMsg string) {
	n := crb.noReadEp - crb.noReadSp
	if n >= need {
		b = crb.bufW[crb.noReadSp : crb.noReadSp+need]
		crb.noReadSp += need
		return b, ""
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
		if targetSize < readBufferWrapInitBufferSize {
			targetSize = readBufferWrapInitBufferSize
		}
		newBuf := make([]byte, len(crb.bufW), targetSize)
		copy(newBuf, crb.bufW)
		crb.bufW = newBuf
	}
	for {
		buf := crb.bufW[len(crb.bufW):cap(crb.bufW)]
		nn, err := crb.R.Read(buf)
		n += nn
		crb.bufW = crb.bufW[:len(crb.bufW)+nn]
		crb.noReadEp += nn
		if err == nil {
			if n < need {
				continue
			}
			b = crb.bufW[crb.noReadSp : crb.noReadSp+need]
			crb.noReadSp += need
			return b, ""
		} else if err == io.EOF {
			if n >= need {
				b = crb.bufW[crb.noReadSp : crb.noReadSp+need]
				crb.noReadSp += need
				return b, ""
			}
			return nil, "unexpected EOF"
		} else {
			return nil, err.Error()
		}
	}
}

func (crb *ReadBufferWrap) ReadByteErrMsg() (b byte, errMsg string) {
	bs, errMsg := crb.ReadBySize(1)
	if errMsg != "" {
		return 0, errMsg
	}
	return bs[0], ""
}

func (crb *ReadBufferWrap) ReadLittleEndUint64() (x uint64, errMsg string) {
	buf, errMsg := crb.ReadBySize(8)
	if errMsg != "" {
		return 0, errMsg
	}
	x = binary.LittleEndian.Uint64(buf)
	return x, ""
}
