// +build windows

package udwW32

import (
	"unsafe"
)

const (
	MIIM_STATE      = 1
	MIIM_ID         = 2
	MIIM_SUBMENU    = 4
	MIIM_CHECKMARKS = 8
	MIIM_TYPE       = 16
	MIIM_DATA       = 32
	MIIM_STRING     = 64
	MIIM_BITMAP     = 128
	MIIM_FTYPE      = 256
)

const (
	MFT_BITMAP       = 4
	MFT_MENUBARBREAK = 32
	MFT_MENUBREAK    = 64
	MFT_OWNERDRAW    = 256
	MFT_RADIOCHECK   = 512
	MFT_RIGHTJUSTIFY = 0x4000
	MFT_SEPARATOR    = 0x800
	MFT_RIGHTORDER   = 0x2000
	MFT_STRING       = 0
)

const (
	MFS_CHECKED   = 8
	MFS_DEFAULT   = 4096
	MFS_DISABLED  = 3
	MFS_ENABLED   = 0
	MFS_GRAYED    = 3
	MFS_HILITE    = 128
	MFS_UNCHECKED = 0
	MFS_UNHILITE  = 0
)

const (
	HBMMENU_CALLBACK        = -1
	HBMMENU_SYSTEM          = 1
	HBMMENU_MBAR_RESTORE    = 2
	HBMMENU_MBAR_MINIMIZE   = 3
	HBMMENU_MBAR_CLOSE      = 5
	HBMMENU_MBAR_CLOSE_D    = 6
	HBMMENU_MBAR_MINIMIZE_D = 7
	HBMMENU_POPUP_CLOSE     = 8
	HBMMENU_POPUP_RESTORE   = 9
	HBMMENU_POPUP_MAXIMIZE  = 10
	HBMMENU_POPUP_MINIMIZE  = 11
)

const (
	MIM_APPLYTOSUBMENUS = 0x80000000
	MIM_BACKGROUND      = 0x00000002
	MIM_HELPID          = 0x00000004
	MIM_MAXHEIGHT       = 0x00000001
	MIM_MENUDATA        = 0x00000008
	MIM_STYLE           = 0x00000010
)

const (
	MNS_AUTODISMISS = 0x10000000
	MNS_CHECKORBMP  = 0x04000000
	MNS_DRAGDROP    = 0x20000000
	MNS_MODELESS    = 0x40000000
	MNS_NOCHECK     = 0x80000000
	MNS_NOTIFYBYPOS = 0x08000000
)

const (
	MF_BYCOMMAND  = 0x00000000
	MF_BYPOSITION = 0x00000400
)

type MENUITEMINFO struct {
	CbSize        uint32
	FMask         uint32
	FType         uint32
	FState        uint32
	WID           uint32
	HSubMenu      HMENU
	HbmpChecked   HBITMAP
	HbmpUnchecked HBITMAP
	DwItemData    uintptr
	DwTypeData    *uint16
	Cch           uint32
	HbmpItem      HBITMAP
}

type MENUINFO struct {
	CbSize          uint32
	FMask           uint32
	DwStyle         uint32
	CyMax           uint32
	HbrBack         HBRUSH
	DwContextHelpID uint32
	DwMenuData      uintptr
}

var procCreatePopupMenu = DllProc{Dll: &moduser32, Name: "CreatePopupMenu"}

func CreatePopupMenu() (hmenu HMENU, errMsg string) {
	ret, _, errMsg := procCreatePopupMenu.CallErrorMsg()

	return HMENU(ret), errMsg
}

func MustCreatePopupMenu() (hmenu HMENU) {
	hmenu, errMsg := CreatePopupMenu()
	if errMsg != "" {
		panic(errMsg)
	}
	return hmenu
}

var procGetMenuInfo = DllProc{Dll: &moduser32, Name: "GetMenuInfo"}

func GetMenuInfo(hmenu HMENU, lpcmi *MENUINFO) (errMsg string) {
	ret, _, errMsg := procGetMenuInfo.CallErrorMsg(
		uintptr(hmenu),
		uintptr(unsafe.Pointer(lpcmi)),
	)
	if ret != 0 {
		return ""
	}
	return errMsg
}

