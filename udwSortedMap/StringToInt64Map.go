package udwSortedMap

import (
	"sort"
	"sync"
)

type StringToInt64Map struct {
	lock sync.RWMutex
	m    map[string]int64
}

func NewStringToInt64Map() *StringToInt64Map {
	return &StringToInt64Map{
		m: map[string]int64{},
	}
}

func (sfm *StringToInt64Map) Set(k string, v int64) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *StringToInt64Map) Del(k string) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *StringToInt64Map) KeysByValueDesc() []string {
	return sfm.Keys(DescSortByValue)
}

func (sfm *StringToInt64Map) KeysByValueAsc() []string {
	return sfm.Keys(AscSortByValue)
}

func (sfm *StringToInt64Map) Keys(st SortType) []string {
	sfm.lock.RLock()
	if len(sfm.m) == 0 {
		sfm.lock.RUnlock()
		return nil
	}
	keys := make([]string, 0, len(sfm.m))
	for k := range sfm.m {
		keys = append(keys, k)
	}
	if st == SortByKey {
		sort.Sort(StringSlice(keys))
		sfm.lock.RUnlock()
		return keys
	}
	sort.Slice(keys, func(a int, b int) bool {
		r := sfm.m[keys[a]] > sfm.m[keys[b]]
		if st == AscSortByValue {
			return !r
		}
		return r
	})
	sfm.lock.RUnlock()
	return keys
}

func (sfm *StringToInt64Map) Get(key string) (v int64, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}
