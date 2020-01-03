package udwSortedMap_Build

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"strings"
)

func DsnBuild() {
	for _, kv := range []kvType{
		{
			k: "string",
			v: "float64",
		},
		{
			k: "string",
			v: "int",
		},
		{
			k: "string",
			v: "string",
		},
		{
			k: "string",
			v: "int64",
		},
		{
			k: "uint64",
			v: "uint64",
		},
	} {
		udwFile.MustWriteFile(
			"src/github.com/tachyon-protocol/udw/udwSortedMap/"+kv.getMapTypeName()+".go",
			kv.getCodeNumberValue())
	}
}

var gGeneratedKeySliceNameMap = map[string]bool{}

type kvType struct {
	k string
	v string
}

func (t kvType) getMapTypeName() string {
	return strings.Title(t.k) + `To` + strings.Title(t.v) + `Map`
}

func (t kvType) getKeySliceTypeName() string {
	return strings.Title(t.k) + `Slice`
}

func (t kvType) getCodeNumberValue() []byte {
	var (
		mapTypeName      = t.getMapTypeName()
		mapTypeLiteral   = `map[` + t.k + `]` + t.v
		keySliceTypeName = t.getKeySliceTypeName()
	)
	buf := udwBytes.NewBufWriter(nil)
	buf.WriteString_(`package udwSortedMap

import (
	"sync"
	"sort"
)

type ` + mapTypeName + ` struct {
	lock sync.RWMutex
	m    ` + mapTypeLiteral + `
}

func New` + mapTypeName + `() *` + mapTypeName + ` {
	return &` + mapTypeName + `{
		m: ` + mapTypeLiteral + `{},
	}
}

func (sfm *` + mapTypeName + `) Set(k ` + t.k + `, v ` + t.v + `) {
	sfm.lock.Lock()
	sfm.m[k] = v
	sfm.lock.Unlock()
}

func (sfm *` + mapTypeName + `) Del(k ` + t.k + `) {
	sfm.lock.Lock()
	delete(sfm.m, k)
	sfm.lock.Unlock()
}

func (sfm *` + mapTypeName + `) KeysByValueDesc() []` + t.k + ` {
	return sfm.Keys(DescSortByValue)
}

func (sfm *` + mapTypeName + `) KeysByValueAsc() []` + t.k + ` {
	return sfm.Keys(AscSortByValue)
}

func (sfm *` + mapTypeName + `) Keys(st SortType) []` + t.k + ` {
	sfm.lock.RLock()
	if len(sfm.m) == 0 {
		sfm.lock.RUnlock()
		return nil
	}
	keys := make([]` + t.k + `, 0, len(sfm.m))
	for k := range sfm.m {
		keys = append(keys, k)
	}
	if st == SortByKey {
		sort.Sort(` + keySliceTypeName + `(keys))
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

func (sfm *` + mapTypeName + `) Get(key ` + t.k + `) (v ` + t.v + `, ok bool) {
	sfm.lock.RLock()
	v, ok = sfm.m[key]
	sfm.lock.RUnlock()
	return v, ok
}
`)
	if !gGeneratedKeySliceNameMap[keySliceTypeName] {
		gGeneratedKeySliceNameMap[keySliceTypeName] = true
		buf.WriteString_(`
type ` + keySliceTypeName + ` []` + t.k + `

func (p ` + keySliceTypeName + `) Len() int           { return len(p) }
func (p ` + keySliceTypeName + `) Less(i, j int) bool { return p[i] < p[j] }
func (p ` + keySliceTypeName + `) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
`)
	}
	return buf.GetBytes()
}
