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
	moduser32 = Dll{Name: "user32.dll"}

	procRegisterClassEx               = DllProc{Dll: &moduser32, Name: "RegisterClassExW"}
	procLoadIcon                      = DllProc{Dll: &moduser32, Name: "LoadIconW"}
	procLoadCursor                    = DllProc{Dll: &moduser32, Name: "LoadCursorW"}
	procShowWindow                    = DllProc{Dll: &moduser32, Name: "ShowWindow"}
	procUpdateWindow                  = DllProc{Dll: &moduser32, Name: "UpdateWindow"}
	procCreateWindowEx                = DllProc{Dll: &moduser32, Name: "CreateWindowExW"}
	procAdjustWindowRect              = DllProc{Dll: &moduser32, Name: "AdjustWindowRect"}
	procAdjustWindowRectEx            = DllProc{Dll: &moduser32, Name: "AdjustWindowRectEx"}
	procDestroyWindow                 = DllProc{Dll: &moduser32, Name: "DestroyWindow"}
	procDefWindowProc                 = DllProc{Dll: &moduser32, Name: "DefWindowProcW"}
	procDefDlgProc                    = DllProc{Dll: &moduser32, Name: "DefDlgProcW"}
	procPostQuitMessage               = DllProc{Dll: &moduser32, Name: "PostQuitMessage"}
	procGetMessage                    = DllProc{Dll: &moduser32, Name: "GetMessageW"}
	procTranslateMessage              = DllProc{Dll: &moduser32, Name: "TranslateMessage"}
	procDispatchMessage               = DllProc{Dll: &moduser32, Name: "DispatchMessageW"}
	procSendMessage                   = DllProc{Dll: &moduser32, Name: "SendMessageW"}
	procSendMessageTimeout            = DllProc{Dll: &moduser32, Name: "SendMessageTimeout"}
	procPostMessage                   = DllProc{Dll: &moduser32, Name: "PostMessageW"}
	procWaitMessage                   = DllProc{Dll: &moduser32, Name: "WaitMessage"}
	procSetWindowText                 = DllProc{Dll: &moduser32, Name: "SetWindowTextW"}
	procGetWindowTextLength           = DllProc{Dll: &moduser32, Name: "GetWindowTextLengthW"}
	procGetWindowText                 = DllProc{Dll: &moduser32, Name: "GetWindowTextW"}
	procGetWindowRect                 = DllProc{Dll: &moduser32, Name: "GetWindowRect"}
	procMoveWindow                    = DllProc{Dll: &moduser32, Name: "MoveWindow"}
	procScreenToClient                = DllProc{Dll: &moduser32, Name: "ScreenToClient"}
	procCallWindowProc                = DllProc{Dll: &moduser32, Name: "CallWindowProcW"}
	procSetWindowLong                 = DllProc{Dll: &moduser32, Name: "SetWindowLongW"}
	procSetWindowLongPtr              = DllProc{Dll: &moduser32, Name: "SetWindowLongW"}
	procGetWindowLong                 = DllProc{Dll: &moduser32, Name: "GetWindowLongW"}
	procGetWindowLongPtr              = DllProc{Dll: &moduser32, Name: "GetWindowLongW"}
	procEnableWindow                  = DllProc{Dll: &moduser32, Name: "EnableWindow"}
	procIsWindowEnabled               = DllProc{Dll: &moduser32, Name: "IsWindowEnabled"}
	procIsWindowVisible               = DllProc{Dll: &moduser32, Name: "IsWindowVisible"}
	procSetFocus                      = DllProc{Dll: &moduser32, Name: "SetFocus"}
	procInvalidateRect                = DllProc{Dll: &moduser32, Name: "InvalidateRect"}
	procGetClientRect                 = DllProc{Dll: &moduser32, Name: "GetClientRect"}
	procGetDC                         = DllProc{Dll: &moduser32, Name: "GetDC"}
	procReleaseDC                     = DllProc{Dll: &moduser32, Name: "ReleaseDC"}
	procSetCapture                    = DllProc{Dll: &moduser32, Name: "SetCapture"}
	procReleaseCapture                = DllProc{Dll: &moduser32, Name: "ReleaseCapture"}
	procGetWindowThreadProcessId      = DllProc{Dll: &moduser32, Name: "GetWindowThreadProcessId"}
	procMessageBox                    = DllProc{Dll: &moduser32, Name: "MessageBoxW"}
	procGetSystemMetrics              = DllProc{Dll: &moduser32, Name: "GetSystemMetrics"}
	procCopyRect                      = DllProc{Dll: &moduser32, Name: "CopyRect"}
	procEqualRect                     = DllProc{Dll: &moduser32, Name: "EqualRect"}
	procInflateRect                   = DllProc{Dll: &moduser32, Name: "InflateRect"}
	procIntersectRect                 = DllProc{Dll: &moduser32, Name: "IntersectRect"}
	procIsRectEmpty                   = DllProc{Dll: &moduser32, Name: "IsRectEmpty"}
	procOffsetRect                    = DllProc{Dll: &moduser32, Name: "OffsetRect"}
	procPtInRect                      = DllProc{Dll: &moduser32, Name: "PtInRect"}
	procSetRect                       = DllProc{Dll: &moduser32, Name: "SetRect"}
	procSetRectEmpty                  = DllProc{Dll: &moduser32, Name: "SetRectEmpty"}
	procSubtractRect                  = DllProc{Dll: &moduser32, Name: "SubtractRect"}
	procUnionRect                     = DllProc{Dll: &moduser32, Name: "UnionRect"}
	procCreateDialogParam             = DllProc{Dll: &moduser32, Name: "CreateDialogParamW"}
	procDialogBoxParam                = DllProc{Dll: &moduser32, Name: "DialogBoxParamW"}
	procGetDlgItem                    = DllProc{Dll: &moduser32, Name: "GetDlgItem"}
	procDrawIcon                      = DllProc{Dll: &moduser32, Name: "DrawIcon"}
	procClientToScreen                = DllProc{Dll: &moduser32, Name: "ClientToScreen"}
	procIsDialogMessage               = DllProc{Dll: &moduser32, Name: "IsDialogMessageW"}
	procIsWindow                      = DllProc{Dll: &moduser32, Name: "IsWindow"}
	procEndDialog                     = DllProc{Dll: &moduser32, Name: "EndDialog"}
	procPeekMessage                   = DllProc{Dll: &moduser32, Name: "PeekMessageW"}
	procTranslateAccelerator          = DllProc{Dll: &moduser32, Name: "TranslateAcceleratorW"}
	procSetWindowPos                  = DllProc{Dll: &moduser32, Name: "SetWindowPos"}
	procFillRect                      = DllProc{Dll: &moduser32, Name: "FillRect"}
	procDrawText                      = DllProc{Dll: &moduser32, Name: "DrawTextW"}
	procAddClipboardFormatListener    = DllProc{Dll: &moduser32, Name: "AddClipboardFormatListener"}
	procRemoveClipboardFormatListener = DllProc{Dll: &moduser32, Name: "RemoveClipboardFormatListener"}
	procOpenClipboard                 = DllProc{Dll: &moduser32, Name: "OpenClipboard"}
	procCloseClipboard                = DllProc{Dll: &moduser32, Name: "CloseClipboard"}
	procEnumClipboardFormats          = DllProc{Dll: &moduser32, Name: "EnumClipboardFormats"}
	procGetClipboardData              = DllProc{Dll: &moduser32, Name: "GetClipboardData"}
	procSetClipboardData              = DllProc{Dll: &moduser32, Name: "SetClipboardData"}
	procEmptyClipboard                = DllProc{Dll: &moduser32, Name: "EmptyClipboard"}
	procGetClipboardFormatName        = DllProc{Dll: &moduser32, Name: "GetClipboardFormatNameW"}
	procIsClipboardFormatAvailable    = DllProc{Dll: &moduser32, Name: "IsClipboardFormatAvailable"}
	procBeginPaint                    = DllProc{Dll: &moduser32, Name: "BeginPaint"}
	procEndPaint                      = DllProc{Dll: &moduser32, Name: "EndPaint"}
	procGetKeyboardState              = DllProc{Dll: &moduser32, Name: "GetKeyboardState"}
	procMapVirtualKey                 = DllProc{Dll: &moduser32, Name: "MapVirtualKeyExW"}
	procGetAsyncKeyState              = DllProc{Dll: &moduser32, Name: "GetAsyncKeyState"}
	procToAscii                       = DllProc{Dll: &moduser32, Name: "ToAscii"}
	procSwapMouseButton               = DllProc{Dll: &moduser32, Name: "SwapMouseButton"}
	procGetCursorPos                  = DllProc{Dll: &moduser32, Name: "GetCursorPos"}
	procSetCursorPos                  = DllProc{Dll: &moduser32, Name: "SetCursorPos"}
	procSetCursor                     = DllProc{Dll: &moduser32, Name: "SetCursor"}
	procCreateIcon                    = DllProc{Dll: &moduser32, Name: "CreateIcon"}
	procDestroyIcon                   = DllProc{Dll: &moduser32, Name: "DestroyIcon"}
	procMonitorFromPoint              = DllProc{Dll: &moduser32, Name: "MonitorFromPoint"}
	procMonitorFromRect               = DllProc{Dll: &moduser32, Name: "MonitorFromRect"}
	procMonitorFromWindow             = DllProc{Dll: &moduser32, Name: "MonitorFromWindow"}
	procGetMonitorInfo                = DllProc{Dll: &moduser32, Name: "GetMonitorInfoW"}
	procEnumDisplayMonitors           = DllProc{Dll: &moduser32, Name: "EnumDisplayMonitors"}
	procEnumDisplaySettingsEx         = DllProc{Dll: &moduser32, Name: "EnumDisplaySettingsExW"}
	procChangeDisplaySettingsEx       = DllProc{Dll: &moduser32, Name: "ChangeDisplaySettingsExW"}
	procSendInput                     = DllProc{Dll: &moduser32, Name: "SendInput"}
	procSetWindowsHookEx              = DllProc{Dll: &moduser32, Name: "SetWindowsHookExW"}
	procUnhookWindowsHookEx           = DllProc{Dll: &moduser32, Name: "UnhookWindowsHookEx"}
	procCallNextHookEx                = DllProc{Dll: &moduser32, Name: "CallNextHookEx"}
	procSetForegroundWindow           = DllProc{Dll: &moduser32, Name: "SetForegroundWindow"}
	procFindWindowW                   = DllProc{Dll: &moduser32, Name: "FindWindowW"}
	procFindWindowExW                 = DllProc{Dll: &moduser32, Name: "FindWindowExW"}
	procGetClassName                  = DllProc{Dll: &moduser32, Name: "GetClassNameW"}
	procEnumChildWindows              = DllProc{Dll: &moduser32, Name: "EnumChildWindows"}
	procSetTimer                      = DllProc{Dll: &moduser32, Name: "SetTimer"}
	procKillTimer                     = DllProc{Dll: &moduser32, Name: "KillTimer"}
	procRedrawWindow                  = DllProc{Dll: &moduser32, Name: "RedrawWindow"}
	procRegisterWindowMessage         = DllProc{Dll: &moduser32, Name: "RegisterWindowMessageW"}
)

