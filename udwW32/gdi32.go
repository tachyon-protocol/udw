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

const (
	DIB_RGB_COLORS = 0
	DIB_PAL_COLORS = 1
)

var (
	modgdi32 = Dll{Name: "gdi32.dll"}

	procGetDeviceCaps             = DllProc{Dll: &modgdi32, Name: "GetDeviceCaps"}
	procGetCurrentObject          = DllProc{Dll: &modgdi32, Name: "GetCurrentObject"}
	procDeleteObject              = DllProc{Dll: &modgdi32, Name: "DeleteObject"}
	procCreateFontIndirect        = DllProc{Dll: &modgdi32, Name: "CreateFontIndirectW"}
	procAbortDoc                  = DllProc{Dll: &modgdi32, Name: "AbortDoc"}
	procBitBlt                    = DllProc{Dll: &modgdi32, Name: "BitBlt"}
	procPatBlt                    = DllProc{Dll: &modgdi32, Name: "PatBlt"}
	procCloseEnhMetaFile          = DllProc{Dll: &modgdi32, Name: "CloseEnhMetaFile"}
	procCopyEnhMetaFile           = DllProc{Dll: &modgdi32, Name: "CopyEnhMetaFileW"}
	procCreateBrushIndirect       = DllProc{Dll: &modgdi32, Name: "CreateBrushIndirect"}
	procCreateCompatibleDC        = DllProc{Dll: &modgdi32, Name: "CreateCompatibleDC"}
	procCreateDC                  = DllProc{Dll: &modgdi32, Name: "CreateDCW"}
	procCreateCompatibleBitmap    = DllProc{Dll: &modgdi32, Name: "CreateCompatibleBitmap"}
	procCreateDIBSection          = DllProc{Dll: &modgdi32, Name: "CreateDIBSection"}
	procCreateEnhMetaFile         = DllProc{Dll: &modgdi32, Name: "CreateEnhMetaFileW"}
	procCreateIC                  = DllProc{Dll: &modgdi32, Name: "CreateICW"}
	procDeleteDC                  = DllProc{Dll: &modgdi32, Name: "DeleteDC"}
	procDeleteEnhMetaFile         = DllProc{Dll: &modgdi32, Name: "DeleteEnhMetaFile"}
	procEllipse                   = DllProc{Dll: &modgdi32, Name: "Ellipse"}
	procEndDoc                    = DllProc{Dll: &modgdi32, Name: "EndDoc"}
	procEndPage                   = DllProc{Dll: &modgdi32, Name: "EndPage"}
	procExtCreatePen              = DllProc{Dll: &modgdi32, Name: "ExtCreatePen"}
	procGetEnhMetaFile            = DllProc{Dll: &modgdi32, Name: "GetEnhMetaFileW"}
	procGetEnhMetaFileHeader      = DllProc{Dll: &modgdi32, Name: "GetEnhMetaFileHeader"}
	procGetObject                 = DllProc{Dll: &modgdi32, Name: "GetObjectW"}
	procGetStockObject            = DllProc{Dll: &modgdi32, Name: "GetStockObject"}
	procGetTextExtentExPoint      = DllProc{Dll: &modgdi32, Name: "GetTextExtentExPointW"}
	procGetTextExtentPoint32      = DllProc{Dll: &modgdi32, Name: "GetTextExtentPoint32W"}
	procGetTextMetrics            = DllProc{Dll: &modgdi32, Name: "GetTextMetricsW"}
	procLineTo                    = DllProc{Dll: &modgdi32, Name: "LineTo"}
	procMoveToEx                  = DllProc{Dll: &modgdi32, Name: "MoveToEx"}
	procPlayEnhMetaFile           = DllProc{Dll: &modgdi32, Name: "PlayEnhMetaFile"}
	procRectangle                 = DllProc{Dll: &modgdi32, Name: "Rectangle"}
	procResetDC                   = DllProc{Dll: &modgdi32, Name: "ResetDCW"}
	procSelectObject              = DllProc{Dll: &modgdi32, Name: "SelectObject"}
	procSetBkMode                 = DllProc{Dll: &modgdi32, Name: "SetBkMode"}
	procSetBrushOrgEx             = DllProc{Dll: &modgdi32, Name: "SetBrushOrgEx"}
	procSetStretchBltMode         = DllProc{Dll: &modgdi32, Name: "SetStretchBltMode"}
	procSetTextColor              = DllProc{Dll: &modgdi32, Name: "SetTextColor"}
	procSetBkColor                = DllProc{Dll: &modgdi32, Name: "SetBkColor"}
	procStartDoc                  = DllProc{Dll: &modgdi32, Name: "StartDocW"}
	procStartPage                 = DllProc{Dll: &modgdi32, Name: "StartPage"}
	procStretchBlt                = DllProc{Dll: &modgdi32, Name: "StretchBlt"}
	procSetDIBitsToDevice         = DllProc{Dll: &modgdi32, Name: "SetDIBitsToDevice"}
	procChoosePixelFormat         = DllProc{Dll: &modgdi32, Name: "ChoosePixelFormat"}
	procDescribePixelFormat       = DllProc{Dll: &modgdi32, Name: "DescribePixelFormat"}
	procGetEnhMetaFilePixelFormat = DllProc{Dll: &modgdi32, Name: "GetEnhMetaFilePixelFormat"}
	procGetPixelFormat            = DllProc{Dll: &modgdi32, Name: "GetPixelFormat"}
	procSetPixelFormat            = DllProc{Dll: &modgdi32, Name: "SetPixelFormat"}
	procSwapBuffers               = DllProc{Dll: &modgdi32, Name: "SwapBuffers"}
)

