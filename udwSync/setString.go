package udwSync

import "sync"

type SetString struct {
	locker sync.RWMutex
	m      map[string]struct{}
}

func (ss *SetString) IsIn(k string) bool {
	ss.locker.RLock()
	if ss.m == nil {
		ss.locker.RUnlock()
		return false
	}
	_, ok := ss.m[k]
	ss.locker.RUnlock()
	return ok
}

func (ss *SetString) Has(k string) bool {
	return ss.IsIn(k)
}

func (ss *SetString) Set(k string) {
	ss.locker.Lock()
	if ss.m == nil {
		ss.m = map[string]struct{}{}
	}
	ss.m[k] = struct{}{}
	ss.locker.Unlock()
}

func (ss *SetString) Add(k string) {
	ss.Set(k)
}

func (ss *SetString) Len() int {
	ss.locker.RLock()
	sz := len(ss.m)
	ss.locker.RUnlock()
	return sz
}

func (ss *SetString) Clear() {
	ss.locker.Lock()
	ss.m = nil
	ss.locker.Unlock()
}

func (ss *SetString) GetStringSliceAndClear() (list []string) {
	ss.locker.Lock()
	data := ss.m
	ss.m = nil
	ss.locker.Unlock()

	if len(data) > 0 {
		list = make([]string, 0, len(data))
		for key := range data {
			list = append(list, key)
		}
	}
	return
}
