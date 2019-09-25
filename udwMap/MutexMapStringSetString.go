package udwMap

import "sync"

type MutexMapStringSetString struct {
	locker sync.RWMutex
	m      map[string]map[string]struct{}
}

func (m *MutexMapStringSetString) Add(k1 string, k2 string) {
	m.locker.Lock()
	if m.m == nil {
		m.m = map[string]map[string]struct{}{}
	}
	k2M := m.m[k1]
	if k2M == nil {
		k2M = map[string]struct{}{}
		m.m[k1] = k2M
	}
	k2M[k2] = struct{}{}
	m.locker.Unlock()

}
func (m *MutexMapStringSetString) Has(k1 string, k2 string) bool {
	m.locker.RLock()
	if m.m == nil {
		m.locker.RUnlock()
		return false
	}
	k2M := m.m[k1]
	if k2M == nil {
		m.locker.RUnlock()
		return false
	}
	_, ok := k2M[k2]
	m.locker.RUnlock()
	return ok
}
func (m *MutexMapStringSetString) GetK1Len(k1 string) int {
	m.locker.RLock()
	if m.m == nil {
		m.locker.RUnlock()
		return 0
	}
	k2M := m.m[k1]
	if k2M == nil {
		m.locker.RUnlock()
		return 0
	}
	l := len(k2M)
	m.locker.RUnlock()
	return l
}
