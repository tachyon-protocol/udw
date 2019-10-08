package udwSort

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

type string_sortSlice []string

func (s string_sortSlice) Len() int { return len(s) }

func (data string_sortSlice) Less(i, j int) bool {
	return (data[i] < data[j])
}

func (s string_sortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func TestT_SortString(ot *testing.T) {
	isSortedFn := func(a []string) bool {
		return sort.IsSorted(string_sortSlice(a))
	}
	getNodeByIntFn := func(id int) (node string) {
		return strconv.Itoa(id)
	}
	getNodeListByIntListFn := func(idList ...int) (nodeList []string) {
		nodeList = make([]string, len(idList))
		for i, id := range idList {
			nodeList[i] = getNodeByIntFn(id)
		}
		return nodeList
	}
	correctTestFn := func(sortFn func(a []string)) {
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
			a = make([]string, randomSize)
			for i := 0; i < randomSize; i++ {
				a[i] = getNodeByIntFn(rand.Intn(100))
			}
			sortFn(a)
			udwTest.Ok(isSortedFn(a))
		}
		{
			const randomSize = 1024
			a = make([]string, randomSize)
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
	correctTestFn(SortString)
	correctTestFn(func(a []string) {
		string_quickSort(a, 0, len(a), 0)
	})
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			SortString(a)
		})
	}
	{
		a := getNodeListByIntListFn(4, 3, 1, 2, 5, 6)
		udwTest.BenchmarkWithRepeatNum(1<<10, func() {
			sort.Sort(string_sortSlice(a))
		})
	}
	{
		size := 1024 * 1024
		a := make([]string, size)
		for i := 0; i < size; i++ {
			a[i] = getNodeByIntFn(size - i)
		}
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(size)
			SortString(a)
		})
		for i := 0; i < size; i++ {
			a[i] = getNodeByIntFn(size - i)
		}
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(size)
			sort.Sort(string_sortSlice(a))
		})
	}
}
