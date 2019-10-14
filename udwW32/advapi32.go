// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modadvapi32 = Dll{Name: "advapi32.dll"}

	procCloseEventLog                = DllProc{Dll: &modadvapi32, Name: "CloseEventLog"}
	procCloseServiceHandle           = DllProc{Dll: &modadvapi32, Name: "CloseServiceHandle"}
	procControlService               = DllProc{Dll: &modadvapi32, Name: "ControlService"}
	procControlTrace                 = DllProc{Dll: &modadvapi32, Name: "ControlTraceW"}
	procInitializeSecurityDescriptor = DllProc{Dll: &modadvapi32, Name: "InitializeSecurityDescriptor"}
	procOpenEventLog                 = DllProc{Dll: &modadvapi32, Name: "OpenEventLogW"}
	procOpenSCManager                = DllProc{Dll: &modadvapi32, Name: "OpenSCManagerW"}
	procOpenService                  = DllProc{Dll: &modadvapi32, Name: "OpenServiceW"}
	procReadEventLog                 = DllProc{Dll: &modadvapi32, Name: "ReadEventLogW"}
	procRegCloseKey                  = DllProc{Dll: &modadvapi32, Name: "RegCloseKey"}
	procRegEnumKeyEx                 = DllProc{Dll: &modadvapi32, Name: "RegEnumKeyExW"}
	procRegGetValue                  = DllProc{Dll: &modadvapi32, Name: "RegGetValueW"}
	procRegOpenKeyEx                 = DllProc{Dll: &modadvapi32, Name: "RegOpenKeyExW"}
	procRegSetValueEx                = DllProc{Dll: &modadvapi32, Name: "RegSetValueExW"}
	procSetSecurityDescriptorDacl    = DllProc{Dll: &modadvapi32, Name: "SetSecurityDescriptorDacl"}
	procStartService                 = DllProc{Dll: &modadvapi32, Name: "StartServiceW"}
	procStartTrace                   = DllProc{Dll: &modadvapi32, Name: "StartTraceW"}
)

var (
	SystemTraceControlGuid = GUID{
		0x9e814aad,
		0x3204,
		0x11d2,
		[8]byte{0x9a, 0x82, 0x00, 0x60, 0x08, 0xa8, 0x69, 0x39},
	}
)

var procRegCreateKeyExW = DllProc{Dll: &modadvapi32, Name: "RegCreateKeyExW"}

func RegCreateKeyEx(key syscall.Handle, subkey *uint16, reserved uint32, class *uint16, options uint32, desired uint32, sa *syscall.SecurityAttributes, result *syscall.Handle, disposition *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegCreateKeyExW.Addr(), 9, uintptr(key), uintptr(unsafe.Pointer(subkey)), uintptr(reserved), uintptr(unsafe.Pointer(class)), uintptr(options), uintptr(desired), uintptr(unsafe.Pointer(sa)), uintptr(unsafe.Pointer(result)), uintptr(unsafe.Pointer(disposition)))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegOpenKeyEx(hKey HKEY, subKey string, samDesired uint32) HKEY {
	var result HKEY
	ret, _, _ := procRegOpenKeyEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(0),
		uintptr(samDesired),
		uintptr(unsafe.Pointer(&result)))

	if ret != ERROR_SUCCESS {
		panic(fmt.Sprintf("RegOpenKeyEx(%d, %s, %d) failed", hKey, subKey, samDesired))
	}
	return result
}

func RegCloseKey(hKey HKEY) error {
	var err error
	ret, _, _ := procRegCloseKey.Call(
		uintptr(hKey))

	if ret != ERROR_SUCCESS {
		err = errors.New("RegCloseKey failed")
	}
	return err
}

func RegGetRaw(hKey HKEY, subKey string, value string) []byte {
	var bufLen uint32
	var valptr unsafe.Pointer
	if len(value) > 0 {
		valptr = unsafe.Pointer(syscall.StringToUTF16Ptr(value))
	}
	procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(valptr),
		uintptr(RRF_RT_ANY),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))

	if bufLen == 0 {
		return nil
	}

	buf := make([]byte, bufLen)
	ret, _, _ := procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(valptr),
		uintptr(RRF_RT_ANY),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))

	if ret != ERROR_SUCCESS {
		return nil
	}

	return buf
}

func RegSetBinary(hKey HKEY, subKey string, value []byte) (errno int) {
	var lptr, vptr unsafe.Pointer
	if len(subKey) > 0 {
		lptr = unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))
	}
	if len(value) > 0 {
		vptr = unsafe.Pointer(&value[0])
	}
	ret, _, _ := procRegSetValueEx.Call(
		uintptr(hKey),
		uintptr(lptr),
		uintptr(0),
		uintptr(REG_BINARY),
		uintptr(vptr),
		uintptr(len(value)))

	return int(ret)
}

func RegGetString(hKey HKEY, subKey string, value string) string {
	var bufLen uint32
	procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))

	if bufLen == 0 {
		return ""
	}

	buf := make([]uint16, bufLen)
	ret, _, _ := procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))

	if ret != ERROR_SUCCESS {
		return ""
	}

	return syscall.UTF16ToString(buf)
}

