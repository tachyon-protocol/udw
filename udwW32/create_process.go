// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"syscall"
	"unsafe"
)

var (
	procCreateProcessW      = DllProc{Dll: &modkernel32, Name: "CreateProcessW"}
	procTerminateProcess    = DllProc{Dll: &modkernel32, Name: "TerminateProcess"}
	procGetExitCodeProcess  = DllProc{Dll: &modkernel32, Name: "GetExitCodeProcess"}
	procWaitForSingleObject = DllProc{Dll: &modkernel32, Name: "WaitForSingleObject"}
)

func CreateProcessW(
	lpApplicationName, lpCommandLine string,
	lpProcessAttributes, lpThreadAttributes *SECURITY_ATTRIBUTES,
	bInheritHandles BOOL,
	dwCreationFlags uint32,
	lpEnvironment unsafe.Pointer,
	lpCurrentDirectory string,
	lpStartupInfo *STARTUPINFOW,
	lpProcessInformation *PROCESS_INFORMATION,
) (e error) {

	var lpAN, lpCL, lpCD *uint16
	if len(lpApplicationName) > 0 {
		lpAN, e = syscall.UTF16PtrFromString(lpApplicationName)
		if e != nil {
			return
		}
	}
	if len(lpCommandLine) > 0 {
		lpCL, e = syscall.UTF16PtrFromString(lpCommandLine)
		if e != nil {
			return
		}
	}
	if len(lpCurrentDirectory) > 0 {
		lpCD, e = syscall.UTF16PtrFromString(lpCurrentDirectory)
		if e != nil {
			return
		}
	}

	ret, _, lastErr := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(lpAN)),
		uintptr(unsafe.Pointer(lpCL)),
		uintptr(unsafe.Pointer(lpProcessAttributes)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
		uintptr(bInheritHandles),
		uintptr(dwCreationFlags),
		uintptr(lpEnvironment),
		uintptr(unsafe.Pointer(lpCD)),
		uintptr(unsafe.Pointer(lpStartupInfo)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

func CreateProcessQuick(cmd string) (pi PROCESS_INFORMATION, e error) {
	si := &STARTUPINFOW{}
	e = CreateProcessW(
		"",
		cmd,
		nil,
		nil,
		0,
		0,
		unsafe.Pointer(nil),
		"",
		si,
		&pi,
	)
	return
}

func TerminateProcess(hProcess HANDLE, exitCode uint32) (e error) {
	ret, _, lastErr := procTerminateProcess.Call(
		uintptr(hProcess),
		uintptr(exitCode),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

func GetExitCodeProcess(hProcess HANDLE) (code uintptr, e error) {
	ret, _, lastErr := procGetExitCodeProcess.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&code)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

func WaitForSingleObject(hHandle HANDLE, msecs uint32) (ok bool, e error) {

	ret, _, lastErr := procWaitForSingleObject.Call(
		uintptr(hHandle),
		uintptr(msecs),
	)

	if ret == WAIT_OBJECT_0 {
		ok = true
		return
	}

	if ret != WAIT_TIMEOUT {
		e = lastErr
	}
	return

}

const (
	WAIT_ABANDONED = 0x00000080
	WAIT_OBJECT_0  = 0x00000000
	WAIT_TIMEOUT   = 0x00000102
	WAIT_FAILED    = 0xFFFFFFFF
	INFINITE       = 0xFFFFFFFF
)

type PROCESS_INFORMATION struct {
	Process   HANDLE
	Thread    HANDLE
	ProcessId uint32
	ThreadId  uint32
}

type STARTUPINFOW struct {
	cb            uint32
	_             *uint16
	Desktop       *uint16
	Title         *uint16
	X             uint32
	Y             uint32
	XSize         uint32
	YSize         uint32
	XCountChars   uint32
	YCountChars   uint32
	FillAttribute uint32
	Flags         uint32
	ShowWindow    uint16
	_             uint16
	_             *uint8
	StdInput      HANDLE
	StdOutput     HANDLE
	StdError      HANDLE
}

type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      BOOL
}
