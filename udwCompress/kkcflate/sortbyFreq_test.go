package kkcflate

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math/rand"
	"sort"
	"testing"
)

type byFreq_sortSlice []literalNode

func (s byFreq_sortSlice) Len() int { return len(s) }

func (s byFreq_sortSlice) Less(i, j int) bool {
	if s[i].freq == s[j].freq {
		return s[i].literal < s[j].literal
	}
	return s[i].freq < s[j].freq
}

func (s byFreq_sortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func TestT_sortbyFreq(ot *testing.T) {
	isSortedFn := func(a []literalNode) bool {
		return sort.IsSorted(byFreq_sortSlice(a))
	}
	getNodeByIntFn := func(id int) (node literalNode) {
		return literalNode{
			literal: uint16(id),
			freq:    int32(id),
		}
	}
	getNodeListByIntListFn := func(idList ...int) (nodeList []literalNode) {
		nodeList = make([]literalNode, len(idList))
		for i, id := range idList {
			nodeList[i] = getNodeByIntFn(id)
		}
		return nodeList
	}
	correctTestFn := func(sortFn func(a []literalNode)) {
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		sortFn(a)
		udwTest.Ok(isSortedFn(a))
		for size := 0; size < 1024; size++ {
			a := getNodeListByIntListFn(rand.Perm(size)...)
			sortFn(a)
			udwTest.Ok(isSortedFn(a))
			sortFn(a)
			udwTest.Ok(isSortedFn(a))
		}
		{
			const randomSize = 1000000
			a = make([]literalNode, randomSize)
			for i := 0; i < randomSize; i++ {
				a[i] = getNodeByIntFn(rand.Intn(100))
			}
			sortFn(a)
			udwTest.Ok(isSortedFn(a))
		}
		{
			const randomSize = 1024
			a = make([]literalNode, randomSize)
			for i := 0; i < randomSize; i++ {
				a[i] = getNodeByIntFn(0)
			}
			sortFn(a)
			udwTest.Ok(isSortedFn(a))
		}
		{
			var ints = getNodeListByIntListFn(74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586)
			sortFn(ints)
			udwTest.Ok(isSortedFn(a))
		}
	}
	correctTestFn(sortbyFreq)
	correctTestFn(func(a []literalNode) {
		byFreq_quickSort(a, 0, len(a), 0)
	})
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			sortbyFreq(a)
		})
	}
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			sort.Sort(byFreq_sortSlice(a))
		})
	}
	{
		size := 1024 * 1024
		a := make([]literalNode, size)
		for i := 0; i < size; i++ {
			a[i] = getNodeByIntFn(size - i)
		}
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(size)
			sortbyFreq(a)
		})
		for i := 0; i < size; i++ {
			a[i] = getNodeByIntFn(size - i)
		}
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(size)
			sort.Sort(byFreq_sortSlice(a))
		})
	}
}
