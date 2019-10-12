// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package udwW32

import (
	"syscall"
	"unsafe"
)

const (
	E_UNEXPECTED   = 0x8000FFFF
	E_NOTIMPL      = 0x80004001
	E_OUTOFMEMORY  = 0x8007000E
	E_INVALIDARG   = 0x80070057
	E_NOINTERFACE  = 0x80004002
	E_POINTER      = 0x80004003
	E_HANDLE       = 0x80070006
	E_ABORT        = 0x80004004
	E_FAIL         = 0x80004005
	E_ACCESSDENIED = 0x80070005
	E_PENDING      = 0x8000000A
)

type IID syscall.GUID
type CLSID syscall.GUID
type REFIID *IID
type REFCLSID *CLSID

var modole32 = Dll{Name: "ole32.dll"}

var procCoInitializeEx = DllProc{Dll: &modole32, Name: "CoInitializeEx"}

func CoInitializeEx(coInit uintptr) HRESULT {
	ret, _, _ := procCoInitializeEx.Call(
		0,
		coInit)

	switch uint32(ret) {
	case E_INVALIDARG:
		panic("CoInitializeEx failed with E_INVALIDARG")
	case E_OUTOFMEMORY:
		panic("CoInitializeEx failed with E_OUTOFMEMORY")
	case E_UNEXPECTED:
		panic("CoInitializeEx failed with E_UNEXPECTED")
	}

	return HRESULT(ret)
}

var procCoInitialize = DllProc{Dll: &modole32, Name: "CoInitialize"}

func CoInitialize() {
	procCoInitialize.Call(0)
}

var procCoUninitialize = DllProc{Dll: &modole32, Name: "CoUninitialize"}

func CoUninitialize() {
	procCoUninitialize.Call()
}

var procCreateStreamOnHGlobal = DllProc{Dll: &modole32, Name: "CreateStreamOnHGlobal"}

func CreateStreamOnHGlobal(hGlobal HGLOBAL, fDeleteOnRelease bool) *IStream {
	stream := new(IStream)
	ret, _, _ := procCreateStreamOnHGlobal.Call(
		uintptr(hGlobal),
		uintptr(BoolToBOOL(fDeleteOnRelease)),
		uintptr(unsafe.Pointer(&stream)))

	switch uint32(ret) {
	case E_INVALIDARG:
		panic("CreateStreamOnHGlobal failed with E_INVALIDARG")
	case E_OUTOFMEMORY:
		panic("CreateStreamOnHGlobal failed with E_OUTOFMEMORY")
	case E_UNEXPECTED:
		panic("CreateStreamOnHGlobal failed with E_UNEXPECTED")
	}

	return stream
}

var procOleInitialize = DllProc{Dll: &modole32, Name: "OleInitialize"}

func OleInitialize() HRESULT {

	ret, _, _ := syscall.Syscall(procOleInitialize.Addr(), 1,
		0,
		0,
		0)

	return HRESULT(ret)
}

var procOleUninitialize = DllProc{Dll: &modole32, Name: "OleUninitialize"}

func OleUninitialize() {
	syscall.Syscall(procOleUninitialize.Addr(), 0,
		0,
		0,
		0)
}

var procCoCreateInstance = DllProc{Dll: &modole32, Name: "CoCreateInstance"}

func CoCreateInstance(rclsid REFCLSID, pUnkOuter *IUnknown, dwClsContext uint32, riid REFIID, ppv *unsafe.Pointer) HRESULT {
	ret, _, _ := syscall.Syscall6(procCoCreateInstance.Addr(), 5,
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(unsafe.Pointer(pUnkOuter)),
		uintptr(dwClsContext),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppv)),
		0)

	return HRESULT(ret)
}

type COAUTHIDENTITY struct {
	User           *uint16
	UserLength     uint32
	Domain         *uint16
	DomainLength   uint32
	Password       *uint16
	PasswordLength uint32
	Flags          uint32
}

