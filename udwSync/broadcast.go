package udwSync

import (
	"sync"
)

type Broadcast struct {
	cond sync.Cond
	cv   int
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		cond: sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (b *Broadcast) Broadcast() {
	b.cond.L.Lock()
	b.cv++
	b.cond.L.Unlock()
	b.cond.Broadcast()
}

func (b *Broadcast) WaitWithCb(fn func() bool) {
	thisV := 0
	b.cond.L.Lock()
	for {
		if b.cv == thisV {
			b.cond.Wait()
		}
		thisV = b.cv
		b.cond.L.Unlock()
		ret := fn()
		if ret {
			return
		}
		b.cond.L.Lock()
	}
}

func (b *Broadcast) WaitWithVersion(v int) int {
	b.cond.L.Lock()
	if b.cv == v {
		b.cond.Wait()
	}
	thisV := b.cv
	b.cond.L.Unlock()
	return thisV
}
