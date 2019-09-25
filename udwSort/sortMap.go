package udwSort

import (
	"reflect"
	"sort"
)

func KeysOfMapSortByKey(keyOfMapMustBeString interface{}) (keyList []string) {
	v := reflect.ValueOf(keyOfMapMustBeString)
	t := reflect.TypeOf(keyOfMapMustBeString)
	if v.Kind() != reflect.Map {
		panic("sortedKeys:accept map only")
	}
	keys := v.MapKeys()
	keyList = make([]string, 0, len(keys))
	if t.Key().Kind() != reflect.String {
		panic("sortedKeys:key must be string")
	}
	for _, key := range keys {
		keyList = append(keyList, key.String())
	}
	SortString(keyList)
	return keyList
}
func KeysOfMapSortByValue(m map[string]int) (keyList []string) {
	type kv struct {
		Key   string
		Value int
	}
	ss := make([]kv, 0, len(m))
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	keyList = make([]string, 0, len(m))
	for i, _ := range ss {
		keyList = append(keyList, ss[i].Key)
	}
	return
}
func KeysOfMapSortByValueFloat64(m map[string]float64) (keyList []string) {
	type kv struct {
		Key   string
		Value float64
	}
	ss := make([]kv, 0, len(m))
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	keyList = make([]string, 0, len(m))
	for i, _ := range ss {
		keyList = append(keyList, ss[i].Key)
	}
	return
}
