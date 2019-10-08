package udwSync

import (
	"sync"
	"time"
)

type Time struct {
	locker sync.RWMutex
	v      time.Time
}

func (s *Time) Get() time.Time {
	s.locker.RLock()
	out := s.v
	s.locker.RUnlock()
	return out
}

func (s *Time) Set(v time.Time) {
	s.locker.Lock()
	s.v = v
	s.locker.Unlock()
}

func NewTime(s time.Time) *Time {
	return &Time{
		v: s,
	}
}