func MustGetMenuInfo(hmenu HMENU, lpcmi *MENUINFO) {
	errMsg := GetMenuInfo(hmenu, lpcmi)
	if errMsg != "" {
		panic(errMsg)
	}
}

var procSetMenuInfo = DllProc{Dll: &moduser32, Name: "SetMenuInfo"}

func SetMenuInfo(hmenu HMENU, lpcmi *MENUINFO) (errMsg string) {
	ret, _, errMsg := procSetMenuInfo.CallErrorMsg(
		uintptr(hmenu),
		uintptr(unsafe.Pointer(lpcmi)),
	)
	if ret != 0 {
		return ""
	}
	return errMsg
}

func MustSetMenuInfo(hmenu HMENU, lpcmi *MENUINFO) {
	errMsg := SetMenuInfo(hmenu, lpcmi)
	if errMsg != "" {
		panic(errMsg)
	}
}

var procInsertMenuItem = DllProc{Dll: &moduser32, Name: "InsertMenuItemW"}

func InsertMenuItem(hMenu HMENU, uItem uint32, fByPosition bool, lpmii *MENUITEMINFO) (errMsg string) {
	ret, _, errMsg := procInsertMenuItem.CallErrorMsg(
		uintptr(hMenu),
		uintptr(uItem),
		uintptr(BoolToBOOL(fByPosition)),
		uintptr(unsafe.Pointer(lpmii)),
	)
	if ret != 0 {
		return ""
	}
	return errMsg
}

func MustInsertMenuItem(hMenu HMENU, uItem uint32, fByPosition bool, lpmii *MENUITEMINFO) {
	errMsg := InsertMenuItem(hMenu, uItem, fByPosition, lpmii)
	if errMsg != "" {
		panic(errMsg)
	}
}

const (
	TPM_CENTERALIGN     = 0x0004
	TPM_LEFTALIGN       = 0x0000
	TPM_RIGHTALIGN      = 0x0008
	TPM_BOTTOMALIGN     = 0x0020
	TPM_TOPALIGN        = 0x0000
	TPM_VCENTERALIGN    = 0x0010
	TPM_NONOTIFY        = 0x0080
	TPM_RETURNCMD       = 0x0100
	TPM_LEFTBUTTON      = 0x0000
	TPM_RIGHTBUTTON     = 0x0002
	TPM_HORNEGANIMATION = 0x0800
	TPM_HORPOSANIMATION = 0x0400
	TPM_NOANIMATION     = 0x4000
	TPM_VERNEGANIMATION = 0x2000
	TPM_VERPOSANIMATION = 0x1000
	TPM_HORIZONTAL      = 0x0000
	TPM_VERTICAL        = 0x0040
)

type TPMPARAMS struct {
	CbSize    uint32
	RcExclude RECT
}

var procTrackPopupMenuEx = DllProc{Dll: &moduser32, Name: "TrackPopupMenuEx"}

func TrackPopupMenuEx(hMenu HMENU, fuFlags uint32, x, y int32, hWnd HWND, lptpm *TPMPARAMS) (actionId uint32, errMsg string) {
	ret, _, errMsg := procTrackPopupMenuEx.CallErrorMsg(uintptr(hMenu),
		uintptr(fuFlags),
		uintptr(x),
		uintptr(y),
		uintptr(hWnd),
		uintptr(unsafe.Pointer(lptpm)),
	)
	return uint32(ret), errMsg
}

func MustTrackPopupMenuEx(hMenu HMENU, fuFlags uint32, x, y int32, hWnd HWND, lptpm *TPMPARAMS) (actionId uint32) {
	actionId, errMsg := TrackPopupMenuEx(hMenu, fuFlags, x, y, hWnd, lptpm)
	if errMsg != "" {
		panic(errMsg)
	}
	return actionId
}

var procDestroyMenu = DllProc{Dll: &moduser32, Name: "DestroyMenu"}

func DestroyMenu(hMenu HMENU) (errMsg string) {
	_, _, errMsg = procDestroyMenu.CallErrorMsg(uintptr(hMenu))
	return errMsg
}
