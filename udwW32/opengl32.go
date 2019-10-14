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
	modopengl32 = Dll{Name: "opengl32.dll"}

	procwglCreateContext      = DllProc{Dll: &modopengl32, Name: "wglCreateContext"}
	procwglCreateLayerContext = DllProc{Dll: &modopengl32, Name: "wglCreateLayerContext"}
	procwglDeleteContext      = DllProc{Dll: &modopengl32, Name: "wglDeleteContext"}
	procwglGetProcAddress     = DllProc{Dll: &modopengl32, Name: "wglGetProcAddress"}
	procwglMakeCurrent        = DllProc{Dll: &modopengl32, Name: "wglMakeCurrent"}
	procwglShareLists         = DllProc{Dll: &modopengl32, Name: "wglShareLists"}
)

func WglCreateContext(hdc HDC) HGLRC {
	ret, _, _ := procwglCreateContext.Call(
		uintptr(hdc),
	)

	return HGLRC(ret)
}

func WglCreateLayerContext(hdc HDC, iLayerPlane int) HGLRC {
	ret, _, _ := procwglCreateLayerContext.Call(
		uintptr(hdc),
		uintptr(iLayerPlane),
	)

	return HGLRC(ret)
}

func WglDeleteContext(hglrc HGLRC) bool {
	ret, _, _ := procwglDeleteContext.Call(
		uintptr(hglrc),
	)

	return ret == TRUE
}

func WglGetProcAddress(szProc string) uintptr {
	ret, _, _ := procwglGetProcAddress.Call(
		uintptr(unsafe.Pointer(syscall.StringBytePtr(szProc))),
	)

	return ret
}

func WglMakeCurrent(hdc HDC, hglrc HGLRC) bool {
	ret, _, _ := procwglMakeCurrent.Call(
		uintptr(hdc),
		uintptr(hglrc),
	)

	return ret == TRUE
}

func WglShareLists(hglrc1, hglrc2 HGLRC) bool {
	ret, _, _ := procwglShareLists.Call(
		uintptr(hglrc1),
		uintptr(hglrc2),
	)

	return ret == TRUE
}
