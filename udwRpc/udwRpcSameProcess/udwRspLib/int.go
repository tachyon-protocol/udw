package udwRspLib

import "unsafe"

func (buf *GoBuffer) ReadInt() int {
	buf.align(sizeOfInt)
	l := *(*int)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += sizeOfInt
	return l
}
func (buf *GoBuffer) WriteInt(i int) {
	buf.align(sizeOfInt)
	buf.writeNeedAddSize(sizeOfInt)
	*(*int)(unsafe.Pointer(&buf.buf[buf.offset])) = i
	buf.offset += sizeOfInt
}

func (buf *GoBuffer) ReadUint8() uint8 {
	buf.align(1)
	l := buf.buf[buf.offset]
	buf.offset += 1
	return l
}
func (buf *GoBuffer) WriteUint8(i uint8) {
	buf.align(1)
	buf.writeNeedAddSize(1)
	buf.buf[buf.offset] = i
	buf.offset += 1
}

func (buf *GoBuffer) ReadUint16() uint16 {
	buf.align(2)
	l := *(*uint16)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 2
	return l
}
func (buf *GoBuffer) WriteUint16(i uint16) {
	buf.align(2)
	buf.writeNeedAddSize(2)
	*(*uint16)(unsafe.Pointer(&buf.buf[buf.offset])) = i
	buf.offset += 2
}

func (buf *GoBuffer) ReadUint32() uint32 {
	buf.align(4)
	l := *(*uint32)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 4
	return l
}
func (buf *GoBuffer) WriteUint32(i uint32) {
	buf.align(4)
	buf.writeNeedAddSize(4)
	*(*uint32)(unsafe.Pointer(&buf.buf[buf.offset])) = i
	buf.offset += 4
}

func (buf *GoBuffer) ReadUint64() uint64 {
	buf.align(8)
	l := *(*uint64)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 8
	return l
}
func (buf *GoBuffer) WriteUint64(i uint64) {
	buf.align(8)
	buf.writeNeedAddSize(8)
	*(*uint64)(unsafe.Pointer(&buf.buf[buf.offset])) = i
	buf.offset += 8
}

func (buf *GoBuffer) ReadInt64() int64 {
	buf.align(8)
	l := *(*int64)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += 8
	return l
}

func (buf *GoBuffer) WriteInt64(i int64) {
	buf.align(8)
	buf.writeNeedAddSize(8)
	*(*int64)(unsafe.Pointer(&buf.buf[buf.offset])) = i
	buf.offset += 8
}
