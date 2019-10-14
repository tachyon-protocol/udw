package udwCache

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwSingleFlight"
	"sync"
	"time"
)

type MemoryTtlCacheV2 struct {
	cache       map[string]ttlCacheEntry
	lock        sync.RWMutex
	singleGroup udwSingleFlight.Group
}

func NewMemoryTtlCacheV2() *MemoryTtlCacheV2 {
	c := &MemoryTtlCacheV2{
		cache: map[string]ttlCacheEntry{},
	}
	go c.GcThread()
	return c
}

func (s *MemoryTtlCacheV2) Do(key string, f func() (value interface{}, ttl time.Duration, err error)) (value interface{}, err error) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, nil
	}
	if err.Error() != errMsgCacheMiss {
		return
	}
	entryi, err := s.singleGroup.Do(key, func() (interface{}, error) {
		value, ttl, err := f()
		return ttlCacheEntry{
			Value:   value,
			Timeout: time.Now().Add(ttl),
		}, err
	})
	if err != nil {
		return nil, err
	}
	entry = entryi.(ttlCacheEntry)
	s.save(key, entry)
	return entry.Value, nil
}

func (s *MemoryTtlCacheV2) Remove(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.cache, key)
}

func (s *MemoryTtlCacheV2) save(key string, entry ttlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *MemoryTtlCacheV2) get(key string) (entry ttlCacheEntry, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now()
	entry, ok := s.cache[key]
	if !ok {
		return entry, errors.New(errMsgCacheMiss)
	}
	if now.After(entry.Timeout) {
		return entry, errors.New(errMsgCacheMiss)
	}
	return entry, nil
}

func (s *MemoryTtlCacheV2) GcThread() {
	for {
		time.Sleep(time.Hour)
		s.lock.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			if now.After(entry.Timeout) {
				delete(s.cache, key)
			}
		}
		s.lock.Unlock()
	}
}
