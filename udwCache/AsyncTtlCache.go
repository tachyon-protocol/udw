package udwCache

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwSingleFlight"
	"sync"
	"time"
)

var errMsgCacheExpire = "cache expire"
var errMsgDoNotNeedCache = "do not need cache"

type AsyncTtlCache struct {
	cache       map[string]ttlCacheEntry
	lock        sync.RWMutex
	singleGroup udwSingleFlight.Group
}

func NewAsyncCache() *AsyncTtlCache {
	return &AsyncTtlCache{
		cache: map[string]ttlCacheEntry{},
	}
}

func (s *AsyncTtlCache) DoWithTtl(key string, f func() (value interface{}, ttl uint32, canSave bool)) (value interface{}, ttl uint32) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, entry.GetTtl()
	}
	updateCache := func() (value interface{}, ttl uint32) {

		entryi, err := s.singleGroup.Do(key, func() (out interface{}, err error) {
			value, ttl, canSave := f()
			timeout := time.Now().Add(time.Duration(ttl) * time.Second)
			out = ttlCacheEntry{
				Value:   value,
				Timeout: timeout,
			}
			if !canSave {
				err = errors.New(errMsgDoNotNeedCache)
			}
			return
		})
		entryn := entryi.(ttlCacheEntry)
		if err == nil {
			s.save(key, entryn)
		}
		ttl = entryn.GetTtl()
		return entryn.Value, ttl
	}
	if err == nil {
		return nil, 0
	}
	errS := err.Error()
	switch errS {
	case errMsgCacheMiss:
		value, ttl := updateCache()
		return value, ttl
	case errMsgCacheExpire:
		go updateCache()
		return entry.Value, 0
	default:
		return nil, 0
	}

}

func (s *AsyncTtlCache) save(key string, entry ttlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *AsyncTtlCache) get(key string) (entry ttlCacheEntry, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now()
	entry, ok := s.cache[key]
	if !ok {
		return entry, errors.New(errMsgCacheMiss)
	}
	if now.After(entry.Timeout) {
		return entry, errors.New(errMsgCacheExpire)
	}
	return entry, nil
}

func (s *AsyncTtlCache) GcThread() {
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

func (s *AsyncTtlCache) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.cache)
}