func MustRegisterClassEx(wndClassEx *WNDCLASSEX) ATOM {
	ret, _, err := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wndClassEx)))
	if IsSyscallErrorHappen(err) {
		panic("MustRegisterClassEx fail " + SyscallErrorToMsg(err))
	}
	return ATOM(ret)
}

func MustLoadIcon(instance HINSTANCE, iconName uintptr) HICON {
	ret, _, err := procLoadIcon.Call(
		uintptr(instance),
		iconName)
	if IsSyscallErrorHappen(err) {
		panic("MustLoadIcon fail " + SyscallErrorToMsg(err))
	}
	return HICON(ret)

}

func MustLoadCursor(instance HINSTANCE, cursorName uintptr) HCURSOR {
	ret, _, err := procLoadCursor.Call(
		uintptr(instance),
		uintptr(cursorName))
	if IsSyscallErrorHappen(err) {
		panic("MustLoadCursor fail " + SyscallErrorToMsg(err))
	}
	return HCURSOR(ret)

}

func GetClassNameW(hwnd HWND) string {
	buf := make([]uint16, 255)
	procGetClassName.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(255))

	return syscall.UTF16ToString(buf)
}

func SetForegroundWindow(hwnd HWND) (isForeground bool) {
	r1, _, _ := procSetForegroundWindow.CallErrorMsg(
		uintptr(hwnd))
	return BOOLToBoolFromUintptr(r1)
}

