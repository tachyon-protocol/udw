package udwSync

import (
	"sync"
)

type String struct {
	locker sync.RWMutex
	v      string
}

func (s *String) Get() string {
	s.locker.RLock()
	out := s.v
	s.locker.RUnlock()
	return out
}

func (s *String) Set(v string) {
	s.locker.Lock()
	s.v = v
	s.locker.Unlock()
}

func NewString(s string) *String {
	return &String{
		v: s,
	}
}
