package udwJsonLib

import (
	"reflect"
	"unsafe"
)

func ReadJsonTmpString(ctx *Context) (s string) {
	bs := readJsonStringToByteSlice(ctx)
	sx := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sx.Data = bx.Data
	sx.Len = bx.Len
	return s
}