func RegEnumKeyEx(hKey HKEY, index uint32) string {
	var bufLen uint32 = 255
	buf := make([]uint16, bufLen)
	procRegEnumKeyEx.Call(
		uintptr(hKey),
		uintptr(index),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)),
		0,
		0,
		0,
		0)
	return syscall.UTF16ToString(buf)
}

func OpenEventLog(servername string, sourcename string) HANDLE {
	ret, _, _ := procOpenEventLog.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servername))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(sourcename))))

	return HANDLE(ret)
}

func ReadEventLog(eventlog HANDLE, readflags, recordoffset uint32, buffer []byte, numberofbytestoread uint32, bytesread, minnumberofbytesneeded *uint32) bool {
	ret, _, _ := procReadEventLog.Call(
		uintptr(eventlog),
		uintptr(readflags),
		uintptr(recordoffset),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(numberofbytestoread),
		uintptr(unsafe.Pointer(bytesread)),
		uintptr(unsafe.Pointer(minnumberofbytesneeded)))

	return ret != 0
}

func CloseEventLog(eventlog HANDLE) bool {
	ret, _, _ := procCloseEventLog.Call(
		uintptr(eventlog))

	return ret != 0
}

func OpenSCManager(lpMachineName, lpDatabaseName string, dwDesiredAccess uint32) (HANDLE, error) {
	var p1, p2 uintptr
	if len(lpMachineName) > 0 {
		p1 = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpMachineName)))
	}
	if len(lpDatabaseName) > 0 {
		p2 = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpDatabaseName)))
	}
	ret, _, _ := procOpenSCManager.Call(
		p1,
		p2,
		uintptr(dwDesiredAccess))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HANDLE(ret), nil
}

func CloseServiceHandle(hSCObject HANDLE) error {
	ret, _, _ := procCloseServiceHandle.Call(uintptr(hSCObject))
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func OpenService(hSCManager HANDLE, lpServiceName string, dwDesiredAccess uint32) (HANDLE, error) {
	ret, _, _ := procOpenService.Call(
		uintptr(hSCManager),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpServiceName))),
		uintptr(dwDesiredAccess))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HANDLE(ret), nil
}

func StartService(hService HANDLE, lpServiceArgVectors []string) error {
	l := len(lpServiceArgVectors)
	var ret uintptr
	if l == 0 {
		ret, _, _ = procStartService.Call(
			uintptr(hService),
			0,
			0)
	} else {
		lpArgs := make([]uintptr, l)
		for i := 0; i < l; i++ {
			lpArgs[i] = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpServiceArgVectors[i])))
		}

		ret, _, _ = procStartService.Call(
			uintptr(hService),
			uintptr(l),
			uintptr(unsafe.Pointer(&lpArgs[0])))
	}

	if ret == 0 {
		return syscall.GetLastError()
	}

	return nil
}

func ControlService(hService HANDLE, dwControl uint32, lpServiceStatus *SERVICE_STATUS) bool {
	if lpServiceStatus == nil {
		panic("ControlService:lpServiceStatus cannot be nil")
	}

	ret, _, _ := procControlService.Call(
		uintptr(hService),
		uintptr(dwControl),
		uintptr(unsafe.Pointer(lpServiceStatus)))

	return ret != 0
}

func ControlTrace(hTrace TRACEHANDLE, lpSessionName string, props *EVENT_TRACE_PROPERTIES, dwControl uint32) (success bool, e error) {

	ret, _, _ := procControlTrace.Call(
		uintptr(unsafe.Pointer(hTrace)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpSessionName))),
		uintptr(unsafe.Pointer(props)),
		uintptr(dwControl))

	if ret == ERROR_SUCCESS {
		return true, nil
	}
	e = errors.New(fmt.Sprintf("error: 0x%x", ret))
	return
}

func StartTrace(lpSessionName string, props *EVENT_TRACE_PROPERTIES) (hTrace TRACEHANDLE, e error) {

	ret, _, _ := procStartTrace.Call(
		uintptr(unsafe.Pointer(&hTrace)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpSessionName))),
		uintptr(unsafe.Pointer(props)))

	if ret == ERROR_SUCCESS {
		return
	}
	e = errors.New(fmt.Sprintf("error: 0x%x", ret))
	return
}

func InitializeSecurityDescriptor(rev uint16) (pSecurityDescriptor *SECURITY_DESCRIPTOR, e error) {

	pSecurityDescriptor = &SECURITY_DESCRIPTOR{}

	ret, _, _ := procInitializeSecurityDescriptor.Call(
		uintptr(unsafe.Pointer(pSecurityDescriptor)),
		uintptr(rev),
	)

	if ret != 0 {
		return
	}
	e = syscall.GetLastError()
	return
}

