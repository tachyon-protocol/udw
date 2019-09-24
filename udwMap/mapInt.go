package udwMap

import (
	"bytes"
	"sort"
	"strconv"
)

type StringIntPair struct {
	S string
	I int
}

type stringIntPairSlice []StringIntPair

func (s stringIntPairSlice) Len() int {
	return len(s)
}
func (s stringIntPairSlice) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s stringIntPairSlice) Less(i int, j int) bool {
	if s[i].I == s[j].I {
		return s[i].S < s[j].S
	}
	return s[i].I > s[j].I
}

func MapStringIntToStringIntPairListByIntDesc(m map[string]int) []StringIntPair {
	pairList := make([]StringIntPair, len(m))
	index := 0
	for s, i := range m {
		pairList[index].S = s
		pairList[index].I = i
		index++
	}
	sort.Sort(stringIntPairSlice(pairList))
	return pairList
}

func MapStringIntGetKeyList(m map[string]int) (list []string) {
	for key := range m {
		list = append(list, key)
	}
	return
}

func MapStringIntGetKeyListByIntDesc(m map[string]int) (list []string) {
	pairList := MapStringIntToStringIntPairListByIntDesc(m)
	for _, pair := range pairList {
		list = append(list, pair.S)
	}
	return list
}

func DebugStringIntPairList(pairList []StringIntPair) string {
	_buf := bytes.Buffer{}
	for _, pair := range pairList {
		_buf.WriteString(pair.S)
		_buf.WriteByte(' ')
		_buf.WriteString(strconv.Itoa(pair.I))
		_buf.WriteByte('\n')
	}
	return _buf.String()
}

func StringIntPairSliceSortByIntDesc(pairList []StringIntPair) {
	sort.Sort(stringIntPairSlice(pairList))
}

func StringIntPairSortByString(pairList []StringIntPair) {
	sort.Slice(pairList, func(i int, j int) bool {
		return pairList[i].S < pairList[j].S
	})
}

func StringIntPairSortByStringDesc(pairList []StringIntPair) {
	sort.Slice(pairList, func(i int, j int) bool {
		return pairList[i].S > pairList[j].S
	})
}
