// +build windows

package udwW32

import (
	"syscall"
)

var procVerSetConditionMask = DllProc{Dll: &modkernel32, Name: "VerSetConditionMask"}

const (
	_WIN32_WINNT_NT4      = 0x0400
	_WIN32_WINNT_WIN2K    = 0x0500
	_WIN32_WINNT_WINXP    = 0x0501
	_WIN32_WINNT_WS03     = 0x0502
	_WIN32_WINNT_WIN6     = 0x0600
	_WIN32_WINNT_VISTA    = 0x0600
	_WIN32_WINNT_WS08     = 0x0600
	_WIN32_WINNT_LONGHORN = 0x0600
	_WIN32_WINNT_WIN7     = 0x0601
	_WIN32_WINNT_WIN8     = 0x0602
	_WIN32_WINNT_WINBLUE  = 0x0603
	_WIN32_WINNT_WIN10    = 0x0A00

	VER_NT_WORKSTATION = 0x0000001

	VER_BUILDNUMBER      = 0x0000004
	VER_MAJORVERSION     = 0x0000002
	VER_MINORVERSION     = 0x0000001
	VER_PLATFORMID       = 0x0000008
	VER_PRODUCT_TYPE     = 0x0000080
	VER_SERVICEPACKMAJOR = 0x0000020
	VER_SERVICEPACKMINOR = 0x0000010
	VER_SUITENAME        = 0x0000040

	VER_EQUAL         = 1
	VER_GREATER       = 2
	VER_GREATER_EQUAL = 3
	VER_LESS          = 4
	VER_LESS_EQUAL    = 5

	ERROR_OLD_WIN_VERSION syscall.Errno = 1150
)

var procVerifyVersionInfoW = DllProc{Dll: &modkernel32, Name: "VerifyVersionInfoW"}

type OSVersionInfoEx struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformId        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserve           byte
}

func IsWinVersionOrGreater(id uint32, wServicePackMajor uint16) bool {
	cm := VerSetConditionMask(0, VER_MAJORVERSION, VER_GREATER_EQUAL)
	cm = VerSetConditionMask(cm, VER_MINORVERSION, VER_GREATER_EQUAL)
	cm = VerSetConditionMask(cm, VER_SERVICEPACKMAJOR, VER_GREATER_EQUAL)
	cm = VerSetConditionMask(cm, VER_SERVICEPACKMINOR, VER_GREATER_EQUAL)
	r, _ := VerifyVersionInfoW(OSVersionInfoEx{
		MajorVersion:     (id >> 8 & 0xff),
		MinorVersion:     (id & 0xff),
		ServicePackMajor: wServicePackMajor,
	}, VER_MAJORVERSION|VER_MINORVERSION|VER_SERVICEPACKMAJOR|VER_SERVICEPACKMINOR, cm)
	return r
}

func IsWindowsXPOrGreater() bool       { return IsWinVersionOrGreater(_WIN32_WINNT_WINXP, 0) }
func IsWindowsXPSP1OrGreater() bool    { return IsWinVersionOrGreater(_WIN32_WINNT_WINXP, 1) }
func IsWindowsXPSP2OrGreater() bool    { return IsWinVersionOrGreater(_WIN32_WINNT_WINXP, 2) }
func IsWindowsXPSP3OrGreater() bool    { return IsWinVersionOrGreater(_WIN32_WINNT_WINXP, 3) }
func IsWindowsVistaOrGreater() bool    { return IsWinVersionOrGreater(_WIN32_WINNT_VISTA, 0) }
func IsWindowsVistaSP1OrGreater() bool { return IsWinVersionOrGreater(_WIN32_WINNT_VISTA, 1) }
func IsWindowsVistaSP2OrGreater() bool { return IsWinVersionOrGreater(_WIN32_WINNT_VISTA, 2) }
func IsWindows7OrGreater() bool        { return IsWinVersionOrGreater(_WIN32_WINNT_WIN7, 0) }
func IsWindows7SP1OrGreater() bool     { return IsWinVersionOrGreater(_WIN32_WINNT_WIN7, 1) }
func IsWindows8OrGreater() bool        { return IsWinVersionOrGreater(_WIN32_WINNT_WIN8, 0) }
func IsWindows8Point1OrGreater() bool  { return IsWinVersionOrGreater(_WIN32_WINNT_WINBLUE, 0) }
func IsWindows10OrGreater() bool       { return IsWinVersionOrGreater(_WIN32_WINNT_WIN10, 0) }

func IsWindowsServer() bool {
	cm := VerSetConditionMask(0, VER_PRODUCT_TYPE, VER_EQUAL)
	r, _ := VerifyVersionInfoW(OSVersionInfoEx{
		ProductType: VER_NT_WORKSTATION,
	}, VER_PRODUCT_TYPE, cm)
	return r
}
