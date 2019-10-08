package udwSync

import (
	"sync"
	"time"
)

type BoolWaiter struct {
	locker           sync.Mutex
	status           bool
	statusChangeChan chan struct{}
}

func (bw *BoolWaiter) Set(b bool) {
	bw.locker.Lock()
	bw.initValue__NOLOCK()
	if bw.status == b {

		bw.locker.Unlock()
		return
	}
	bw.status = b
	close(bw.statusChangeChan)
	bw.statusChangeChan = make(chan struct{})
	bw.locker.Unlock()
}
func (bw *BoolWaiter) Get() bool {
	bw.locker.Lock()
	b := bw.status
	bw.locker.Unlock()
	return b
}

func (bw *BoolWaiter) initValue__NOLOCK() {
	if bw.statusChangeChan == nil {
		bw.statusChangeChan = make(chan struct{})
	}
}

func (bw *BoolWaiter) WaitTrueWithTimeout(dur time.Duration) bool {
	return bw.WaitWithTimeout(true, dur)
}

func (bw *BoolWaiter) WaitWithTimeout(needStatus bool, dur time.Duration) bool {
	waitChan := time.After(dur)
	for {
		bw.locker.Lock()
		bw.initValue__NOLOCK()
		b := bw.status
		thisChan := bw.statusChangeChan
		bw.locker.Unlock()
		if b == needStatus {
			return true
		}
		select {
		case <-thisChan:
			return true
		case <-waitChan:
			return false
		}
	}
}

func (bw *BoolWaiter) Wait(needStatus bool) {
	for {
		bw.locker.Lock()
		bw.initValue__NOLOCK()
		b := bw.status
		thisChan := bw.statusChangeChan
		bw.locker.Unlock()
		if b == needStatus {
			return
		}
		select {
		case <-thisChan:
			return
		}
	}
}
