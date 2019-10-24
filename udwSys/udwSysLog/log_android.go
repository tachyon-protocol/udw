// +build android

package udwSysLog

/*
#cgo LDFLAGS: -landroid -llog
#include <android/log.h>
#include <string.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

var (
	ctag = C.CString("GoLog")
)

func log(s string) {
	cstr := C.CString(s)
	C.__android_log_write(C.ANDROID_LOG_INFO, ctag, cstr)
	C.free(unsafe.Pointer(cstr))
}
