// +build amd64

package udwIpCountryV2Map

import (
	"unsafe"
)

func (r *Reader) ReadNode(nodeNumber uint32, index uint8) uint32 {
	offset := uintptr(nodeNumber*6) + uintptr(index*3)

	out := *(*uint32)(unsafe.Pointer(&r.Buf[offset]))

	return out & 0xffffff

}