func ShowWindow(hwnd HWND, cmdshow int32) bool {
	ret, _, _ := procShowWindow.Call(
		uintptr(hwnd),
		uintptr(cmdshow))

	return ret != 0

}

func UpdateWindow(hwnd HWND) (errMsg string) {
	_, _, err := procUpdateWindow.Call(
		uintptr(hwnd))
	if IsSyscallErrorHappen(err) {
		return "UpdateWindow " + SyscallErrorToMsg(err)
	}
	return ""
}

func MustCreateWindowEx(exStyle uint,
	className, windowName *uint16,
	style uint32,
	x, y, width, height int32,
	parent HWND, menu HMENU,
	instance HINSTANCE, param uintptr) HWND {
	ret, _, errMsg := procCreateWindowEx.CallErrorMsg(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		uintptr(param))
	if errMsg != "" {
		panic(errMsg)
	}
	return HWND(ret)
}

func FindWindowExW(hwndParent, hwndChildAfter HWND, className, windowName *uint16) HWND {
	ret, _, _ := procFindWindowExW.Call(
		uintptr(hwndParent),
		uintptr(hwndChildAfter),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)))

	return HWND(ret)
}

func FindWindowW(className, windowName *uint16) HWND {
	ret, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)))

	return HWND(ret)
}

func EnumChildWindows(hWndParent HWND, lpEnumFunc WNDENUMPROC, lParam LPARAM) bool {
	ret, _, _ := procEnumChildWindows.Call(
		uintptr(hWndParent),
		uintptr(syscall.NewCallback(lpEnumFunc)),
		uintptr(lParam),
	)

	return ret != 0
}

func AdjustWindowRectEx(rect *RECT, style uint, menu bool, exStyle uint) bool {
	ret, _, _ := procAdjustWindowRectEx.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(style),
		uintptr(BoolToBOOL(menu)),
		uintptr(exStyle))

	return ret != 0
}

func MustAdjustWindowRect(rect *RECT, style uint32, menu bool) {
	ret, _, err := procAdjustWindowRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(style),
		uintptr(BoolToBOOL(menu)))
	if ret == 0 {
		panic("MustAdjustWindowRect fail " + SyscallErrorToMsg(err))
	}
	return
}

func DestroyWindow(hwnd HWND) (errMsg string) {
	_, _, err := procDestroyWindow.Call(
		uintptr(hwnd))
	if IsSyscallErrorHappen(err) {
		return "DestroyWindow fail " + SyscallErrorToMsg(err)
	}
	return ""
}

func MustDestroyWindow(hwnd HWND) {
	errMsg := DestroyWindow(hwnd)
	if errMsg != "" {
		panic(errMsg)
	}
}

func DefWindowProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProc.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func DefDlgProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefDlgProc.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func PostQuitMessage(exitCode int) {
	procPostQuitMessage.Call(
		uintptr(exitCode))
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))

	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret != 0

}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret

}

func SendMessage(hwnd HWND, msg uint32, wParam, lParam uintptr) (ret uintptr, errMsg string) {
	ret, _, errMsg = procSendMessage.CallErrorMsg(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)
	if errMsg != "" {
		return ret, errMsg
	}
	return ret, ""
}

