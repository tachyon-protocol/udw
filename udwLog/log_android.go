// +build android

package udwLog

/*
#cgo LDFLAGS: -landroid -llog
#include <android/log.h>
#include <string.h>
*/
import "C"

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"unsafe"
)

var (
	ctag = C.CString("GoLog")
)

func log(bufW *udwBytes.BufWriter) {
	bufW.WriteByte_(0)
	blist := bufW.GetBytes()
	C.__android_log_write(C.ANDROID_LOG_INFO, ctag, (*C.char)(unsafe.Pointer(&blist[0])))
}
