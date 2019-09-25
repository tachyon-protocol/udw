package udwMap

import "sort"

type KeyValueInt64Pair struct {
	Key   string
	Value int64
}

type keyValueInt64PairSlice []KeyValueInt64Pair

func (s keyValueInt64PairSlice) Len() int {
	return len(s)
}
func (s keyValueInt64PairSlice) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s keyValueInt64PairSlice) Less(i int, j int) bool {
	if s[i].Value == s[j].Value {
		return s[i].Key < s[j].Key
	}
	return s[i].Value > s[j].Value
}

func MapStringInt64ToStringIntPairListByIntDesc(m map[string]int64) []KeyValueInt64Pair {
	pairList := make([]KeyValueInt64Pair, len(m))
	index := 0
	for s, i := range m {
		pairList[index].Key = s
		pairList[index].Value = i
		index++
	}
	sort.Sort(keyValueInt64PairSlice(pairList))
	return pairList
}

func StringInt64PairSliceSort(pairList []KeyValueInt64Pair) {
	sort.Sort(keyValueInt64PairSlice(pairList))
}