func MustSendMessage(hwnd HWND, msg uint32, wParam, lParam uintptr) (ret uintptr) {
	ret, errMsg := SendMessage(hwnd, msg, wParam, lParam)
	if errMsg != "" {
		panic(errMsg)
	}
	return ret
}

func SendMessageTimeout(hwnd HWND, msg uint32, wParam, lParam uintptr, fuFlags, uTimeout uint32) uintptr {
	ret, _, _ := procSendMessageTimeout.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
		uintptr(fuFlags),
		uintptr(uTimeout))

	return ret
}

const SC_RESTORE = 0xF120

func PostMessage(hwnd HWND, msg uint32, wParam, lParam uintptr) (errMsg string) {
	_, _, err := procPostMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)
	if IsSyscallErrorHappen(err) {
		return "PostMessage fail " + SyscallErrorToMsg(err)
	}
	return ""
}

func WaitMessage() bool {
	ret, _, _ := procWaitMessage.Call()
	return ret != 0
}

func SetWindowText(hwnd HWND, text string) (err error) {
	_, _, err = procSetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
	return err
}

func MustSetWindowText(hwnd HWND, text string) {
	err := SetWindowText(hwnd, text)
	if err != nil {
		panic("MustSetWindowText fail f7vwt8g28b " + err.Error())
	}
}

func GetWindowTextLength(hwnd HWND) (retI int, errMsg string) {
	ret, _, errMsg := procGetWindowTextLength.CallErrorMsg(
		uintptr(hwnd))

	return int(ret), errMsg
}

func GetWindowText(hwnd HWND) (ret string, errMsg string) {
	textLen, errMsg := GetWindowTextLength(hwnd)
	if errMsg != "" {
		return "", errMsg
	}
	textLen += 1
	buf := make([]uint16, textLen)
	_, _, errMsg = procGetWindowText.CallErrorMsg(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))
	if errMsg != "" {
		return "", errMsg
	}
	return syscall.UTF16ToString(buf), ""
}

func MustGetWindowText(hwnd HWND) (ret string) {
	ret, errMsg := GetWindowText(hwnd)
	if errMsg != "" {
		panic(errMsg)
	}
	return ret
}

func GetWindowRect(hwnd HWND) *RECT {
	var rect RECT
	procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	return &rect
}

func MoveWindow(hwnd HWND, x, y, width, height int, repaint bool) bool {
	ret, _, _ := procMoveWindow.Call(
		uintptr(hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(BoolToBOOL(repaint)))

	return ret != 0

}

func ScreenToClient(hwnd HWND, x, y int) (X, Y int, ok bool) {
	pt := POINT{X: int32(x), Y: int32(y)}
	ret, _, _ := procScreenToClient.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&pt)))

	return int(pt.X), int(pt.Y), ret != 0
}

func CallWindowProc(preWndProc uintptr, hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procCallWindowProc.Call(
		preWndProc,
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func SetWindowLong(hwnd HWND, index int, value uint32) uint32 {
	ret, _, _ := procSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index),
		uintptr(value))

	return uint32(ret)
}

func SetWindowLongPtr(hwnd HWND, index int, value uintptr) uintptr {
	ret, _, _ := procSetWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(index),
		value)

	return ret
}

func GetWindowLong(hwnd HWND, index int) int32 {
	ret, _, _ := procGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index))

	return int32(ret)
}

func GetWindowLongPtr(hwnd HWND, index int) uintptr {
	ret, _, _ := procGetWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(index))

	return ret
}

func EnableWindow(hwnd HWND, b bool) bool {
	ret, _, _ := procEnableWindow.Call(
		uintptr(hwnd),
		uintptr(BoolToBOOL(b)))
	return ret != 0
}

func IsWindowEnabled(hwnd HWND) bool {
	ret, _, _ := procIsWindowEnabled.Call(
		uintptr(hwnd))

	return ret != 0
}

func IsWindowVisible(hwnd HWND) bool {
	ret, _, _ := procIsWindowVisible.Call(
		uintptr(hwnd))

	return ret != 0
}

func SetFocus(hwnd HWND) (oldHwnd HWND, errMsg string) {
	ret, _, err := procSetFocus.Call(
		uintptr(hwnd))
	if IsSyscallErrorHappen(err) {
		return 0, "SetFocus failed " + strconv.Itoa(int(hwnd)) + " " + SyscallErrorToMsg(err)
	}
	return HWND(ret), ""
}

func MustSetFocus(hwnd HWND) (oldHwnd HWND) {
	oldHwnd, errMsg := SetFocus(hwnd)
	if errMsg != "" {
		panic(errMsg)
	}
	return oldHwnd
}

func InvalidateRect(hwnd HWND, rect *RECT, erase bool) bool {
	ret, _, _ := procInvalidateRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
		uintptr(BoolToBOOL(erase)))

	return ret != 0
}

func GetClientRect(hwnd HWND, rect *RECT) (errMsg string) {
	_, _, errMsg = procGetClientRect.CallErrorMsg(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)))
	return errMsg
}

func MustGetClientRect(hwnd HWND) *RECT {
	var rect RECT
	ret, _, err := procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		panic("GetClientRect failed " + strconv.Itoa(int(hwnd)) + " " + SyscallErrorToMsg(err))
	}

	return &rect
}

