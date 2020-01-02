package udwRspLib

import (
	"unsafe"
)

func (buf *GoBuffer) WriteBool(b bool) {

	buf.writeNeedAddSize(1)
	if b == true {
		buf.buf[buf.offset] = 1
	} else {
		buf.buf[buf.offset] = 0
	}
	buf.offset += 1
}
func (buf *GoBuffer) ReadBool() bool {
	out := (buf.buf[buf.offset] == 1)
	buf.offset += 1
	return out
}

func (buf *GoBuffer) WriteFloat64(f float64) {
	buf.align(8)
	buf.writeNeedAddSize(8)
	*(*float64)(unsafe.Pointer(&buf.buf[buf.offset])) = f
	buf.offset += 8
}
func (buf *GoBuffer) ReadFloat64() float64 {
	buf.align(8)
	l := *(*float64)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 8
	return l
}
func (buf *GoBuffer) WriteFloat32(f float32) {
	buf.align(4)
	buf.writeNeedAddSize(4)
	*(*float32)(unsafe.Pointer(&buf.buf[buf.offset])) = f
	buf.offset += 4
}
func (buf *GoBuffer) ReadFloat32() float32 {
	buf.align(4)
	l := *(*float32)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 4
	return l
}

func (buf *GoBuffer) WriteByteSlice(s []byte) {
	buf.align(sizeOfInt)
	buf.writeNeedAddSize(len(s) + sizeOfInt)
	*(*int)(unsafe.Pointer(&buf.buf[buf.offset])) = len(s)
	buf.offset += sizeOfInt
	copy(buf.buf[buf.offset:], s)
	buf.offset += len(s)
}

func (buf *GoBuffer) ReadByteSlice() []byte {
	buf.align(sizeOfInt)
	l := *(*int)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += sizeOfInt
	s := make([]byte, l)
	copy(s, buf.buf[buf.offset:buf.offset+l])
	buf.offset += l
	return s
}