func SetSecurityDescriptorDacl(pSecurityDescriptor *SECURITY_DESCRIPTOR, pDacl *ACL) (e error) {

	if pSecurityDescriptor == nil {
		return errors.New("null descriptor")
	}

	var ret uintptr
	if pDacl == nil {
		ret, _, _ = procSetSecurityDescriptorDacl.Call(
			uintptr(unsafe.Pointer(pSecurityDescriptor)),
			uintptr(1),
			uintptr(0),
			uintptr(0),
		)
	} else {
		ret, _, _ = procSetSecurityDescriptorDacl.Call(
			uintptr(unsafe.Pointer(pSecurityDescriptor)),
			uintptr(1),
			uintptr(unsafe.Pointer(pDacl)),
			uintptr(0),
		)
	}

	if ret != 0 {
		return
	}
	e = syscall.GetLastError()
	return
}

var procRegDeleteKeyW = DllProc{Dll: &modadvapi32, Name: "RegDeleteKeyW"}

func RegDeleteKey(key syscall.Handle, subkey *uint16) (regerrno error) {
	r0, _, _ := syscall.Syscall(procRegDeleteKeyW.Addr(), 2, uintptr(key), uintptr(unsafe.Pointer(subkey)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

var procRegSetValueExW = DllProc{Dll: &modadvapi32, Name: "RegSetValueExW"}

func RegSetValueEx(key syscall.Handle, valueName *uint16, reserved uint32, vtype uint32, buf *byte, bufsize uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall6(procRegSetValueExW.Addr(), 6, uintptr(key), uintptr(unsafe.Pointer(valueName)), uintptr(reserved), uintptr(vtype), uintptr(unsafe.Pointer(buf)), uintptr(bufsize))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

var procRegEnumValueW = DllProc{Dll: &modadvapi32, Name: "RegEnumValueW"}
var procRegDeleteValueW = DllProc{Dll: &modadvapi32, Name: "RegDeleteValueW"}
var procRegLoadMUIStringW = DllProc{Dll: &modadvapi32, Name: "RegLoadMUIStringW"}

func RegEnumValue(key syscall.Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegEnumValueW.Addr(), 8, uintptr(key), uintptr(index), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(valtype)), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(buflen)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegDeleteValue(key syscall.Handle, name *uint16) (regerrno error) {
	r0, _, _ := syscall.Syscall(procRegDeleteValueW.Addr(), 2, uintptr(key), uintptr(unsafe.Pointer(name)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func RegLoadMUIString(key syscall.Handle, name *uint16, buf *uint16, buflen uint32, buflenCopied *uint32, flags uint32, dir *uint16) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegLoadMUIStringW.Addr(), 7, uintptr(key), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(buf)), uintptr(buflen), uintptr(unsafe.Pointer(buflenCopied)), uintptr(flags), uintptr(unsafe.Pointer(dir)), 0, 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

var SECURITY_NT_AUTHORITY = SID_IDENTIFIER_AUTHORITY{
	Value: [6]byte{0, 0, 0, 0, 0, 5},
}

const SECURITY_BUILTIN_DOMAIN_RID = 0x00000020
const DOMAIN_ALIAS_RID_ADMINS = 0x00000220

var procAllocateAndInitializeSid = DllProc{Dll: &modadvapi32, Name: "AllocateAndInitializeSid"}

func AllocateAndInitializeSid(
	pIdentifierAuthority *SID_IDENTIFIER_AUTHORITY,
	nSubAuthorityCount uint8,
	nSubAuthority0 uint32,
	nSubAuthority1 uint32,
	nSubAuthority2 uint32,
	nSubAuthority3 uint32,
	nSubAuthority4 uint32,
	nSubAuthority5 uint32,
	nSubAuthority6 uint32,
	nSubAuthority7 uint32,
	sid **SID,
) (errMsg string) {
	_, _, errMsg = procAllocateAndInitializeSid.CallErrorMsg(
		uintptr(unsafe.Pointer(pIdentifierAuthority)),
		uintptr(nSubAuthorityCount),
		uintptr(nSubAuthority0),
		uintptr(nSubAuthority1),
		uintptr(nSubAuthority2),
		uintptr(nSubAuthority3),
		uintptr(nSubAuthority4),
		uintptr(nSubAuthority5),
		uintptr(nSubAuthority6),
		uintptr(nSubAuthority7),
		uintptr(unsafe.Pointer(sid)),
	)
	return errMsg
}

var procCheckTokenMembership = DllProc{Dll: &modadvapi32, Name: "CheckTokenMembership"}

func CheckTokenMembership(TokenHandle HANDLE, SidToCheck *SID, IsMember *uint8) (errMsg string) {
	_, _, errMsg = procCheckTokenMembership.CallErrorMsg(
		uintptr(TokenHandle),
		uintptr(unsafe.Pointer(SidToCheck)),
		uintptr(unsafe.Pointer(IsMember)),
	)
	return errMsg
}

var procFreeSid = DllProc{Dll: &modadvapi32, Name: "FreeSid"}

func FreeSid(sid *SID) {

	procFreeSid.CallErrorMsg(
		uintptr(unsafe.Pointer(sid)),
	)
	return
}
