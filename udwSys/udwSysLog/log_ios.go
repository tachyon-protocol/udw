// +build ios macAppStore

package udwSysLog

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
void cGoNSLog(const char *buf){
	// 加入 autoreleasepool 减少内存占用.
	@autoreleasepool {
		NSLog(@"%@",[[NSString alloc] initWithCString:buf encoding:NSUTF8StringEncoding]);
	}
}
*/
import "C"
import (
	"unsafe"
)

func log(s string) {
	cstr := C.CString(s)
	C.cGoNSLog(cstr)
	C.free(unsafe.Pointer(cstr))
}

func NSLogWithByteSliceAndCEnding(p uintptr) {
	C.cGoNSLog((*C.char)(unsafe.Pointer(p)))
}
