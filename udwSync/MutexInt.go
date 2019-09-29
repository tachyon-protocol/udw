package udwSync

import "sync"

type Int struct {
	locker sync.RWMutex
	v      int
}

func (s *Int) Get() int {
	s.locker.RLock()
	out := s.v
	s.locker.RUnlock()
	return out
}

func (s *Int) Set(v int) {
	s.locker.Lock()
	s.v = v
	s.locker.Unlock()
}

func (s *Int) Add(toAdd int) {
	s.locker.Lock()
	s.v += toAdd
	s.locker.Unlock()
}

func (s *Int) AddAndReturnNew(toAdd int) int {
	s.locker.Lock()
	v := s.v + toAdd
	s.v = v
	s.locker.Unlock()
	return v
}

func NewInt(s int) *Int {
	return &Int{
		v: s,
	}
}
