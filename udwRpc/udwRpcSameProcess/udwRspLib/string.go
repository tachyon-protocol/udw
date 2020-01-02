package udwRspLib

import (
	"reflect"
	"unicode/utf16"
	"unsafe"
)

func (buf *GoBuffer) WriteStringUTF8(s string) {
	buf.align(sizeOfInt)
	buf.writeNeedAddSize(len(s) + sizeOfInt)
	*(*int)(unsafe.Pointer(&buf.buf[buf.offset])) = len(s)
	buf.offset += sizeOfInt
	copy(buf.buf[buf.offset:], s)
	buf.offset += len(s)
}
func (buf *GoBuffer) ReadStringUTF8() string {
	buf.align(sizeOfInt)
	l := *(*int)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += sizeOfInt
	s := string(buf.buf[buf.offset : buf.offset+l])
	buf.offset += l
	return s
}
func (buf *GoBuffer) WriteStringUTF16(s string) {
	if len(s) == 0 {
		buf.WriteInt(0)
		return
	}
	buf.align(sizeOfInt)
	runeList := []rune(s)
	utf16S := utf16.Encode(runeList)
	buf.writeNeedAddSize(len(utf16S)*2 + sizeOfInt)
	*(*int)(unsafe.Pointer(&buf.buf[buf.offset])) = len(utf16S)
	buf.offset += sizeOfInt
	var utf16Bytes []byte
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&(utf16Bytes)))
	bx.Data = uintptr(unsafe.Pointer(&utf16S[0]))
	bx.Len = len(utf16S) * 2
	bx.Cap = len(utf16S) * 2
	copy(buf.buf[buf.offset:], utf16Bytes)
	buf.offset += 2 * len(utf16S)
}
func (buf *GoBuffer) ReadStringUTF16() string {
	buf.align(sizeOfInt)
	l := *(*int)(unsafe.Pointer(&buf.buf[buf.offset]))
	buf.offset += sizeOfInt
	if l == 0 {
		return ""
	}
	var utf16S []uint16
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&(utf16S)))
	bx.Data = uintptr(unsafe.Pointer(&buf.buf[buf.offset]))
	bx.Len = l
	bx.Cap = l
	runeList := utf16.Decode(utf16S)
	s := string(runeList)
	buf.offset += l * 2
	return s
}
