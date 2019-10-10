// +build linux
// +build !appengine

package udwFileFastWalk

import (
	"bytes"
	"syscall"
	"unsafe"
)

func direntNamlen(dirent *syscall.Dirent) uint64 {
	const fixedHdr = uint16(unsafe.Offsetof(syscall.Dirent{}.Name))
	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
	const nameBufLen = uint16(len(nameBuf))
	limit := dirent.Reclen - fixedHdr
	if limit > nameBufLen {
		limit = nameBufLen
	}
	nameLen := bytes.IndexByte(nameBuf[:limit], 0)
	if nameLen < 0 {
		panic("failed to find terminating 0 byte in dirent")
	}
	return uint64(nameLen)
}