func MustGetDeviceCaps(hdc HDC, index int) int {
	ret, _, err := procGetDeviceCaps.Call(
		uintptr(hdc),
		uintptr(index))
	if err != nil {
		panic(err)
	}
	return int(ret)
}

func GetCurrentObject(hdc HDC, uObjectType uint32) HGDIOBJ {
	ret, _, _ := procGetCurrentObject.Call(
		uintptr(hdc),
		uintptr(uObjectType))

	return HGDIOBJ(ret)
}

func MustDeleteObject(hObject HGDIOBJ) {
	ret, _, err := procDeleteObject.Call(
		uintptr(hObject))
	if ret == 0 {
		panic("DeleteObject fail " + SyscallErrorToMsg(err))
	}
	return
}

func CreateFontIndirect(logFont *LOGFONT) HFONT {
	ret, _, _ := procCreateFontIndirect.Call(
		uintptr(unsafe.Pointer(logFont)))

	return HFONT(ret)
}

func AbortDoc(hdc HDC) int {
	ret, _, _ := procAbortDoc.Call(
		uintptr(hdc))

	return int(ret)
}

func MustBitBlt(hdcDest HDC, nXDest, nYDest, nWidth, nHeight int32, hdcSrc HDC, nXSrc, nYSrc int32, dwRop uint) {
	ret, _, err := procBitBlt.Call(
		uintptr(hdcDest),
		uintptr(nXDest),
		uintptr(nYDest),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hdcSrc),
		uintptr(nXSrc),
		uintptr(nYSrc),
		uintptr(dwRop))

	if ret == 0 {
		panic("BitBlt failed " + SyscallErrorToMsg(err))
	}
}

func PatBlt(hdc HDC, nXLeft, nYLeft, nWidth, nHeight int, dwRop uint) {
	ret, _, _ := procPatBlt.Call(
		uintptr(hdc),
		uintptr(nXLeft),
		uintptr(nYLeft),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(dwRop))

	if ret == 0 {
		panic("PatBlt failed")
	}
}

func CloseEnhMetaFile(hdc HDC) HENHMETAFILE {
	ret, _, _ := procCloseEnhMetaFile.Call(
		uintptr(hdc))

	return HENHMETAFILE(ret)
}

func CopyEnhMetaFile(hemfSrc HENHMETAFILE, lpszFile *uint16) HENHMETAFILE {
	ret, _, _ := procCopyEnhMetaFile.Call(
		uintptr(hemfSrc),
		uintptr(unsafe.Pointer(lpszFile)))

	return HENHMETAFILE(ret)
}

func CreateBrushIndirect(lplb *LOGBRUSH) HBRUSH {
	ret, _, _ := procCreateBrushIndirect.Call(
		uintptr(unsafe.Pointer(lplb)))

	return HBRUSH(ret)
}

func MustCreateCompatibleDC(hdc HDC) HDC {
	ret, _, err := procCreateCompatibleDC.Call(
		uintptr(hdc))

	if ret == 0 {
		panic("Create compatible DC failed " + SyscallErrorToMsg(err))
	}

	return HDC(ret)
}

