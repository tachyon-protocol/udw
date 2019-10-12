// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"unsafe"
)

var (
	modcomdlg32 = Dll{Name: "comdlg32.dll"}

	procGetSaveFileName      = DllProc{Dll: &modcomdlg32, Name: "GetSaveFileNameW"}
	procGetOpenFileName      = DllProc{Dll: &modcomdlg32, Name: "GetOpenFileNameW"}
	procCommDlgExtendedError = DllProc{Dll: &modcomdlg32, Name: "CommDlgExtendedError"}
)

func GetOpenFileName(ofn *OPENFILENAME) bool {
	ret, _, _ := procGetOpenFileName.Call(
		uintptr(unsafe.Pointer(ofn)))

	return ret != 0
}

func GetSaveFileName(ofn *OPENFILENAME) bool {
	ret, _, _ := procGetSaveFileName.Call(
		uintptr(unsafe.Pointer(ofn)))

	return ret != 0
}

func CommDlgExtendedError() uint {
	ret, _, _ := procCommDlgExtendedError.Call()

	return uint(ret)
}