func MustGetDC(hwnd HWND) HDC {
	ret, _, err := procGetDC.Call(
		uintptr(hwnd))
	if ret == 0 {
		panic("MustGetDC fail " + SyscallErrorToMsg(err))
	}
	return HDC(ret)
}

func MustReleaseDC(hwnd HWND, hDC HDC) {
	_, _, err := procReleaseDC.Call(
		uintptr(hwnd),
		uintptr(hDC))
	if IsSyscallErrorHappen(err) {
		panic("ReleaseDC " + SyscallErrorToMsg(err))
	}
	return
}

func SetCapture(hwnd HWND) HWND {
	ret, _, _ := procSetCapture.Call(
		uintptr(hwnd))

	return HWND(ret)
}

func ReleaseCapture() bool {
	ret, _, _ := procReleaseCapture.Call()

	return ret != 0
}

func GetWindowThreadProcessId(hwnd HWND, pProcessId *uint32) (threadId uint32) {
	ret, _, _ := procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pProcessId)))

	return uint32(ret)
}

func MessageBox(hwnd HWND, title, caption string, flags uint32) (ret int32, errMsg string) {
	r1, _, errMsg := procMessageBox.CallErrorMsg(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(flags))

	return int32(r1), errMsg
}

func MustMessageBox(hwnd HWND, title, caption string, flags uint32) (ret int32) {
	ret, errMsg := MessageBox(hwnd, title, caption, flags)
	if errMsg != "" {
		panic(errMsg)
	}
	return ret
}

func GetSystemMetrics(index int) int {
	ret, _, _ := procGetSystemMetrics.Call(
		uintptr(index))

	return int(ret)
}

func CopyRect(dst, src *RECT) bool {
	ret, _, _ := procCopyRect.Call(
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src)))

	return ret != 0
}

func EqualRect(rect1, rect2 *RECT) bool {
	ret, _, _ := procEqualRect.Call(
		uintptr(unsafe.Pointer(rect1)),
		uintptr(unsafe.Pointer(rect2)))

	return ret != 0
}

func InflateRect(rect *RECT, dx, dy int) bool {
	ret, _, _ := procInflateRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(dx),
		uintptr(dy))

	return ret != 0
}

func IntersectRect(dst, src1, src2 *RECT) bool {
	ret, _, _ := procIntersectRect.Call(
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src1)),
		uintptr(unsafe.Pointer(src2)))

	return ret != 0
}

func IsRectEmpty(rect *RECT) bool {
	ret, _, _ := procIsRectEmpty.Call(
		uintptr(unsafe.Pointer(rect)))

	return ret != 0
}

func OffsetRect(rect *RECT, dx, dy int) bool {
	ret, _, _ := procOffsetRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(dx),
		uintptr(dy))

	return ret != 0
}

func PtInRect(rect *RECT, x, y int) bool {
	pt := POINT{X: int32(x), Y: int32(y)}
	ret, _, _ := procPtInRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(unsafe.Pointer(&pt)))

	return ret != 0
}

func SetRect(rect *RECT, left, top, right, bottom int) bool {
	ret, _, _ := procSetRect.Call(
		uintptr(unsafe.Pointer(rect)),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom))

	return ret != 0
}

func SetRectEmpty(rect *RECT) bool {
	ret, _, _ := procSetRectEmpty.Call(
		uintptr(unsafe.Pointer(rect)))

	return ret != 0
}

func SubtractRect(dst, src1, src2 *RECT) bool {
	ret, _, _ := procSubtractRect.Call(
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src1)),
		uintptr(unsafe.Pointer(src2)))

	return ret != 0
}

func UnionRect(dst, src1, src2 *RECT) bool {
	ret, _, _ := procUnionRect.Call(
		uintptr(unsafe.Pointer(dst)),
		uintptr(unsafe.Pointer(src1)),
		uintptr(unsafe.Pointer(src2)))

	return ret != 0
}

func CreateDialog(hInstance HINSTANCE, lpTemplate *uint16, hWndParent HWND, lpDialogProc uintptr) HWND {
	ret, _, _ := procCreateDialogParam.Call(
		uintptr(hInstance),
		uintptr(unsafe.Pointer(lpTemplate)),
		uintptr(hWndParent),
		lpDialogProc,
		0)

	return HWND(ret)
}

func DialogBox(hInstance HINSTANCE, lpTemplateName *uint16, hWndParent HWND, lpDialogProc uintptr) int {
	ret, _, _ := procDialogBoxParam.Call(
		uintptr(hInstance),
		uintptr(unsafe.Pointer(lpTemplateName)),
		uintptr(hWndParent),
		lpDialogProc,
		0)

	return int(ret)
}

func GetDlgItem(hDlg HWND, nIDDlgItem int) HWND {
	ret, _, _ := procGetDlgItem.Call(
		uintptr(unsafe.Pointer(hDlg)),
		uintptr(nIDDlgItem))

	return HWND(ret)
}

func DrawIcon(hDC HDC, x, y int, hIcon HICON) bool {
	ret, _, _ := procDrawIcon.Call(
		uintptr(unsafe.Pointer(hDC)),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(hIcon)))

	return ret != 0
}

