package udwSync

import (
	"fmt"
	"sync"
)

type OneInstance struct {
	locker sync.Mutex
	isRun  bool
	Name   string
}

func (oi *OneInstance) RunAndLogIfNotRun(fn func()) {
	ret := oi.Run(fn)
	if ret == false {
		fmt.Println("[udwSync.OneInstance] too many instance2 " + oi.Name)
	}
}

func (oi *OneInstance) IsRun() (isRun bool) {
	oi.locker.Lock()
	isRun = oi.isRun
	oi.locker.Unlock()
	return isRun
}

func (oi *OneInstance) Run(fn func()) bool {
	oi.locker.Lock()
	if oi.isRun {
		oi.locker.Unlock()
		return false
	}
	oi.isRun = true
	oi.locker.Unlock()
	defer func() {
		oi.locker.Lock()
		oi.isRun = false
		oi.locker.Unlock()
	}()
	fn()
	return true
}

func (oi *OneInstance) MustRun(fn func()) {
	oi.locker.Lock()
	if oi.isRun {
		oi.locker.Unlock()
		panic("[udwSync.OneInstance] too many instance1 " + oi.Name)

	}
	oi.isRun = true
	oi.locker.Unlock()
	defer func() {
		oi.locker.Lock()
		oi.isRun = false
		oi.locker.Unlock()
	}()
	fn()
	return
}

func (oi *OneInstance) MustStart() {
	oi.locker.Lock()
	if oi.isRun {
		oi.locker.Unlock()
		panic(`axwwu4p8tz too many instance1`)
	}
	oi.isRun = true
	oi.locker.Unlock()
}

func (oi *OneInstance) MustStop() {
	oi.locker.Lock()
	if !oi.isRun {
		oi.locker.Unlock()
		panic(`tmszrajzfx `)
	}
	oi.isRun = false
	oi.locker.Unlock()
}
