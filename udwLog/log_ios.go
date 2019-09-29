// +build ios macAppStore

package udwLog

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
void cGoNSLog(const char *buf){
	@autoreleasepool {
		NSLog(@"%@",[[NSString alloc] initWithCString:buf encoding:NSUTF8StringEncoding]);
	}
}
*/
import "C"
import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"unsafe"
)

func log(bufW *udwBytes.BufWriter) {
	bufW.WriteByte_(0)
	blist := bufW.GetBytes()
	C.cGoNSLog((*C.char)(unsafe.Pointer(&blist[0])))
}
