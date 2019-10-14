// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"strconv"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = Dll{Name: "kernel32.dll"}

	procGetModuleHandle    = DllProc{Dll: &modkernel32, Name: "GetModuleHandleW"}
	procMulDiv             = DllProc{Dll: &modkernel32, Name: "MulDiv"}
	procGetCurrentProcess  = DllProc{Dll: &modkernel32, Name: "GetCurrentProcess"}
	procGetConsoleWindow   = DllProc{Dll: &modkernel32, Name: "GetConsoleWindow"}
	procGetCurrentThread   = DllProc{Dll: &modkernel32, Name: "GetCurrentThread"}
	procGetLogicalDrives   = DllProc{Dll: &modkernel32, Name: "GetLogicalDrives"}
	procGetUserDefaultLCID = DllProc{Dll: &modkernel32, Name: "GetUserDefaultLCID"}
	procLstrlen            = DllProc{Dll: &modkernel32, Name: "lstrlenW"}
	procLstrcpy            = DllProc{Dll: &modkernel32, Name: "lstrcpyW"}
	procGlobalAlloc        = DllProc{Dll: &modkernel32, Name: "GlobalAlloc"}
	procGlobalFree         = DllProc{Dll: &modkernel32, Name: "GlobalFree"}
	procGlobalLock         = DllProc{Dll: &modkernel32, Name: "GlobalLock"}
	procGlobalUnlock       = DllProc{Dll: &modkernel32, Name: "GlobalUnlock"}
	procMoveMemory         = DllProc{Dll: &modkernel32, Name: "RtlMoveMemory"}
	procFindResource       = DllProc{Dll: &modkernel32, Name: "FindResourceW"}
	procSizeofResource     = DllProc{Dll: &modkernel32, Name: "SizeofResource"}
	procLockResource       = DllProc{Dll: &modkernel32, Name: "LockResource"}
	procLoadResource       = DllProc{Dll: &modkernel32, Name: "LoadResource"}
	procGetLastError       = DllProc{Dll: &modkernel32, Name: "GetLastError"}

	procCloseHandle                = DllProc{Dll: &modkernel32, Name: "CloseHandle"}
	procCreateToolhelp32Snapshot   = DllProc{Dll: &modkernel32, Name: "CreateToolhelp32Snapshot"}
	procModule32First              = DllProc{Dll: &modkernel32, Name: "Module32FirstW"}
	procModule32Next               = DllProc{Dll: &modkernel32, Name: "Module32NextW"}
	procGetSystemTimes             = DllProc{Dll: &modkernel32, Name: "GetSystemTimes"}
	procGetConsoleScreenBufferInfo = DllProc{Dll: &modkernel32, Name: "GetConsoleScreenBufferInfo"}
	procSetConsoleTextAttribute    = DllProc{Dll: &modkernel32, Name: "SetConsoleTextAttribute"}
	procGetDiskFreeSpaceEx         = DllProc{Dll: &modkernel32, Name: "GetDiskFreeSpaceExW"}
	procGetProcessTimes            = DllProc{Dll: &modkernel32, Name: "GetProcessTimes"}
	procSetSystemTime              = DllProc{Dll: &modkernel32, Name: "SetSystemTime"}
	procGetSystemTime              = DllProc{Dll: &modkernel32, Name: "GetSystemTime"}
	procVirtualAllocEx             = DllProc{Dll: &modkernel32, Name: "VirtualAllocEx"}
	procVirtualFreeEx              = DllProc{Dll: &modkernel32, Name: "VirtualFreeEx"}
	procWriteProcessMemory         = DllProc{Dll: &modkernel32, Name: "WriteProcessMemory"}
	procReadProcessMemory          = DllProc{Dll: &modkernel32, Name: "ReadProcessMemory"}
	procQueryPerformanceCounter    = DllProc{Dll: &modkernel32, Name: "QueryPerformanceCounter"}
	procQueryPerformanceFrequency  = DllProc{Dll: &modkernel32, Name: "QueryPerformanceFrequency"}
	procQueryDosDevice             = DllProc{Dll: &modkernel32, Name: "QueryDosDeviceW"}
	procIsWow64Process             = DllProc{Dll: &modkernel32, Name: "IsWow64Process"}
	procGetUserGeoID               = DllProc{Dll: &modkernel32, Name: "GetUserGeoID"}
	procGetGeoInfo                 = DllProc{Dll: &modkernel32, Name: "GetGeoInfoW"}
)

