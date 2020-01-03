package udwSortedMap

import (
	"sort"
	"sync"
)

type Uint64ToUint64Map struct {
	lock sync.RWMutex
	m    map[uint64]uint64
}

func NewUint64ToUint64Map() *Uint64ToUint64Map {
	return &Uint64ToUint64Map{
		m: map[uint64]uint64{},
	}
}

func (sfm *Uint64ToUint64Map) Set(k uint64, v uint64) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *Uint64ToUint64Map) Del(k uint64) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *Uint64ToUint64Map) KeysByValueDesc() []uint64 {
	return sfm.Keys(DescSortByValue)
}

func (sfm *Uint64ToUint64Map) KeysByValueAsc() []uint64 {
	return sfm.Keys(AscSortByValue)
}

func (sfm *Uint64ToUint64Map) Keys(st SortType) []uint64 {
	sfm.lock.RLock()
	if len(sfm.m) == 0 {
		sfm.lock.RUnlock()
		return nil
	}
	keys := make([]uint64, 0, len(sfm.m))
	for k := range sfm.m {
		keys = append(keys, k)
	}
	if st == SortByKey {
		sort.Sort(Uint64Slice(keys))
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

func (sfm *Uint64ToUint64Map) Get(key uint64) (v uint64, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