type COAUTHINFO struct {
	dwAuthnSvc           uint32
	dwAuthzSvc           uint32
	pwszServerPrincName  *uint16
	dwAuthnLevel         uint32
	dwImpersonationLevel uint32
	pAuthIdentityData    *COAUTHIDENTITY
	dwCapabilities       uint32
}

type COSERVERINFO struct {
	dwReserved1 uint32
	pwszName    *uint16
	pAuthInfo   *COAUTHINFO
	dwReserved2 uint32
}

var (
	IID_IClassFactory             = IID{0x00000001, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IConnectionPointContainer = IID{0xB196B284, 0xBAB4, 0x101A, [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07}}
	IID_IOleClientSite            = IID{0x00000118, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IOleInPlaceObject         = IID{0x00000113, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IOleInPlaceSite           = IID{0x00000119, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IOleObject                = IID{0x00000112, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	IID_IUnknown                  = IID{0x00000000, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
)

var procCoGetClassObject = DllProc{Dll: &modole32, Name: "CoGetClassObject"}

func CoGetClassObject(rclsid REFCLSID, dwClsContext uint32, pServerInfo *COSERVERINFO, riid REFIID, ppv *unsafe.Pointer) HRESULT {

	ret, _, _ := syscall.Syscall6(procCoGetClassObject.Addr(), 5,
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(dwClsContext),
		uintptr(unsafe.Pointer(pServerInfo)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppv)),
		0)

	return HRESULT(ret)
}

type IClassFactoryVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	CreateInstance uintptr
	LockServer     uintptr
}

type IClassFactory struct {
	LpVtbl *IClassFactoryVtbl
}

func (cf *IClassFactory) Release() uint32 {
	ret, _, _ := syscall.Syscall(cf.LpVtbl.Release, 1,
		uintptr(unsafe.Pointer(cf)),
		0,
		0)

	return uint32(ret)
}

func (cf *IClassFactory) CreateInstance(pUnkOuter *IUnknown, riid REFIID, ppvObject *unsafe.Pointer) HRESULT {
	ret, _, _ := syscall.Syscall6(cf.LpVtbl.CreateInstance, 4,
		uintptr(unsafe.Pointer(cf)),
		uintptr(unsafe.Pointer(pUnkOuter)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)),
		0,
		0)

	return HRESULT(ret)
}

type IOleClientSiteVtbl struct {
	QueryInterface         uintptr
	AddRef                 uintptr
	Release                uintptr
	SaveObject             uintptr
	GetMoniker             uintptr
	GetContainer           uintptr
	ShowObject             uintptr
	OnShowWindow           uintptr
	RequestNewObjectLayout uintptr
}

type IOleClientSite struct {
	LpVtbl *IOleClientSiteVtbl
}

type IOleObjectVtbl struct {
	QueryInterface   uintptr
	AddRef           uintptr
	Release          uintptr
	SetClientSite    uintptr
	GetClientSite    uintptr
	SetHostNames     uintptr
	Close            uintptr
	SetMoniker       uintptr
	GetMoniker       uintptr
	InitFromData     uintptr
	GetClipboardData uintptr
	DoVerb           uintptr
	EnumVerbs        uintptr
	Update           uintptr
	IsUpToDate       uintptr
	GetUserClassID   uintptr
	GetUserType      uintptr
	SetExtent        uintptr
	GetExtent        uintptr
	Advise           uintptr
	Unadvise         uintptr
	EnumAdvise       uintptr
	GetMiscStatus    uintptr
	SetColorScheme   uintptr
}

type IOleObject struct {
	LpVtbl *IOleObjectVtbl
}

func (obj *IOleObject) QueryInterface(riid REFIID, ppvObject *unsafe.Pointer) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.QueryInterface, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)))

	return HRESULT(ret)
}

func (obj *IOleObject) Release() uint32 {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.Release, 1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)

	return uint32(ret)
}

func (obj *IOleObject) SetClientSite(pClientSite *IOleClientSite) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.SetClientSite, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pClientSite)),
		0)

	return HRESULT(ret)
}