const (
	GEOCLASS_NATION = 16
	GEOCLASS_REGION = 14
	GEOCLASS_ALL    = 0
	GEO_ISO2        = 0x0004
)

func GetUserGeoID() string {
	ret, _, err := procGetUserGeoID.Call(uintptr(GEOCLASS_NATION))
	if ret == 0 {
		panic("GetUserGeoID " + SyscallErrorToMsg(err))
	}
	outData := make([]uint16, 10)
	ret, _, err = procGetGeoInfo.Call(ret, uintptr(GEO_ISO2), uintptr(unsafe.Pointer(&outData[0])), 10, 0)
	if ret == 0 {
		panic("GetGeoInfo " + SyscallErrorToMsg(err))
	}
	return syscall.UTF16ToString(outData)
}

func IsSys64() (result bool) {
	ret, _, err := procIsWow64Process.Call(GetCurrentProcess(), uintptr(unsafe.Pointer(&result)))
	if ret == 0 {
		panic("IsSys64 " + SyscallErrorToMsg(err))
	}
	return result
}

func MustQueryDosDevice(path string) string {
	var mn uintptr
	if path == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path)))
	}
	lpTargetPath := make([]uint16, 260)
	ret, _, err := procQueryDosDevice.Call(mn, uintptr(unsafe.Pointer(&lpTargetPath[0])), 260)
	if ret == 0 {
		panic("MustQueryDosDevice " + SyscallErrorToMsg(err))
	}
	return syscall.UTF16ToString(lpTargetPath)
}

func MustGetModuleHandle(modulename string) HINSTANCE {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, err := procGetModuleHandle.Call(mn)
	if ret == 0 {
		panic("MustGetModuleHandle " + SyscallErrorToMsg(err))
	}
	return HINSTANCE(ret)
}

func MulDiv(number, numerator, denominator int) int {
	ret, _, _ := procMulDiv.Call(
		uintptr(number),
		uintptr(numerator),
		uintptr(denominator))

	return int(ret)
}

func GetConsoleWindow() HWND {
	ret, _, _ := procGetConsoleWindow.Call()

	return HWND(ret)
}

func GetCurrentProcess() uintptr {
	ret, _, _ := procGetCurrentProcess.Call()

	return ret
}

func GetCurrentThread() HANDLE {
	ret, _, _ := procGetCurrentThread.Call()

	return HANDLE(ret)
}

func GetLogicalDrives() uint32 {
	ret, _, _ := procGetLogicalDrives.Call()

	return uint32(ret)
}

func GetUserDefaultLCID() uint32 {
	ret, _, _ := procGetUserDefaultLCID.Call()

	return uint32(ret)
}

func Lstrlen(lpString *uint16) int {
	ret, _, _ := procLstrlen.Call(uintptr(unsafe.Pointer(lpString)))

	return int(ret)
}

func Lstrcpy(buf []uint16, lpString *uint16) {
	procLstrcpy.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(lpString)))
}

func GlobalAlloc(uFlags uint, dwBytes uint32) HGLOBAL {
	ret, _, _ := procGlobalAlloc.Call(
		uintptr(uFlags),
		uintptr(dwBytes))

	if ret == 0 {
		panic("GlobalAlloc failed")
	}

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) {
	ret, _, _ := procGlobalFree.Call(uintptr(hMem))

	if ret != 0 {
		panic("GlobalFree failed")
	}
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := procGlobalLock.Call(uintptr(hMem))

	if ret == 0 {
		panic("GlobalLock failed")
	}

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := procGlobalUnlock.Call(uintptr(hMem))

	return ret != 0
}

func MoveMemory(destination, source unsafe.Pointer, length uint32) {
	procMoveMemory.Call(
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func FindResource(hModule HMODULE, lpName, lpType *uint16) (HRSRC, error) {
	ret, _, _ := procFindResource.Call(
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpName)),
		uintptr(unsafe.Pointer(lpType)))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HRSRC(ret), nil
}

func SizeofResource(hModule HMODULE, hResInfo HRSRC) uint32 {
	ret, _, _ := procSizeofResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("SizeofResource failed")
	}

	return uint32(ret)
}

func LockResource(hResData HGLOBAL) unsafe.Pointer {
	ret, _, _ := procLockResource.Call(uintptr(hResData))

	if ret == 0 {
		panic("LockResource failed")
	}

	return unsafe.Pointer(ret)
}

