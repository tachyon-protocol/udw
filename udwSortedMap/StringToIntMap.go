package udwSortedMap

import (
	"sort"
	"sync"
)

type StringToIntMap struct {
	lock sync.RWMutex
	m    map[string]int
}

func NewStringToIntMap() *StringToIntMap {
	return &StringToIntMap{
		m: map[string]int{},
	}
}

func (sfm *StringToIntMap) Set(k string, v int) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *StringToIntMap) Del(k string) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *StringToIntMap) KeysByValueDesc() []string {
	return sfm.Keys(DescSortByValue)
}

func (sfm *StringToIntMap) KeysByValueAsc() []string {
	return sfm.Keys(AscSortByValue)
}

func (sfm *StringToIntMap) Keys(st SortType) []string {
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

func (sfm *StringToIntMap) Get(key string) (v int, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}
