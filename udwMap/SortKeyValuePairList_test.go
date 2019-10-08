package udwMap

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sort"
	"strconv"
	"testing"
)

type stringKeyValueSlice []KeyValuePair

func (s stringKeyValueSlice) Len() int {
	return len(s)
}
func (s stringKeyValueSlice) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s stringKeyValueSlice) Less(i int, j int) bool {
	return s[i].Key < s[j].Key
}

func TestSortKeyValuePairList(t *testing.T) {
	kvPairList := []KeyValuePair{
		{"5", "v5"},
		{"1", "v1"},
		{"3", "v3"},
		{"4", "v4"},
		{"1", "v1"},
	}
	SortKeyValuePairList(kvPairList)
	udwTest.Equal(len(kvPairList), 5)
	udwTest.Equal(kvPairList[0].Key, "1")
	udwTest.Equal(kvPairList[1].Key, "1")
	udwTest.Equal(kvPairList[2].Key, "3")
	udwTest.Equal(kvPairList[3].Key, "4")
	udwTest.Equal(kvPairList[4].Key, "5")

	const size = 1024
	kvPairList = make([]KeyValuePair, size)
	for i := 0; i < size; i++ {
		kvPairList[i].Key = strconv.Itoa(size - i)
		kvPairList[i].Value = strconv.Itoa(size - i)
	}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		sort.Sort(stringKeyValueSlice(kvPairList))
	})
	for i := 0; i < size; i++ {
		kvPairList[i].Key = strconv.Itoa(size - i)
		kvPairList[i].Value = strconv.Itoa(size - i)
	}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		SortKeyValuePairList(kvPairList)
	})
}
