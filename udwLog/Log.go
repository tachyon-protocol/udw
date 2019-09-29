package udwLog

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTime"
	"sync"
	"time"
)

func Log(a ...interface{}) {
	bufW := gBufWriterPool.Get()
	outB := time.Now().In(udwTime.GetUtc8Zone()).AppendFormat(bufW.GetBytes(), udwTime.FormatUdwLog)
	bufW.Write_(outB)
	bufW.WriteByte_(' ')
	_, _ = fmt.Fprintln(bufW, a...)
	writeToUdwTee(bufW.GetBytes())
	bufW.AddPos(-1)
	log(bufW)
	if bufW.GetLen() > 1024*4 {

		bufW.ResetWithBuffer(nil)
	}
	gBufWriterPool.Put(bufW)
}

func LogAsync(a ...interface{}) {
	go func() {
		gLogAsyncLocker.Lock()
		Log(a...)
		gLogAsyncLocker.Unlock()
	}()
}

var gLogAsyncLocker sync.Mutex

var gBufWriterPool udwBytes.BufWriterPool

type LogPrintlnFunc func(a ...interface{})

func (f LogPrintlnFunc) LogLn(a ...interface{}) {
	if f == nil {
		Log(a...)
		return
	}
	f(a...)
}

type LogRow struct {
	Cat  string
	Time time.Time

	Data []interface{}
}

var gUdwTeeFn func(buf []byte)
var gUdwTeeFnLocker sync.RWMutex

func SetUdwTeeFn(fn func(buf []byte)) {
	gUdwTeeFnLocker.Lock()
	gUdwTeeFn = fn
	gUdwTeeFnLocker.Unlock()
}

func writeToUdwTee(buf []byte) {
	gUdwTeeFnLocker.RLock()
	if gUdwTeeFn != nil {
		gUdwTeeFn(buf)
	}
	gUdwTeeFnLocker.RUnlock()
}