func CreateDC(lpszDriver, lpszDevice, lpszOutput *uint16, lpInitData *DEVMODE) HDC {
	ret, _, _ := procCreateDC.Call(
		uintptr(unsafe.Pointer(lpszDriver)),
		uintptr(unsafe.Pointer(lpszDevice)),
		uintptr(unsafe.Pointer(lpszOutput)),
		uintptr(unsafe.Pointer(lpInitData)))

	return HDC(ret)
}

func MustCreateCompatibleBitmap(hdc HDC, width int32, height int32) HBITMAP {
	ret, _, err := procCreateCompatibleBitmap.Call(
		uintptr(hdc),
		uintptr(width),
		uintptr(height))
	if ret == 0 || IsSyscallErrorHappen(err) {
		panic("MustCreateCompatibleBitmap " + SyscallErrorToMsg(err) + " " + strconv.Itoa(int(ret)))
	}
	return HBITMAP(ret)
}

func MustCreateDIBSection(hdc HDC, pbmi *BITMAPINFO, iUsage uint, ppvBits **byte, hSection HANDLE, dwOffset uint) HBITMAP {
	ret, _, err := procCreateDIBSection.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pbmi)),
		uintptr(iUsage),
		uintptr(unsafe.Pointer(ppvBits)),
		uintptr(hSection),
		uintptr(dwOffset))
	if ret == 0 {
		panic("MustCreateDIBSection " + SyscallErrorToMsg(err))
	}
	return HBITMAP(ret)
}

func CreateEnhMetaFile(hdcRef HDC, lpFilename *uint16, lpRect *RECT, lpDescription *uint16) HDC {
	ret, _, _ := procCreateEnhMetaFile.Call(
		uintptr(hdcRef),
		uintptr(unsafe.Pointer(lpFilename)),
		uintptr(unsafe.Pointer(lpRect)),
		uintptr(unsafe.Pointer(lpDescription)))

	return HDC(ret)
}

func CreateIC(lpszDriver, lpszDevice, lpszOutput *uint16, lpdvmInit *DEVMODE) HDC {
	ret, _, _ := procCreateIC.Call(
		uintptr(unsafe.Pointer(lpszDriver)),
		uintptr(unsafe.Pointer(lpszDevice)),
		uintptr(unsafe.Pointer(lpszOutput)),
		uintptr(unsafe.Pointer(lpdvmInit)))

	return HDC(ret)
}

func MustDeleteDC(hdc HDC) {
	ret, _, err := procDeleteDC.Call(
		uintptr(hdc))
	if ret == 0 {
		panic("DeleteDC fail " + SyscallErrorToMsg(err))
	}
	return
}

func DeleteEnhMetaFile(hemf HENHMETAFILE) bool {
	ret, _, _ := procDeleteEnhMetaFile.Call(
		uintptr(hemf))

	return ret != 0
}

func Ellipse(hdc HDC, nLeftRect, nTopRect, nRightRect, nBottomRect int) bool {
	ret, _, _ := procEllipse.Call(
		uintptr(hdc),
		uintptr(nLeftRect),
		uintptr(nTopRect),
		uintptr(nRightRect),
		uintptr(nBottomRect))

	return ret != 0
}

func EndDoc(hdc HDC) int {
	ret, _, _ := procEndDoc.Call(
		uintptr(hdc))

	return int(ret)
}

func EndPage(hdc HDC) int {
	ret, _, _ := procEndPage.Call(
		uintptr(hdc))

	return int(ret)
}

func ExtCreatePen(dwPenStyle, dwWidth uint, lplb *LOGBRUSH, dwStyleCount uint, lpStyle *uint) HPEN {
	ret, _, _ := procExtCreatePen.Call(
		uintptr(dwPenStyle),
		uintptr(dwWidth),
		uintptr(unsafe.Pointer(lplb)),
		uintptr(dwStyleCount),
		uintptr(unsafe.Pointer(lpStyle)))

	return HPEN(ret)
}

func GetEnhMetaFile(lpszMetaFile *uint16) HENHMETAFILE {
	ret, _, _ := procGetEnhMetaFile.Call(
		uintptr(unsafe.Pointer(lpszMetaFile)))

	return HENHMETAFILE(ret)
}

func GetEnhMetaFileHeader(hemf HENHMETAFILE, cbBuffer uint, lpemh *ENHMETAHEADER) uint {
	ret, _, _ := procGetEnhMetaFileHeader.Call(
		uintptr(hemf),
		uintptr(cbBuffer),
		uintptr(unsafe.Pointer(lpemh)))

	return uint(ret)
}