func ClientToScreen(hwnd HWND, x, y int) (int, int) {
	pt := POINT{X: int32(x), Y: int32(y)}

	procClientToScreen.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&pt)))

	return int(pt.X), int(pt.Y)
}

func IsDialogMessage(hwnd HWND, msg *MSG) bool {
	ret, _, _ := procIsDialogMessage.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(msg)))

	return ret != 0
}

func IsWindow(hwnd HWND) bool {
	ret, _, _ := procIsWindow.Call(
		uintptr(hwnd))

	return ret != 0
}

func EndDialog(hwnd HWND, nResult uintptr) bool {
	ret, _, _ := procEndDialog.Call(
		uintptr(hwnd),
		nResult)

	return ret != 0
}

func PeekMessage(lpMsg *MSG, hwnd HWND, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	ret, _, _ := procPeekMessage.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		uintptr(wRemoveMsg))

	return ret != 0
}

func TranslateAccelerator(hwnd HWND, hAccTable HACCEL, lpMsg *MSG) bool {
	ret, _, _ := procTranslateAccelerator.Call(
		uintptr(hwnd),
		uintptr(hAccTable),
		uintptr(unsafe.Pointer(lpMsg)))

	return ret != 0
}

func SetWindowPos(hwnd, hWndInsertAfter HWND, x, y, cx, cy int, uFlags uint) (errMsg string) {
	_, _, errMsg = procSetWindowPos.CallErrorMsg(
		uintptr(hwnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags))
	return errMsg
}

func MustSetWindowPos(hwnd, hWndInsertAfter HWND, x, y, cx, cy int, uFlags uint) {
	errMsg := SetWindowPos(hwnd, hWndInsertAfter, x, y, cx, cy, uFlags)
	if errMsg != "" {
		panic(errMsg)
	}
	return
}

func MustFillRect(hDC HDC, lprc *RECT, hbr HBRUSH) {
	_, _, err := procFillRect.Call(
		uintptr(hDC),
		uintptr(unsafe.Pointer(lprc)),
		uintptr(hbr))
	if IsSyscallErrorHappen(err) {
		panic("FillRect fail " + SyscallErrorToMsg(err))
	}
	return
}

func MustDrawText(hDC HDC, text string, uCount int, lpRect *RECT, uFormat uint) int {
	ret, _, err := procDrawText.Call(
		uintptr(hDC),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(uCount),
		uintptr(unsafe.Pointer(lpRect)),
		uintptr(uFormat))
	if ret == 0 {
		panic("DrawText fail " + SyscallErrorToMsg(err))
	}
	return int(ret)
}

func AddClipboardFormatListener(hwnd HWND) bool {
	ret, _, _ := procAddClipboardFormatListener.Call(
		uintptr(hwnd))
	return ret != 0
}

func RemoveClipboardFormatListener(hwnd HWND) bool {
	ret, _, _ := procRemoveClipboardFormatListener.Call(
		uintptr(hwnd))
	return ret != 0
}

func OpenClipboard(hWndNewOwner HWND) bool {
	ret, _, _ := procOpenClipboard.Call(
		uintptr(hWndNewOwner))
	return ret != 0
}

func CloseClipboard() bool {
	ret, _, _ := procCloseClipboard.Call()
	return ret != 0
}

func EnumClipboardFormats(format uint) uint {
	ret, _, _ := procEnumClipboardFormats.Call(
		uintptr(format))
	return uint(ret)
}

func GetClipboardData(uFormat uint) HANDLE {
	ret, _, _ := procGetClipboardData.Call(
		uintptr(uFormat))
	return HANDLE(ret)
}

func SetClipboardData(uFormat uint, hMem HANDLE) HANDLE {
	ret, _, _ := procSetClipboardData.Call(
		uintptr(uFormat),
		uintptr(hMem))
	return HANDLE(ret)
}

func EmptyClipboard() bool {
	ret, _, _ := procEmptyClipboard.Call()
	return ret != 0
}

func GetClipboardFormatName(format uint) (string, bool) {
	cchMaxCount := 255
	buf := make([]uint16, cchMaxCount)
	ret, _, _ := procGetClipboardFormatName.Call(
		uintptr(format),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(cchMaxCount))

	if ret > 0 {
		return syscall.UTF16ToString(buf), true
	}

	return "Requested format does not exist or is predefined", false
}

func IsClipboardFormatAvailable(format uint) bool {
	ret, _, _ := procIsClipboardFormatAvailable.Call(uintptr(format))
	return ret != 0
}

func MustBeginPaint(hwnd HWND, paint *PAINTSTRUCT) HDC {
	ret, _, err := procBeginPaint.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(paint)))
	if ret == 0 {
		panic("[MustBeginPaint] ret==0 " + SyscallErrorToMsg(err))
	}
	return HDC(ret)
}

func EndPaint(hwnd HWND, paint *PAINTSTRUCT) {
	procEndPaint.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(paint)))
}

func GetKeyboardState(lpKeyState *[]byte) bool {
	ret, _, _ := procGetKeyboardState.Call(
		uintptr(unsafe.Pointer(&(*lpKeyState)[0])))
	return ret != 0
}

func MapVirtualKeyEx(uCode, uMapType uint, dwhkl HKL) uint {
	ret, _, _ := procMapVirtualKey.Call(
		uintptr(uCode),
		uintptr(uMapType),
		uintptr(dwhkl))
	return uint(ret)
}

