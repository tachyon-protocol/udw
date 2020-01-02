package udwRspLib

/*
#include <stdlib.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

var sizeOfInt = int(unsafe.Sizeof(int(0)))

type GoBuffer struct {
	buf    []byte
	offset int
}

func (buf *GoBuffer) writeNeedAddSize(l int) {
	needSize := buf.offset + l
	if len(buf.buf) >= needSize {
		return
	}

	toAllocSize := len(buf.buf)
	if toAllocSize < 64 {
		toAllocSize = 64
	}
	for {
		if toAllocSize >= needSize+64 {
			break
		}
		toAllocSize = toAllocSize * 2
	}
	if buf.buf == nil {
		buf.buf = unsafeRealloc(nil, toAllocSize)

	} else {
		buf.buf = unsafeRealloc(unsafe.Pointer(&buf.buf[0]), toAllocSize)
	}
}

func (buf *GoBuffer) align(alignment int) {
	pad := buf.offset % alignment
	if pad > 0 {
		buf.offset += alignment - pad
	}
}

func unsafeRealloc(oldPointer unsafe.Pointer, toAllocSize int) []byte {
	allocPointer := C.realloc(oldPointer, C.size_t(toAllocSize))
	var out []byte
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&(out)))
	bx.Data = uintptr(unsafe.Pointer(allocPointer))
	bx.Len = toAllocSize
	bx.Cap = toAllocSize
	return out
}

func (buf *GoBuffer) ToC() (p uintptr, c int) {
	if buf.buf == nil {
		buf.writeNeedAddSize(1)
	}

	pBuf := &buf.buf[0]
	p = uintptr(unsafe.Pointer(pBuf))

	c = cap(buf.buf)
	return
}

func (buf *GoBuffer) SetFromC(p uintptr, c int) {
	pBuf := &(buf.buf)
	bx := (*reflect.SliceHeader)(unsafe.Pointer(pBuf))
	bx.Data = p
	bx.Cap = c
	bx.Len = c
	buf.offset = 0
}

func (buf *GoBuffer) ResetToWrite() {
	buf.offset = 0
}

func (buf *GoBuffer) ResetToRead() {
	buf.offset = 0
}

func (buf *GoBuffer) FreeFromGo() {
	if buf.buf != nil {

		pBuf := &buf.buf[0]
		C.free(unsafe.Pointer(pBuf))

		buf.buf = nil
	}
}

func (buf *GoBuffer) GetOffset() int {
	return buf.offset
}

func (buf *GoBuffer) CopyToByteSlice() []byte {
	newBuf := make([]byte, buf.offset)
	copy(newBuf, buf.buf[:buf.offset])
	return newBuf
}

func NewGoBufferFromByteSlice(slice []byte) *GoBuffer {
	buf := &GoBuffer{}
	buf.writeNeedAddSize(len(slice))
	copy(buf.buf, slice)
	return buf
}

func NewGoBufferFromC(p uintptr, cap int) *GoBuffer {
	buf := &GoBuffer{}
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&(buf.buf)))
	bx.Data = p
	bx.Cap = cap
	bx.Len = cap
	buf.offset = 0
	return buf
}
