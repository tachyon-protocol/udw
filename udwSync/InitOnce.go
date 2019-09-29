package udwSync

import (
	"github.com/tachyon-protocol/udw/udwSingleFlight"
	"sync"
)

type InitOnceWithError struct {
	single     udwSingleFlight.One
	hasRunOnce bool
	lastError  string
	locker     sync.Mutex
}

func (once *InitOnceWithError) Do(fn func() (errS string)) (errS string) {
	once.locker.Lock()
	if once.hasRunOnce && once.lastError == "" {
		once.locker.Unlock()
		return ""
	}
	once.lastError = "NotInit"
	once.hasRunOnce = true
	once.locker.Unlock()
	once.single.DoNoReturn(func() {
		lastError := fn()
		once.locker.Lock()
		once.lastError = lastError
		once.locker.Unlock()
	})
	once.locker.Lock()
	lastErr := once.lastError
	once.locker.Unlock()
	return lastErr
}

func (once *InitOnceWithError) HasInitSucc() bool {
	once.locker.Lock()
	out := (once.hasRunOnce && once.lastError == "")
	once.locker.Unlock()
	return out
}