func GetAsyncKeyState(vKey int) uint16 {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return uint16(ret)
}

func ToAscii(uVirtKey, uScanCode uint, lpKeyState *byte, lpChar *uint16, uFlags uint) int {
	ret, _, _ := procToAscii.Call(
		uintptr(uVirtKey),
		uintptr(uScanCode),
		uintptr(unsafe.Pointer(lpKeyState)),
		uintptr(unsafe.Pointer(lpChar)),
		uintptr(uFlags))
	return int(ret)
}

func SwapMouseButton(fSwap bool) bool {
	ret, _, _ := procSwapMouseButton.Call(
		uintptr(BoolToBOOL(fSwap)))
	return ret != 0
}

func MustGetCursorPos(pt *POINT) {
	_, _, errMsg := procGetCursorPos.CallErrorMsg(uintptr(unsafe.Pointer(pt)))
	if errMsg != "" {
		panic(errMsg)
	}

}

func SetCursorPos(x, y int) bool {
	ret, _, _ := procSetCursorPos.Call(
		uintptr(x),
		uintptr(y),
	)
	return ret != 0
}

func SetCursor(cursor HCURSOR) HCURSOR {
	ret, _, _ := procSetCursor.Call(
		uintptr(cursor),
	)
	return HCURSOR(ret)
}

func CreateIcon(instance HINSTANCE, nWidth, nHeight int, cPlanes, cBitsPerPixel byte, ANDbits, XORbits *byte) HICON {
	ret, _, _ := procCreateIcon.Call(
		uintptr(instance),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(cPlanes),
		uintptr(cBitsPerPixel),
		uintptr(unsafe.Pointer(ANDbits)),
		uintptr(unsafe.Pointer(XORbits)),
	)
	return HICON(ret)
}

func DestroyIcon(icon HICON) bool {
	ret, _, _ := procDestroyIcon.Call(
		uintptr(icon),
	)
	return ret != 0
}

func MonitorFromPoint(x, y int, dwFlags uint32) HMONITOR {
	ret, _, _ := procMonitorFromPoint.Call(
		uintptr(x),
		uintptr(y),
		uintptr(dwFlags),
	)
	return HMONITOR(ret)
}

func MonitorFromRect(rc *RECT, dwFlags uint32) HMONITOR {
	ret, _, _ := procMonitorFromRect.Call(
		uintptr(unsafe.Pointer(rc)),
		uintptr(dwFlags),
	)
	return HMONITOR(ret)
}

func MonitorFromWindow(hwnd HWND, dwFlags uint32) HMONITOR {
	ret, _, _ := procMonitorFromWindow.Call(
		uintptr(hwnd),
		uintptr(dwFlags),
	)
	return HMONITOR(ret)
}

func GetMonitorInfo(hMonitor HMONITOR, lmpi *MONITORINFO) bool {
	ret, _, _ := procGetMonitorInfo.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lmpi)),
	)
	return ret != 0
}

func EnumDisplayMonitors(hdc HDC, clip *RECT, fnEnum, dwData uintptr) bool {
	ret, _, _ := procEnumDisplayMonitors.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(clip)),
		fnEnum,
		dwData,
	)
	return ret != 0
}

func EnumDisplaySettingsEx(szDeviceName *uint16, iModeNum uint32, devMode *DEVMODE, dwFlags uint32) bool {
	ret, _, _ := procEnumDisplaySettingsEx.Call(
		uintptr(unsafe.Pointer(szDeviceName)),
		uintptr(iModeNum),
		uintptr(unsafe.Pointer(devMode)),
		uintptr(dwFlags),
	)
	return ret != 0
}

func ChangeDisplaySettingsEx(szDeviceName *uint16, devMode *DEVMODE, hwnd HWND, dwFlags uint32, lParam uintptr) int32 {
	ret, _, _ := procChangeDisplaySettingsEx.Call(
		uintptr(unsafe.Pointer(szDeviceName)),
		uintptr(unsafe.Pointer(devMode)),
		uintptr(hwnd),
		uintptr(dwFlags),
		lParam,
	)
	return int32(ret)
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func SetTimer(hwnd HWND, nIDEvent uint32, uElapse uint32, lpTimerProc uintptr) uintptr {
	ret, _, _ := procSetTimer.Call(
		uintptr(hwnd),
		uintptr(nIDEvent),
		uintptr(uElapse),
		lpTimerProc,
	)
	return ret
}

func KillTimer(hwnd HWND, nIDEvent uint32) bool {
	ret, _, _ := procKillTimer.Call(
		uintptr(hwnd),
		uintptr(nIDEvent),
	)
	return ret != 0
}

func MustRedrawWindow(hWnd HWND, lpRect *RECT, hrgnUpdate HRGN, flag uint32) {
	errMsg := RedrawWindow(hWnd, lpRect, hrgnUpdate, flag)
	if errMsg != "" {
		panic(errMsg)
	}
	return
}

func RedrawWindow(hWnd HWND, lpRect *RECT, hrgnUpdate HRGN, flag uint32) (errMsg string) {
	_, _, errMsg = procRedrawWindow.CallErrorMsg(
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lpRect)),
		uintptr(hrgnUpdate),
		uintptr(flag),
	)
	return errMsg
}

