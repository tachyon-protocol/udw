// +build windows

package udwW32

import "syscall"

func MustUdwDefaultMultiByteToUtf8String(output []byte) string {
	output1 := make([]uint16, len(output)*2)
	outSize := MustMultiByteToWideChar(
		CP_ACP,
		MB_PRECOMPOSED,
		&output[0],
		int32(len(output)),
		&output1[0],
		int32(len(output1)),
	)
	outputUtf16 := output1[:outSize]
	return syscall.UTF16ToString(outputUtf16)
}