func LoadResource(hModule HMODULE, hResInfo HRSRC) HGLOBAL {
	ret, _, _ := procLoadResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("LoadResource failed")
	}

	return HGLOBAL(ret)
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}

func CreateToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

func Module32First(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func GetSystemTimes(lpIdleTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(lpIdleTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetProcessTimes(hProcess HANDLE, lpCreationTime, lpExitTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetProcessTimes.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(lpCreationTime)),
		uintptr(unsafe.Pointer(lpExitTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetConsoleScreenBufferInfo(hConsoleOutput HANDLE) *CONSOLE_SCREEN_BUFFER_INFO {
	var csbi CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(hConsoleOutput),
		uintptr(unsafe.Pointer(&csbi)))
	if ret == 0 {
		return nil
	}
	return &csbi
}

func SetConsoleTextAttribute(hConsoleOutput HANDLE, wAttributes uint16) bool {
	ret, _, _ := procSetConsoleTextAttribute.Call(
		uintptr(hConsoleOutput),
		uintptr(wAttributes))
	return ret != 0
}

func GetDiskFreeSpaceEx(dirName string) (r bool,
	freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64) {
	ret, _, _ := procGetDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dirName))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)))
	return ret != 0,
		freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes
}

func GetSystemTime() *SYSTEMTIME {
	var time SYSTEMTIME
	procGetSystemTime.Call(
		uintptr(unsafe.Pointer(&time)))
	return &time
}

func SetSystemTime(time *SYSTEMTIME) bool {
	ret, _, _ := procSetSystemTime.Call(
		uintptr(unsafe.Pointer(time)))
	return ret != 0
}

func VirtualAllocEx(hProcess HANDLE, lpAddress, dwSize uintptr, flAllocationType, flProtect uint32) uintptr {
	ret, _, _ := procVirtualAllocEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(flAllocationType),
		uintptr(flProtect),
	)

	return ret
}

func VirtualFreeEx(hProcess HANDLE, lpAddress, dwSize uintptr, dwFreeType uint32) bool {
	ret, _, _ := procVirtualFreeEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(dwFreeType),
	)

	return ret != 0
}

func WriteProcessMemory(hProcess HANDLE, lpBaseAddress, lpBuffer, nSize uintptr) (int, bool) {
	var nBytesWritten int
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		lpBuffer,
		nSize,
		uintptr(unsafe.Pointer(&nBytesWritten)),
	)

	return nBytesWritten, ret != 0
}

func ReadProcessMemory(hProcess HANDLE, lpBaseAddress, nSize uintptr) (lpBuffer []uint16, lpNumberOfBytesRead int, ok bool) {

	var nBytesRead int
	buf := make([]uint16, nSize)
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&buf[0])),
		nSize,
		uintptr(unsafe.Pointer(&nBytesRead)),
	)

	return buf, nBytesRead, ret != 0
}

func QueryPerformanceCounter() uint64 {
	result := uint64(0)
	procQueryPerformanceCounter.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}

func QueryPerformanceFrequency() uint64 {
	result := uint64(0)
	procQueryPerformanceFrequency.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}

const CP_ACP = 0
const CP_UTF8 = 65001
const CP_CHINESE_SIMPLE = 936
const CP_GB2312 = 936
const CP_US = 437

var procSetConsoleOutputCP = DllProc{Dll: &modkernel32, Name: "SetConsoleOutputCP"}

func MustSetConsoleOutputCP(wCodePageID uint32) {
	_, _, err := procSetConsoleOutputCP.Call(uintptr(wCodePageID))
	if IsSyscallErrorHappen(err) {
		panic("MustSetConsoleOutputCP " + err.Error())
	}
}

const MB_PRECOMPOSED = 1

var procMultiByteToWideChar = DllProc{Dll: &modkernel32, Name: "MultiByteToWideChar"}

func MustMultiByteToWideChar(
	wCodePage uint32,
	dwFlag uint32,
	lpMultiByteStr *uint8,
	cbMultiByte int32,
	lpWideCharStr *uint16,
	cchWideChar int32,
) int32 {
	r1, _, err := procMultiByteToWideChar.Call(
		uintptr(wCodePage),
		uintptr(dwFlag),
		uintptr(unsafe.Pointer(lpMultiByteStr)),
		uintptr(cbMultiByte),
		uintptr(unsafe.Pointer(lpWideCharStr)),
		uintptr(cchWideChar),
	)
	if IsSyscallErrorHappen(err) {
		panic("MustMultiByteToWideChar " + err.Error())
	}
	return int32(r1)
}

