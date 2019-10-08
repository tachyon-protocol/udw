package udwSync

import "sync"

type MapStringString struct {
	locker sync.RWMutex
	m      map[string]string
}

func (m *MapStringString) Get(k string) string {
	m.locker.RLock()
	if m.m == nil {
		m.locker.RUnlock()
		return ""
	}
	v := m.m[k]
	m.locker.RUnlock()
	return v
}

func (m *MapStringString) Set(k string, v string) {
	m.locker.Lock()
	if m.m == nil {
		m.m = map[string]string{}
	}
	m.m[k] = v
	m.locker.Unlock()
}
