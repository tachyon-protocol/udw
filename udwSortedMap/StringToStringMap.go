package udwSortedMap

import (
	"sort"
	"sync"
)

type StringToStringMap struct {
	lock sync.RWMutex
	m    map[string]string
}

func NewStringToStringMap() *StringToStringMap {
	return &StringToStringMap{
		m: map[string]string{},
	}
}

func (sfm *StringToStringMap) Set(k string, v string) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *StringToStringMap) Del(k string) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *StringToStringMap) KeysByValueDesc() []string {
	return sfm.Keys(DescSortByValue)
}

func (sfm *StringToStringMap) KeysByValueAsc() []string {
	return sfm.Keys(AscSortByValue)
}

func (sfm *StringToStringMap) Keys(st SortType) []string {
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

func (sfm *StringToStringMap) Get(key string) (v string, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}
