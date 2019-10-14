// +build windows

package udwW32

import (
	"syscall"
	"unsafe"
)

var modmsimg32 = Dll{Name: "msimg32.dll"}

var procAlphaBlend = DllProc{Dll: &modmsimg32, Name: "AlphaBlend"}

type BLENDFUNCTION struct {
	BlendOp             byte
	BlendFlags          byte
	SourceConstantAlpha byte
	AlphaFormat         byte
}

const AC_SRC_OVER = 0x00
const AC_SRC_ALPHA = 0x01

func MustAlphaBlend(dcdest HDC, xoriginDest int32, yoriginDest int32, wDest int32, hDest int32, dcsrc HDC, xoriginSrc int32, yoriginSrc int32, wsrc int32, hsrc int32, ftn BLENDFUNCTION) {
	r1, _, e1 := syscall.Syscall12(procAlphaBlend.Addr(), 11,
		uintptr(dcdest),
		uintptr(xoriginDest),
		uintptr(yoriginDest),
		uintptr(wDest),
		uintptr(hDest),
		uintptr(dcsrc),
		uintptr(xoriginSrc),
		uintptr(yoriginSrc),
		uintptr(wsrc),
		uintptr(hsrc),
		uintptr(*((*uintptr)(unsafe.Pointer(&ftn)))),
		0)
	if r1 == 0 {
		if e1 != 0 {
			panic("MustAlphaBlend fail " + SyscallErrorToMsg(e1))
		} else {
			panic("MustAlphaBlend fail syscall.EINVAL")
		}
	}
	return
}
