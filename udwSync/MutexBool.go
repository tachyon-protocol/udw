package udwSync

import (
	"sync"
)

type Bool struct {
	locker sync.RWMutex
	v      bool
}

func (s *Bool) Get() bool {
	s.locker.RLock()
	out := s.v
	s.locker.RUnlock()
	return out
}

func (s *Bool) Set(v bool) {
	s.locker.Lock()
	s.v = v
	s.locker.Unlock()
}

func NewBool(s bool) *Bool {
	return &Bool{
		v: s,
	}
}
