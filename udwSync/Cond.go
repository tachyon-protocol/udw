package udwSync

import (
	"sync"
)

type Cond struct {
	c             sync.Cond
	l             sync.Mutex
	initOnce      sync.Once
	isClose       bool
	isCloseLocker sync.Mutex
}

func (c *Cond) WaitCheck(isTrue func() bool) {
	c.initValue()
	c.l.Lock()
	for {
		if isTrue() {
			c.l.Unlock()
			return
		}

		if c.isClose {
			c.l.Unlock()
			return
		}
		c.c.Wait()
	}
}

func (c *Cond) InLockAndSignal(f func()) {
	c.initValue()
	c.l.Lock()
	defer c.l.Unlock()
	f()
	c.c.Signal()
}

func (c *Cond) WaitCheckAndDo(isTrue func() bool, doAfterTrue func()) {
	c.WaitCheck(func() bool {
		ret := isTrue()
		if ret {
			doAfterTrue()
		}
		return ret
	})
}

func (c *Cond) InLock(f func()) {
	c.initValue()
	c.l.Lock()
	defer c.l.Unlock()
	f()
}

func (c *Cond) Close() {
	c.initValue()
	c.l.Lock()
	c.isClose = true
	c.l.Unlock()
	c.Broadcast()
}

func (c *Cond) Signal() {
	c.initValue()
	c.c.Signal()
}

func (c *Cond) Broadcast() {
	c.initValue()
	c.c.Broadcast()
}

func (c *Cond) initValue() {
	c.initOnce.Do(func() {
		c.c.L = &c.l
	})
}
