// +build windows

package udwW32

import (
	"os"
	"unsafe"
)

func unpackUint64(cm uint64) (m1, m2 uintptr) {
	return uintptr(cm & 0xffffffff), uintptr(cm >> 32)
}

func packUint64(m1, m2 uintptr) uint64 {
	return uint64(m1) | (uint64(m2) << 32)
}

func VerSetConditionMask(lConditionMask uint64, typeBitMask uint32, conditionMask uint8) uint64 {
	m1, m2 := unpackUint64(lConditionMask)

	r1, r2, _ := procVerSetConditionMask.Call(m1, m2, uintptr(typeBitMask), uintptr(conditionMask))
	return packUint64(r1, r2)
}

func VerifyVersionInfoW(vi OSVersionInfoEx, typeMask uint32, conditionMask uint64) (bool, error) {
	vi.OSVersionInfoSize = uint32(unsafe.Sizeof(vi))
	cm1, cm2 := unpackUint64(conditionMask)

	r1, _, e1 := procVerifyVersionInfoW.Call(uintptr(unsafe.Pointer(&vi)), uintptr(typeMask), uintptr(cm1), uintptr(cm2))
	if r1 != 0 {
		return true, nil
	}
	if r1 == 0 && e1 == ERROR_OLD_WIN_VERSION {
		return false, nil
	}
	return false, os.NewSyscallError("VerifyVersionInfoW", e1)
}