var procGetCurrentThreadId = DllProc{Dll: &modkernel32, Name: "GetCurrentThreadId"}

func GetCurrentThreadId() uint32 {
	ret, _, _ := procGetCurrentThreadId.Call()
	return uint32(ret)
}

var procExpandEnvironmentStringsW = DllProc{Dll: &modkernel32, Name: "ExpandEnvironmentStringsW"}

func ExpandEnvironmentStrings(src *uint16, dst *uint16, size uint32) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall(procExpandEnvironmentStringsW.Addr(), 3, uintptr(unsafe.Pointer(src)), uintptr(unsafe.Pointer(dst)), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

var procQueryFullProcessImageNameW = DllProc{Dll: &modkernel32, Name: "QueryFullProcessImageNameW"}

func QueryFullProcessImageNameW(hProcess HANDLE,
	dwFlags uint32,
	lpExeName *uint16,
	lpdwSize *uint32) (errMsg string) {
	ret, _, errMsg := procQueryFullProcessImageNameW.CallErrorMsg(uintptr(hProcess), uintptr(dwFlags), uintptr(unsafe.Pointer(lpExeName)), uintptr(unsafe.Pointer(lpdwSize)))
	if ret != 0 {
		return ""
	}
	return errMsg
}

const PROCESS_NAME_NATIVE = 1

var procGetFullPathNameW = DllProc{Dll: &modkernel32, Name: "GetFullPathNameW"}

func GetFullPathNameW(lpFileName *uint16,
	nBufferLength uint32,
	lpBuffer *uint16,
	lpFilePart *uint16) (size uint32, errMsg string) {
	ret, _, errMsg := procGetFullPathNameW.CallErrorMsg(uintptr(unsafe.Pointer(lpFileName)), uintptr(nBufferLength), uintptr(unsafe.Pointer(lpBuffer)), uintptr(unsafe.Pointer(lpFilePart)))
	if ret == 0 {
		return 0, errMsg
	}
	return uint32(ret), errMsg
}

func GetFullPathNameGo(filename string) (out string, errMsg string) {
	bufSize := uint32(256)
	buf := make([]uint16, bufSize)
	size, errMsg := GetFullPathNameW(syscall.StringToUTF16Ptr(filename), bufSize, &buf[0], nil)
	if errMsg != "" {
		return "", errMsg
	}
	if size > 1024*1024 {
		return "", "32vb7wfwjn " + strconv.Itoa(int(size))
	}
	if size > bufSize {
		bufSize = size
		buf := make([]uint16, bufSize)
		size, errMsg = GetFullPathNameW(syscall.StringToUTF16Ptr(filename), bufSize, &buf[0], nil)
		if errMsg != "" {
			return "", errMsg
		}
		if size > bufSize {
			return "", "z4h6txar2z " + strconv.Itoa(int(size))
		}
	}
	return syscall.UTF16ToString(buf), ""
}

var procGetLongPathNameW = DllProc{Dll: &modkernel32, Name: "GetLongPathNameW"}

func GetLongPathNameW(lpszShortPath *uint16,
	lpszLongPath *uint16,
	cchBuffer uint32) (size uint32, errMsg string) {
	ret, _, errMsg := procGetLongPathNameW.CallErrorMsg(uintptr(unsafe.Pointer(lpszShortPath)), uintptr(unsafe.Pointer(lpszLongPath)), uintptr(cchBuffer))
	if ret == 0 {
		return 0, errMsg
	}
	return uint32(ret), errMsg
}

func GetLongPathNameGo(filename string) (out string, errMsg string) {
	bufSize := uint32(256)
	buf := make([]uint16, bufSize)
	size, errMsg := GetLongPathNameW(syscall.StringToUTF16Ptr(filename), &buf[0], bufSize)
	if errMsg != "" {
		return "", errMsg
	}
	if size > 1024*1024 {
		return "", "gpr2pknw8f " + strconv.Itoa(int(size))
	}
	if size > bufSize {
		bufSize = size
		buf := make([]uint16, bufSize)
		size, errMsg = GetLongPathNameW(syscall.StringToUTF16Ptr(filename), &buf[0], bufSize)
		if errMsg != "" {
			return "", errMsg
		}
		if size > bufSize {
			return "", "dt6xxc8rg4 " + strconv.Itoa(int(size))
		}
	}
	return syscall.UTF16ToString(buf), ""
}
