package udwCache

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwMath"
	"github.com/tachyon-protocol/udw/udwSingleFlight"
	"sync"
	"time"
)

var errMsgCacheMiss = "cache miss"

type ttlCacheEntry struct {
	Value   interface{}
	Timeout time.Time
}

func (entry ttlCacheEntry) GetTtl() uint32 {
	ttlDur := entry.Timeout.Sub(time.Now())
	if ttlDur < 0 {
		ttlDur = 0
	}
	return uint32(udwMath.CeilToInt(ttlDur.Seconds()))
}

type TtlCache struct {
	cache       map[string]ttlCacheEntry
	lock        sync.RWMutex
	singleGroup udwSingleFlight.Group
}

func NewTtlCache() *TtlCache {
	return &TtlCache{
		cache: map[string]ttlCacheEntry{},
	}
}

func (s *TtlCache) DoWithTtl(key string, f func() (value interface{}, ttl uint32, err error)) (value interface{}, ttl uint32, err error) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, entry.GetTtl(), nil
	}
	if err.Error() != errMsgCacheMiss {
		return
	}
	entryi, err := s.singleGroup.Do(key, func() (interface{}, error) {
		value, ttl, err := f()
		timeout := time.Now().Add(time.Duration(ttl) * time.Second)
		return ttlCacheEntry{
			Value:   value,
			Timeout: timeout,
		}, err
	})
	entry = entryi.(ttlCacheEntry)
	ttl = entry.GetTtl()
	if err == nil && ttl > 0 {
		s.save(key, entry)
	}
	return entry.Value, ttl, nil
}

func (s *TtlCache) save(key string, entry ttlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *TtlCache) get(key string) (entry ttlCacheEntry, err error) {
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

func (s *TtlCache) GcThread() {
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
