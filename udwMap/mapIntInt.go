package udwMap

import (
	"github.com/tachyon-protocol/udw/udwSort"
)

type IntIntPair struct {
	K int `json:",omitempty"`
	V int `json:",omitempty"`
}

func MapIntIntToPairListByKeyIntDesc(m map[int]int) []IntIntPair {
	pairList := make([]IntIntPair, len(m))
	index := 0
	for s, i := range m {
		pairList[index].K = s
		pairList[index].V = i
		index++
	}
	udwSort.InterfaceCallbackSortWithIndexLess(pairList, func(a int, b int) bool {
		return pairList[a].K > pairList[b].K
	})
	return pairList
}

func MapIntIntToPairListByKeyIntAes(m map[int]int) []IntIntPair {
	pairList := make([]IntIntPair, len(m))
	index := 0
	for s, i := range m {
		pairList[index].K = s
		pairList[index].V = i
		index++
	}
	udwSort.InterfaceCallbackSortWithIndexLess(pairList, func(a int, b int) bool {
		return pairList[a].K < pairList[b].K
	})
	return pairList
}
