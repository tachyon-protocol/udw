package udwSortedMap

import (
	"sort"
	"sync"
)

type StringToFloat64Map struct {
	lock sync.RWMutex
	m    map[string]float64
}

func NewStringToFloat64Map() *StringToFloat64Map {
	return &StringToFloat64Map{
		m: map[string]float64{},
	}
}

func (sfm *StringToFloat64Map) Set(k string, v float64) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *StringToFloat64Map) Del(k string) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *StringToFloat64Map) KeysByValueDesc() []string {
	return sfm.Keys(DescSortByValue)
}

func (sfm *StringToFloat64Map) KeysByValueAsc() []string {
	return sfm.Keys(AscSortByValue)
}

func (sfm *StringToFloat64Map) Keys(st SortType) []string {
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

func (sfm *StringToFloat64Map) Get(key string) (v float64, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}

type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
