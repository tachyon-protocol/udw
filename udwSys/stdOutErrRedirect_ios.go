// +build ios macAppStore

package udwSys

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
const char *cGoGetCacheDir(){
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSCachesDirectory, NSUserDomainMask, YES);
    NSString *cachePath = [paths objectAtIndex:0];
    return [cachePath UTF8String];
}
*/
import "C"

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwSys/udwSysLog"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"
	"unsafe"
)

func StdOutErrRedirectToNSLog() {
	udwSysLog.StdOutRedirectIos()
	return
}

func StdOutErrRedirectToNSLogV2() {
	dir := C.GoString(C.cGoGetCacheDir())
	NSLog(fmt.Sprintln("wd", udwFile.MustGetWd(), dir))
	outLogPath := filepath.Join(dir, "out.log")
	errLogPath := filepath.Join(dir, "err.log")
	os.Stdout.Close()
	os.Stderr.Close()
	stdErrRediectToNSLog(outLogPath, "StdOut")
	stdErrRediectToNSLog(errLogPath, "StdErr")
}
func stdErrRediectToNSLog(logPath string, typeName string) {
	if udwFile.MustFileExist(logPath) {
		NSLog(fmt.Sprintln("last", typeName, "============================\n"))
		readfileToNsLog(logPath)
		NSLog(fmt.Sprintln("\nend of ", logPath, "===================================="))
	}
	f, err := os.Create(logPath)
	if err != nil {
		NSLog(fmt.Sprintln(typeName, "log file open fail", err))
		return
	}
	if typeName == "StdOut" {
		os.Stdout = f
	} else if typeName == "StdErr" {
		os.Stderr = f
	}
	now := time.Now()
	fmt.Fprintln(f, typeName, "logStart", now)
}

func readfileToNsLog(logPath string) {
	const tmpStackBufSize = 4096
	var tmpStackBuf [tmpStackBufSize]byte
	var stackBuf []byte
	tmpStackBufp := uintptr(unsafe.Pointer(&tmpStackBuf))
	tmpStackBufp2 := *(*uintptr)(unsafe.Pointer(&tmpStackBufp))
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&stackBuf))
	bx.Data = uintptr(tmpStackBufp2)
	bx.Len = tmpStackBufSize
	bx.Cap = tmpStackBufSize

	f, err := os.Open(logPath)
	if err != nil {
		NSLog(fmt.Sprintln("os.Open fail", err))
		return
	}
	defer f.Close()
	for {
		n, err := f.Read(stackBuf[:tmpStackBufSize-1])
		if err == io.EOF {
			return
		}
		if err != nil {
			NSLog(fmt.Sprintln("f.Read fail", err))
			return
		}
		stackBuf[n] = 0
		nSLogWithByteSliceAndCEnding(tmpStackBufp2)

	}
}

func NSLog(s string) {
	udwLog.Log(s)
}

func nSLogWithByteSliceAndCEnding(p uintptr) {
	udwSysLog.NSLogWithByteSliceAndCEnding(p)
}