func MustRegisterWindowMessage(s string) uint32 {
	lpString := syscall.StringToUTF16Ptr(s)
	ret, _, err := procRegisterWindowMessage.Call(
		uintptr(unsafe.Pointer(lpString)),
	)
	if IsSyscallErrorHappen(err) {
		panic(SyscallErrorToMsg(err))
	}
	return uint32(ret)
}

var procSystemParametersInfo = DllProc{Dll: &moduser32, Name: "SystemParametersInfoW"}

const SPI_GETWORKAREA = 0x0030
const SPI_SETFOREGROUNDLOCKTIMEOUT = 0x2001
const SPIF_SENDWININICHANGE = 0x0002
const SPIF_UPDATEINIFILE = 0x0001

func SystemParametersInfo(uiAction uint32, uiParam uint32, pvParam uintptr, fWinIni uint32) (errMsg string) {
	_, _, errMsg = procSystemParametersInfo.CallErrorMsg(
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(pvParam),
		uintptr(fWinIni),
	)
	return errMsg
}

func MustSystemParametersInfo(uiAction uint32, uiParam uint32, pvParam uintptr, fWinIni uint32) {
	errMsg := SystemParametersInfo(uiAction, uiParam, pvParam, fWinIni)
	if errMsg != "" {
		panic(errMsg)
	}
}

var procPostThreadMessage = DllProc{Dll: &moduser32, Name: "PostThreadMessageW"}

func PostThreadMessage(idThread uint32, msg uint32, wParam, lParam uintptr) (errMsg string) {
	_, _, errMsg = procPostThreadMessage.CallErrorMsg(
		uintptr(idThread),
		uintptr(msg),
		wParam,
		lParam)
	return errMsg
}
func MustPostThreadMessage(idThread uint32, msg uint32, wParam, lParam uintptr) {
	errMsg := PostThreadMessage(idThread, msg, wParam, lParam)
	if errMsg != "" {
		panic(errMsg)
	}
	return
}

var procSetActiveWindow = DllProc{Dll: &moduser32, Name: "SetActiveWindow"}

func SetActiveWindow(hwnd HWND) (originHwnd HWND, errMsg string) {
	r1, _, errMsg := procSetActiveWindow.CallErrorMsg(
		uintptr(hwnd))
	return HWND(r1), errMsg
}

func MustSetActiveWindow(hwnd HWND) (originHwnd HWND) {
	originHwnd, errMsg := SetActiveWindow(hwnd)
	if errMsg != "" {
		panic(errMsg)
	}
	return originHwnd
}

var procBringWindowToTop = DllProc{Dll: &moduser32, Name: "BringWindowToTop"}

func BringWindowToTop(hwnd HWND) (errMsg string) {
	_, _, err := procBringWindowToTop.Call(
		uintptr(hwnd))
	if IsSyscallErrorHappen(err) {
		return "BringWindowToTop fail " + SyscallErrorToMsg(err)
	}
	return ""
}

type ICONINFO struct {
	FIcon    BOOL
	XHotspot uint32
	YHotspot uint32
	HbmMask  HBITMAP
	HbmColor HBITMAP
}

var procCreateIconIndirect = DllProc{Dll: &moduser32, Name: "CreateIconIndirect"}

func MustCreateIconIndirect(piconinfo *ICONINFO) HICON {
	ret, _, errMsg := procCreateIconIndirect.CallErrorMsg(uintptr(unsafe.Pointer(piconinfo)))
	if errMsg != "" {
		panic(errMsg)
	}
	return HICON(ret)
}

var procGetForegroundWindow = DllProc{Dll: &moduser32, Name: "GetForegroundWindow"}

func GetForegroundWindow() HWND {
	ret, _, _ := procGetForegroundWindow.Call()
	return HWND(ret)
}

var procAttachThreadInput = DllProc{Dll: &moduser32, Name: "AttachThreadInput"}

func AttachThreadInput(idAttach uint32, idAttachTo uint32, fAttach bool) (errMsg string) {
	_, _, err := procAttachThreadInput.Call(
		uintptr(idAttach),
		uintptr(idAttachTo),
		uintptr(BoolToBOOL(fAttach)),
	)
	if IsSyscallErrorHappen(err) {
		return "AttachThreadInput fail" + SyscallErrorToMsg(err)
	}
	return ""
}

func MustAttachThreadInput(idAttach uint32, idAttachTo uint32, fAttach bool) {
	errMsg := AttachThreadInput(idAttach, idAttachTo, fAttach)
	if errMsg != "" {
		panic(errMsg)
	}
}

type CIEXYZ struct {
	CiexyzX, CiexyzY, CiexyzZ int32
}

type CIEXYZTRIPLE struct {
	CiexyzRed, CiexyzGreen, CiexyzBlue CIEXYZ
}

var procSwitchToThisWindow = DllProc{Dll: &moduser32, Name: "SwitchToThisWindow"}

func SwitchToThisWindow(hwnd HWND, fAltTab bool) {
	procSwitchToThisWindow.Call(
		uintptr(hwnd),
		uintptr(BoolToBOOL(fAltTab)),
	)
}

var procGetDpiForSystem = DllProc{Dll: &moduser32, Name: "GetDpiForSystem"}

func MustGetDpiForSystem() uint32 {
	r1, _, err := procGetDpiForSystem.Call()
	if err != nil {
		panic(err)
	}
	return uint32(r1)
}
