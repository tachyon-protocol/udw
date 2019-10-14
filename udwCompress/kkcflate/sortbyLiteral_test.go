package kkcflate

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math/rand"
	"sort"
	"testing"
)

type byLiteral_sortSlice []literalNode

func (s byLiteral_sortSlice) Len() int { return len(s) }

func (data byLiteral_sortSlice) Less(i, j int) bool {
	return (data[i].literal < data[j].literal)
}

func (s byLiteral_sortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func TestT_sortbyLiteral(ot *testing.T) {
	isSortedFn := func(a []literalNode) bool {
		return sort.IsSorted(byLiteral_sortSlice(a))
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
	correctTestFn(sortbyLiteral)
	correctTestFn(func(a []literalNode) {
		byLiteral_quickSort(a, 0, len(a), 0)
	})
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			sortbyLiteral(a)
		})
	}
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			sort.Sort(byLiteral_sortSlice(a))
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
			sortbyLiteral(a)
		})
		for i := 0; i < size; i++ {
			a[i] = getNodeByIntFn(size - i)
		}
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(size)
			sort.Sort(byLiteral_sortSlice(a))
		})
	}
}