func GetObject(hgdiobj HGDIOBJ, cbBuffer uintptr, lpvObject unsafe.Pointer) int {
	ret, _, _ := procGetObject.Call(
		uintptr(hgdiobj),
		uintptr(cbBuffer),
		uintptr(lpvObject))

	return int(ret)
}

func MustGetStockObject(fnObject int) HGDIOBJ {
	ret, _, err := procGetStockObject.Call(
		uintptr(fnObject))
	if ret == 0 {
		panic("GetStockObject fail " + SyscallErrorToMsg(err))
	}
	return HGDIOBJ(ret)
}

func GetTextExtentExPoint(hdc HDC, lpszStr *uint16, cchString, nMaxExtent int, lpnFit, alpDx *int, lpSize *SIZE) bool {
	ret, _, _ := procGetTextExtentExPoint.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpszStr)),
		uintptr(cchString),
		uintptr(nMaxExtent),
		uintptr(unsafe.Pointer(lpnFit)),
		uintptr(unsafe.Pointer(alpDx)),
		uintptr(unsafe.Pointer(lpSize)))

	return ret != 0
}

func MustGetTextExtentPoint32(hdc HDC, msg string, lpSize *SIZE) {
	stringArr := syscall.StringToUTF16(msg)
	c := len(stringArr)
	lpString := &stringArr[0]
	_, _, err := procGetTextExtentPoint32.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpString)),
		uintptr(c),
		uintptr(unsafe.Pointer(lpSize)))
	if IsSyscallErrorHappen(err) {
		panic("[MustGetTextExtentPoint32] " + SyscallErrorToMsg(err))
	}
	return
}

func GetTextMetrics(hdc HDC, lptm *TEXTMETRIC) bool {
	ret, _, _ := procGetTextMetrics.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lptm)))

	return ret != 0
}

func LineTo(hdc HDC, nXEnd, nYEnd int) bool {
	ret, _, _ := procLineTo.Call(
		uintptr(hdc),
		uintptr(nXEnd),
		uintptr(nYEnd))

	return ret != 0
}

func MoveToEx(hdc HDC, x, y int, lpPoint *POINT) bool {
	ret, _, _ := procMoveToEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpPoint)))

	return ret != 0
}

func PlayEnhMetaFile(hdc HDC, hemf HENHMETAFILE, lpRect *RECT) bool {
	ret, _, _ := procPlayEnhMetaFile.Call(
		uintptr(hdc),
		uintptr(hemf),
		uintptr(unsafe.Pointer(lpRect)))

	return ret != 0
}

func Rectangle(hdc HDC, nLeftRect, nTopRect, nRightRect, nBottomRect int) bool {
	ret, _, _ := procRectangle.Call(
		uintptr(hdc),
		uintptr(nLeftRect),
		uintptr(nTopRect),
		uintptr(nRightRect),
		uintptr(nBottomRect))

	return ret != 0
}

func ResetDC(hdc HDC, lpInitData *DEVMODE) HDC {
	ret, _, _ := procResetDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpInitData)))

	return HDC(ret)
}

func MustSelectObject(hdc HDC, hgdiobj HGDIOBJ) HGDIOBJ {
	ret, _, err := procSelectObject.Call(
		uintptr(hdc),
		uintptr(hgdiobj))

	if ret == 0 {
		panic("SelectObject failed " + SyscallErrorToMsg(err))
	}

	return HGDIOBJ(ret)
}

func SetBkMode(hdc HDC, iBkMode int) int {
	ret, _, _ := procSetBkMode.Call(
		uintptr(hdc),
		uintptr(iBkMode))

	if ret == 0 {
		panic("SetBkMode failed")
	}

	return int(ret)
}

func SetBrushOrgEx(hdc HDC, nXOrg, nYOrg int, lppt *POINT) bool {
	ret, _, _ := procSetBrushOrgEx.Call(
		uintptr(hdc),
		uintptr(nXOrg),
		uintptr(nYOrg),
		uintptr(unsafe.Pointer(lppt)))

	return ret != 0
}

func SetStretchBltMode(hdc HDC, iStretchMode int) int {
	ret, _, _ := procSetStretchBltMode.Call(
		uintptr(hdc),
		uintptr(iStretchMode))

	return int(ret)
}