func (obj *IOleObject) SetHostNames(szContainerApp, szContainerObj *uint16) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.SetHostNames, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(szContainerApp)),
		uintptr(unsafe.Pointer(szContainerObj)))

	return HRESULT(ret)
}

func (obj *IOleObject) Close(dwSaveOption uint32) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.Close, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(dwSaveOption),
		0)

	return HRESULT(ret)
}

func (obj *IOleObject) DoVerb(iVerb int32, lpmsg *MSG, pActiveSite *IOleClientSite, lindex int32, hwndParent HWND, lprcPosRect *RECT) HRESULT {
	ret, _, _ := syscall.Syscall9(obj.LpVtbl.DoVerb, 7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(iVerb),
		uintptr(unsafe.Pointer(lpmsg)),
		uintptr(unsafe.Pointer(pActiveSite)),
		uintptr(lindex),
		uintptr(hwndParent),
		uintptr(unsafe.Pointer(lprcPosRect)),
		0,
		0)

	return HRESULT(ret)
}

func EqualREFIID(a, b REFIID) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Data1 != b.Data1 || a.Data2 != b.Data2 || a.Data3 != b.Data3 {
		return false
	}
	for i := 0; i < 8; i++ {
		if a.Data4[i] != b.Data4[i] {
			return false
		}
	}
	return true
}

type IOleInPlaceFrameVtbl struct {
	QueryInterface       uintptr
	AddRef               uintptr
	Release              uintptr
	GetWindow            uintptr
	ContextSensitiveHelp uintptr
	GetBorder            uintptr
	RequestBorderSpace   uintptr
	SetBorderSpace       uintptr
	SetActiveObject      uintptr
	InsertMenus          uintptr
	SetMenu              uintptr
	RemoveMenus          uintptr
	SetStatusText        uintptr
	EnableModeless       uintptr
	TranslateAccelerator uintptr
}

type IOleInPlaceFrame struct {
	LpVtbl *IOleInPlaceFrameVtbl
}

type IOleInPlaceSiteVtbl struct {
	QueryInterface       uintptr
	AddRef               uintptr
	Release              uintptr
	GetWindow            uintptr
	ContextSensitiveHelp uintptr
	CanInPlaceActivate   uintptr
	OnInPlaceActivate    uintptr
	OnUIActivate         uintptr
	GetWindowContext     uintptr
	Scroll               uintptr
	OnUIDeactivate       uintptr
	OnInPlaceDeactivate  uintptr
	DiscardUndoState     uintptr
	DeactivateAndUndo    uintptr
	OnPosRectChange      uintptr
}

type IOleInPlaceSite struct {
	LpVtbl *IOleInPlaceSiteVtbl
}

type OLEINPLACEFRAMEINFO struct {
	Cb            uint32
	FMDIApp       BOOL
	HwndFrame     HWND
	Haccel        HACCEL
	CAccelEntries uint32
}

type IOleInPlaceObjectVtbl struct {
	QueryInterface       uintptr
	AddRef               uintptr
	Release              uintptr
	GetWindow            uintptr
	ContextSensitiveHelp uintptr
	InPlaceDeactivate    uintptr
	UIDeactivate         uintptr
	SetObjectRects       uintptr
	ReactivateAndUndo    uintptr
}

type IOleInPlaceObject struct {
	LpVtbl *IOleInPlaceObjectVtbl
}

func (obj *IOleInPlaceObject) Release() uint32 {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.Release, 1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)

	return uint32(ret)
}

func (obj *IOleInPlaceObject) SetObjectRects(lprcPosRect, lprcClipRect *RECT) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.SetObjectRects, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(lprcPosRect)),
		uintptr(unsafe.Pointer(lprcClipRect)))

	return HRESULT(ret)
}

const (
	OLEIVERB_PRIMARY          = 0
	OLEIVERB_SHOW             = -1
	OLEIVERB_OPEN             = -2
	OLEIVERB_HIDE             = -3
	OLEIVERB_UIACTIVATE       = -4
	OLEIVERB_INPLACEACTIVATE  = -5
	OLEIVERB_DISCARDUNDOSTATE = -6
)
