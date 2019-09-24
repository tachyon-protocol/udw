package udwMap

import "sync"

type MutexSetString struct {
	locker sync.RWMutex
	m      map[string]struct{}
}

func (ss *MutexSetString) Has(k string) bool {
	ss.locker.RLock()
	if ss.m == nil {
		ss.locker.RUnlock()
		return false
	}
	_, ok := ss.m[k]
	ss.locker.RUnlock()
	return ok
}

func (ss *MutexSetString) Set(k string) {
	ss.locker.Lock()
	if ss.m == nil {
		ss.m = map[string]struct{}{}
	}
	ss.m[k] = struct{}{}
	ss.locker.Unlock()
}

func (ss *MutexSetString) SetAllByMap(m map[string]struct{}) {
	ss.locker.Lock()
	ss.m = m
	ss.locker.Unlock()
}

func (ss *MutexSetString) Len() int {
	ss.locker.RLock()
	sz := len(ss.m)
	ss.locker.RUnlock()
	return sz
}

func (ss *MutexSetString) Clear() {
	ss.locker.Lock()
	ss.m = nil
	ss.locker.Unlock()
}
