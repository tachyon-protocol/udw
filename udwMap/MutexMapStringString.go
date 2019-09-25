package udwMap

import "sync"

type MutexMapStringString struct {
	lock sync.RWMutex
	m    map[string]string
}

func (m *MutexMapStringString) Set(k string, v string) {
	m.lock.Lock()
	if m.m == nil {
		m.m = map[string]string{}
	}
	m.m[k] = v
	m.lock.Unlock()
}

func (m *MutexMapStringString) Get(k string) string {
	m.lock.RLock()
	if m.m == nil {
		m.lock.RUnlock()
		return ""
	}
	v := m.m[k]
	m.lock.RUnlock()
	return v
}

func (m *MutexMapStringString) Del(k string) {
	m.lock.Lock()
	if m.m == nil {
		m.lock.Unlock()
		return
	}
	delete(m.m, k)
	m.lock.Unlock()
}
