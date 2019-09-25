package udwStrings

import (
	"reflect"
	"unsafe"
)

func GetStringFromByteArrayNoAlloc(b []byte) string {
	var s string
	sx := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sx.Data = bx.Data
	sx.Len = bx.Len
	return s
}

func GetByteArrayFromStringNoAlloc(s string) []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&s))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(s)
	bx.Cap = bx.Len
	return b
}
