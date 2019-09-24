package udwClose

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Closer struct {
	closeFuncLocker sync.Mutex
	fieldLocker     sync.Mutex
	closeChan       chan struct{}
	timeoutTimer    *time.Timer
	startTime       time.Time
	isClosed        bool
	closeFnList     []func()
}

func NewCloser() *Closer {
	return &Closer{
		startTime: time.Now(),
	}
}

func (c *Closer) IsClose() bool {

	c.fieldLocker.Lock()
	isClose := c.isClosed
	c.fieldLocker.Unlock()
	return isClose
}

func (c *Closer) CloseWithCallback(addCloser func()) {
	c.closeFuncLocker.Lock()
	if c.IsClose() {
		c.closeFuncLocker.Unlock()
		return
	}
	c.fieldLocker.Lock()
	c.isClosed = true
	if c.closeChan != nil {
		close(c.closeChan)
	}
	closerFnList := c.closeFnList
	c.closeFnList = nil
	if addCloser != nil {
		closerFnList = append(closerFnList, addCloser)
	}

	c.fieldLocker.Unlock()

	c.closeFuncLocker.Unlock()
	for i := len(closerFnList) - 1; i >= 0; i-- {
		closerFnList[i]()
	}
}

func (c *Closer) Close() {
	c.CloseWithCallback(nil)
}

func (c *Closer) GetCloseChan() <-chan struct{} {
	c.fieldLocker.Lock()
	if c.closeChan == nil {
		c.closeChan = make(chan struct{})

		if c.isClosed {
			close(c.closeChan)
		}
	}
	thisChan := c.closeChan
	c.fieldLocker.Unlock()
	return thisChan
}

func (c *Closer) AddOnClose(fn func()) {
	if c == nil {
		return
	}
	c.fieldLocker.Lock()
	if c.isClosed {
		c.fieldLocker.Unlock()
		fn()
		return
	}
	c.closeFnList = append(c.closeFnList, fn)
	c.fieldLocker.Unlock()
}

func (c *Closer) AddIoCloserOnCloseWithErrorLog(closer io.Closer) {
	c.AddOnClose(func() {
		err := closer.Close()
		if err != nil {
			fmt.Println("2chs7bstn8 " + err.Error())
		}
	})
}

func (c *Closer) LoopUntilCloseFirstRun(dur time.Duration, f func()) {
	timer := time.NewTimer(dur)
	closeChan := c.GetCloseChan()
	for {
		f()
		timer.Reset(dur)
		select {
		case <-timer.C:
		case <-closeChan:
			return
		}
	}
}

func (c *Closer) LoopUntilCloseFirstSleep(dur time.Duration, f func()) {
	timer := time.NewTimer(dur)
	closeChan := c.GetCloseChan()
	for {
		timer.Reset(dur)
		select {
		case <-timer.C:
		case <-closeChan:
			return
		}
		f()
	}
}

func (c *Closer) WaitClose() {
	closeChan := c.GetCloseChan()
	select {
	case <-closeChan:
		return
	}
}

func (c *Closer) SleepDur(dur time.Duration) bool {
	timer := time.NewTimer(dur)
	closeChan := c.GetCloseChan()
	select {
	case <-timer.C:
		return true
	case <-closeChan:
		return false
	}
}

func (c *Closer) SleepUtil(t time.Time) {
	startTime := time.Now()
	dur := t.Sub(startTime)
	if dur <= 0 {
		return
	}
	timer := time.NewTimer(dur)
	closeChan := c.GetCloseChan()
	select {
	case <-timer.C:
	case <-closeChan:
		return
	}
}

func (c *Closer) AddUpperCloser(upperCloser *Closer) {
	if upperCloser == nil {

		return
	}
	if c.IsClose() {

		return
	}

	thisC := c
	thisCLocker := sync.Mutex{}
	upperCloser.AddOnClose(func() {
		thisCLocker.Lock()
		if thisC != nil {
			thisC2 := thisC
			thisC = nil
			thisCLocker.Unlock()
			thisC2.Close()
			return
		}
		thisCLocker.Unlock()
	})
	c.AddOnClose(func() {
		thisCLocker.Lock()
		thisC = nil
		thisCLocker.Unlock()
	})
}
