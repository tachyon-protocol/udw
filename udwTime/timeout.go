package udwTime

import (
	"time"
)

func Timeout(f func(), dur time.Duration) (hasTimeout bool) {
	finishChan := make(chan struct{})
	go func() {
		f()
		finishChan <- struct{}{}
	}()
	select {
	case <-finishChan:
		return false
	case <-time.After(dur):
		return true
	}
}

func MustNotTimeout(f func(), dur time.Duration) {
	finishChan := make(chan struct{})
	go func() {
		select {
		case <-finishChan:
			return
		case <-time.After(dur):
			panic("time out")
		}
	}()
	f()
	finishChan <- struct{}{}
	return
}
