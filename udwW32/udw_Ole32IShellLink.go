// +build windows

package udwW32

import (
	"strconv"
	"syscall"
	"unsafe"
)

type IShellLinkVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPath             uintptr
	GetIDList           uintptr
	SetIDList           uintptr
	GetDescription      uintptr
	SetDescription      uintptr
	GetWorkingDirectory uintptr
	SetWorkingDirectory uintptr
	GetArguments        uintptr
	SetArguments        uintptr
	GetHotkey           uintptr
	SetHotkey           uintptr
	GetShowCmd          uintptr
	SetShowCmd          uintptr
	GetIconLocation     uintptr
	SetIconLocation     uintptr
	SetRelativePath     uintptr
	Resolve             uintptr
	SetPath             uintptr
}

type IShellLink struct {
	LpVtbl *IShellLinkVtbl
}

func (obj *IShellLink) Release() uint32 {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.Release, 1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)

	return uint32(ret)
}
func (obj *IShellLink) SetDescription(path string) (errMsg string) {
	p := syscall.StringToUTF16Ptr(path)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.SetDescription, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		0)

	if hres != S_OK {
		return "[IShellLink.SetDescription] fail " + strconv.Itoa(int(hres))
	}

	return ""
}
func (obj *IShellLink) SetArguments(path string) (errMsg string) {
	p := syscall.StringToUTF16Ptr(path)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.SetArguments, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		0)

	if hres != S_OK {
		return "[IShellLink.SetArguments] fail " + strconv.Itoa(int(hres))
	}

	return ""
}
func (obj *IShellLink) GetDescription() (path string, errMsg string) {
	p := [MAX_PATH]uint16{}
	lpp := &p[0]
	hres, _, _ := syscall.Syscall(obj.LpVtbl.GetDescription, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(lpp)),
		uintptr(MAX_PATH))

	if hres != S_OK {
		return "", "[IShellLink.GetDescription] fail " + strconv.Itoa(int(hres))
	}
	path = syscall.UTF16ToString(p[:])

	return path, ""
}
func (obj *IShellLink) SetPath(path string) (errMsg string) {
	p := syscall.StringToUTF16Ptr(path)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.SetPath, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		0)

	if hres != S_OK {
		return "[IShellLink.SetPath] fail " + strconv.Itoa(int(hres))
	}

	return ""
}
func (obj *IShellLink) SetIconLocation(path string, index int) (errMsg string) {
	p := syscall.StringToUTF16Ptr(path)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.SetIconLocation, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		uintptr(index))

	if hres != S_OK {
		return "[IShellLink.SetIconLocation] fail " + strconv.Itoa(int(hres))
	}

	return ""
}
func (obj *IShellLink) SetWorkingDirectory(path string) (errMsg string) {
	p := syscall.StringToUTF16Ptr(path)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.SetWorkingDirectory, 2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		0)

	if hres != S_OK {
		return "[IShellLink.SetWorkingDirectory] fail " + strconv.Itoa(int(hres))
	}

	return ""
}

func (obj *IShellLink) QueryInterface(riid REFIID, ppvObject *unsafe.Pointer) HRESULT {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.QueryInterface, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)))

	return HRESULT(ret)
}

func IShellLinkQueryInterfaceToIPersistFile(inObj *IShellLink) (outObj *IPersistFile, errMsg string) {
	var cpcPtr unsafe.Pointer
	hres := inObj.QueryInterface(&IID_IPersistFile, &cpcPtr)
	if hres < 0 {
		return nil, "[IShellLinkQueryInterfaceToIPersistFile] fail " + strconv.Itoa(int(hres))
	}
	outObj = (*IPersistFile)(cpcPtr)

	return outObj, ""
}

var CLSID_ShellLink = CLSID{0x00021401, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}

var IID_IShellLinkW = IID{0x000214F9, 0, 0, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}

func NewIShellLink() (psl *IShellLink, errMsg string) {
	var cfcObj unsafe.Pointer
	hres := CoCreateInstance(&CLSID_ShellLink, nil, CLSCTX_INPROC_SERVER, &IID_IShellLinkW, &cfcObj)
	if hres != S_OK {
		return nil, "[NewIShellLink] fail " + strconv.Itoa(int(hres))
	}
	psl = (*IShellLink)(cfcObj)

	return psl, ""
}

type IPersistFileVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetClassID uintptr

	IsDirty       uintptr
	Load          uintptr
	Save          uintptr
	SaveCompleted uintptr
	GetCurFile    uintptr
}
type IPersistFile struct {
	LpVtbl *IPersistFileVtbl
}

func (obj *IPersistFile) Release() uint32 {
	ret, _, _ := syscall.Syscall(obj.LpVtbl.Release, 1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)

	return uint32(ret)
}
func (obj *IPersistFile) Save(pszFileName string, fRemember bool) (errMsg string) {
	p := syscall.StringToUTF16Ptr(pszFileName)
	hres, _, _ := syscall.Syscall(obj.LpVtbl.Save, 3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(p)),
		uintptr(BoolToBOOL(fRemember)),
	)

	if hres != S_OK {
		return "[NewIShellLink] fail " + strconv.Itoa(int(hres))
	}

	return ""
}

var IID_IPersistFile = IID{0x0000010b, 0, 0, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
