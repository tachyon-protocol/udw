package udwTask

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"sync"
)

func RunFunctionListConcurrent(fnList ...func()) {
	if len(fnList) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(fnList))
	for _, fn := range fnList {
		fn := fn
		go func() {
			fn()
			wg.Done()
		}()
	}
	wg.Wait()
}

func MustRunFunctionListConcurrentPassByPanic(fnList ...func()) {
	if len(fnList) == 0 {
		return
	}
	var panicMsgLocker sync.Mutex
	var panicMsg string
	wg := sync.WaitGroup{}
	wg.Add(len(fnList))
	for _, fn := range fnList {
		fn := fn
		go func() {
			thisErrorMsg := udwErr.PanicToErrorMsgWithStack(fn)
			if thisErrorMsg != "" {
				panicMsgLocker.Lock()
				panicMsg += thisErrorMsg + "\n"
				panicMsgLocker.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if panicMsg != "" {
		panic(panicMsg)
	}
}
