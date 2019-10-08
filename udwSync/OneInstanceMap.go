package udwSync

import "sync"

type OneInstanceMap struct {
	m      map[string]struct{}
	locker sync.Mutex
}

func (m *OneInstanceMap) IsRunning(key string) bool {
	m.locker.Lock()
	if m.m == nil {
		m.m = map[string]struct{}{}
	}
	_, isRunning := m.m[key]
	m.locker.Unlock()
	return isRunning
}

func (m *OneInstanceMap) Run(key string, fn func()) bool {
	m.locker.Lock()
	if m.m == nil {
		m.m = map[string]struct{}{}
	}
	_, isRunning := m.m[key]
	if isRunning {
		m.locker.Unlock()
		return false
	}
	m.m[key] = struct{}{}
	m.locker.Unlock()
	defer func() {
		m.locker.Lock()
		delete(m.m, key)
		m.locker.Unlock()
	}()
	fn()
	return true
}
