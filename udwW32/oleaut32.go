// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"syscall"
	"unsafe"
)

type DISPID int32
type VARTYPE uint16

const (
	DISP_E_MEMBERNOTFOUND = 0x80020003
)

var (
	modoleaut32 = Dll{Name: "oleaut32"}

	procVariantInit        = DllProc{Dll: &modoleaut32, Name: "VariantInit"}
	procSysAllocString     = DllProc{Dll: &modoleaut32, Name: "SysAllocString"}
	procSysFreeString      = DllProc{Dll: &modoleaut32, Name: "SysFreeString"}
	procSysStringLen       = DllProc{Dll: &modoleaut32, Name: "SysStringLen"}
	procCreateDispTypeInfo = DllProc{Dll: &modoleaut32, Name: "CreateDispTypeInfo"}
	procCreateStdDispatch  = DllProc{Dll: &modoleaut32, Name: "CreateStdDispatch"}
)

func VariantInit(v *VARIANT) {
	hr, _, _ := procVariantInit.Call(uintptr(unsafe.Pointer(v)))
	if hr != 0 {
		panic("Invoke VariantInit error.")
	}
	return
}

func SysAllocString(v string) (ss *uint16) {
	pss, _, _ := procSysAllocString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(v))))
	ss = (*uint16)(unsafe.Pointer(pss))
	return
}

func SysFreeString(v *uint16) {
	hr, _, _ := procSysFreeString.Call(uintptr(unsafe.Pointer(v)))
	if hr != 0 {
		panic("Invoke SysFreeString error.")
	}
	return
}

func SysStringLen(v *int16) uint {
	l, _, _ := procSysStringLen.Call(uintptr(unsafe.Pointer(v)))
	return uint(l)
}

func StringToVariantBSTR(value string) *VAR_BSTR {

	return &VAR_BSTR{vt: VT_BSTR, bstrVal: SysAllocString(value)}
}

type VAR_I4 struct {
	vt        VARTYPE
	reserved1 [6]byte
	lVal      int32
	reserved2 [4]byte
}

func IntToVariantI4(value int32) *VAR_I4 {
	return &VAR_I4{vt: VT_I4, lVal: value}
}
