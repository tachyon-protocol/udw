package udwTime

import (
	"sync"
	"time"
)

type TtlStringTimeMap struct {
	m           map[string]time.Time
	locker      sync.Mutex
	ttlDuration time.Duration
	gcNum       int
}

func TtlStringTimeMapNew(ttlDuration time.Duration) *TtlStringTimeMap {
	return &TtlStringTimeMap{
		m:           map[string]time.Time{},
		ttlDuration: ttlDuration,
	}
}

func (s *TtlStringTimeMap) Add(key string, t time.Time, now time.Time) (has bool) {
	s.locker.Lock()
	oldT, ok := s.m[key]
	if ok {
		deleteTime := now.Add(-s.ttlDuration)
		if oldT.Before(deleteTime) {
			delete(s.m, key)
		} else {
			s.locker.Unlock()
			return true
		}
	}
	s.m[key] = t
	if len(s.m) > 1000 {
		if s.gcNum%1000 == 0 {
			deleteTime := now.Add(-s.ttlDuration)
			for token, t := range s.m {
				if t.Before(deleteTime) {
					delete(s.m, token)
				}
			}
		}
		s.gcNum++
	}
	s.locker.Unlock()
	return false
}

func (s *TtlStringTimeMap) AddNow(key string, now time.Time) (has bool) {
	return s.Add(key, now, now)
}
