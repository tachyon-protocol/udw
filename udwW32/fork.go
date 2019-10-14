// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"fmt"
	"unsafe"
)

var (
	procRtlCloneUserProcess = DllProc{Dll: &modntdll, Name: "RtlCloneUserProcess"}
	procAllocConsole        = DllProc{Dll: &modkernel32, Name: "AllocConsole"}
	procOpenProcess         = DllProc{Dll: &modkernel32, Name: "OpenProcess"}
	procOpenThread          = DllProc{Dll: &modkernel32, Name: "OpenThread"}
	procResumeThread        = DllProc{Dll: &modkernel32, Name: "ResumeThread"}
)

func OpenProcess(desiredAccess int, inheritHandle bool, processId uintptr) (h HANDLE, e error) {
	inherit := uintptr(0)
	if inheritHandle {
		inherit = 1
	}

	ret, _, lastErr := procOpenProcess.Call(
		uintptr(desiredAccess),
		inherit,
		uintptr(processId),
	)

	if ret == 0 {
		e = lastErr
	}

	h = HANDLE(ret)
	return
}

func OpenThread(desiredAccess int, inheritHandle bool, threadId uintptr) (h HANDLE, e error) {
	inherit := uintptr(0)
	if inheritHandle {
		inherit = 1
	}

	ret, _, lastErr := procOpenThread.Call(
		uintptr(desiredAccess),
		inherit,
		uintptr(threadId),
	)

	if ret == 0 {
		e = lastErr
	}

	h = HANDLE(ret)
	return
}

func ResumeThread(ht HANDLE) (e error) {

	ret, _, lastErr := procResumeThread.Call(
		uintptr(ht),
	)
	if ret == ^uintptr(0) {
		e = lastErr
	}
	return
}

func AllocConsole() (e error) {
	ret, _, lastErr := procAllocConsole.Call()
	if ret != ERROR_SUCCESS {
		e = lastErr
	}
	return
}

func RtlCloneUserProcess(
	ProcessFlags uint32,
	ProcessSecurityDescriptor, ThreadSecurityDescriptor *SECURITY_DESCRIPTOR,
	DebugPort HANDLE,
	ProcessInformation *RTL_USER_PROCESS_INFORMATION,
) (status uintptr) {

	status, _, _ = procRtlCloneUserProcess.Call(
		uintptr(ProcessFlags),
		uintptr(unsafe.Pointer(ProcessSecurityDescriptor)),
		uintptr(unsafe.Pointer(ThreadSecurityDescriptor)),
		uintptr(DebugPort),
		uintptr(unsafe.Pointer(ProcessInformation)),
	)

	return
}

func Fork() (pid uintptr, e error) {

	pi := &RTL_USER_PROCESS_INFORMATION{}

	ret := RtlCloneUserProcess(
		RTL_CLONE_PROCESS_FLAGS_CREATE_SUSPENDED|RTL_CLONE_PROCESS_FLAGS_INHERIT_HANDLES,
		nil,
		nil,
		HANDLE(0),
		pi,
	)

	switch ret {
	case RTL_CLONE_PARENT:
		pid = pi.ClientId.UniqueProcess
		ht, err := OpenThread(THREAD_ALL_ACCESS, false, pi.ClientId.UniqueThread)
		if err != nil {
			e = fmt.Errorf("OpenThread: %s", err)
		}
		err = ResumeThread(ht)
		if err != nil {
			e = fmt.Errorf("ResumeThread: %s", err)
		}
		CloseHandle(ht)
	case RTL_CLONE_CHILD:
		pid = 0
		err := AllocConsole()
		if err != nil {
			e = fmt.Errorf("AllocConsole: %s", err)
		}
	default:
		e = fmt.Errorf("0x%x", ret)
	}
	return
}