func SetTextColor(hdc HDC, crColor COLORREF) COLORREF {
	ret, _, _ := procSetTextColor.Call(
		uintptr(hdc),
		uintptr(crColor))

	if ret == CLR_INVALID {
		panic("SetTextColor failed")
	}

	return COLORREF(ret)
}

func MustSetBkColor(hdc HDC, crColor COLORREF) COLORREF {
	ret, _, _ := procSetBkColor.Call(
		uintptr(hdc),
		uintptr(crColor))

	if ret == CLR_INVALID {
		panic("SetBkColor failed")
	}

	return COLORREF(ret)
}

func StartDoc(hdc HDC, lpdi *DOCINFO) int {
	ret, _, _ := procStartDoc.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpdi)))

	return int(ret)
}

func StartPage(hdc HDC) int {
	ret, _, _ := procStartPage.Call(
		uintptr(hdc))

	return int(ret)
}

func StretchBlt(hdcDest HDC, nXOriginDest, nYOriginDest, nWidthDest, nHeightDest int, hdcSrc HDC, nXOriginSrc, nYOriginSrc, nWidthSrc, nHeightSrc int, dwRop uint) {
	ret, _, _ := procStretchBlt.Call(
		uintptr(hdcDest),
		uintptr(nXOriginDest),
		uintptr(nYOriginDest),
		uintptr(nWidthDest),
		uintptr(nHeightDest),
		uintptr(hdcSrc),
		uintptr(nXOriginSrc),
		uintptr(nYOriginSrc),
		uintptr(nWidthSrc),
		uintptr(nHeightSrc),
		uintptr(dwRop))

	if ret == 0 {
		panic("StretchBlt failed")
	}
}

func SetDIBitsToDevice(hdc HDC, xDest, yDest, dwWidth, dwHeight, xSrc, ySrc int, uStartScan, cScanLines uint, lpvBits []byte, lpbmi *BITMAPINFO, fuColorUse uint) int {
	ret, _, _ := procSetDIBitsToDevice.Call(
		uintptr(hdc),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(dwWidth),
		uintptr(dwHeight),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(uStartScan),
		uintptr(cScanLines),
		uintptr(unsafe.Pointer(&lpvBits[0])),
		uintptr(unsafe.Pointer(lpbmi)),
		uintptr(fuColorUse))

	return int(ret)
}

func ChoosePixelFormat(hdc HDC, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, _ := procChoosePixelFormat.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pfd)),
	)
	return int(ret)
}

func DescribePixelFormat(hdc HDC, iPixelFormat int, nBytes uint, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, _ := procDescribePixelFormat.Call(
		uintptr(hdc),
		uintptr(iPixelFormat),
		uintptr(nBytes),
		uintptr(unsafe.Pointer(pfd)),
	)
	return int(ret)
}

func GetEnhMetaFilePixelFormat(hemf HENHMETAFILE, cbBuffer uint32, pfd *PIXELFORMATDESCRIPTOR) uint {
	ret, _, _ := procGetEnhMetaFilePixelFormat.Call(
		uintptr(hemf),
		uintptr(cbBuffer),
		uintptr(unsafe.Pointer(pfd)),
	)
	return uint(ret)
}

func GetPixelFormat(hdc HDC) int {
	ret, _, _ := procGetPixelFormat.Call(
		uintptr(hdc),
	)
	return int(ret)
}

func SetPixelFormat(hdc HDC, iPixelFormat int, pfd *PIXELFORMATDESCRIPTOR) bool {
	ret, _, _ := procSetPixelFormat.Call(
		uintptr(hdc),
		uintptr(iPixelFormat),
		uintptr(unsafe.Pointer(pfd)),
	)
	return ret == TRUE
}

func SwapBuffers(hdc HDC) bool {
	ret, _, _ := procSwapBuffers.Call(uintptr(hdc))
	return ret == TRUE
}

const ERROR_INVALID_PARAMETER = 87

var procCreateBitmap = DllProc{Dll: &modgdi32, Name: "CreateBitmap"}

func MustCreateBitmap(nWidth, nHeight int32, cPlanes, cBitsPerPel uint32, lpvBits unsafe.Pointer) HBITMAP {
	ret, _, _ := procCreateBitmap.Call(
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(cPlanes),
		uintptr(cBitsPerPel),
		uintptr(lpvBits))
	if ret == ERROR_INVALID_PARAMETER || ret == 0 {
		panic("CreateBitmap fail")
	}
	return HBITMAP(ret)
}
