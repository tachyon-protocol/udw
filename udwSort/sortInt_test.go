package udwSort

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math/rand"
	"sort"
	"testing"
)

func TestSortInt(t *testing.T) {
	correctTestFn := func(sortFn func(a []int)) {
		a := []int{4, 3, 1, 2, 5, 6}
		sortFn(a)
		for i, v := range a {
			udwTest.Ok(v == i+1)
		}
		udwTest.Ok(sort.IntsAreSorted(a))
		for size := 0; size < 1024; size++ {
			a := rand.Perm(size)
			sortFn(a)
			udwTest.Ok(sort.IntsAreSorted(a))
			for i, v := range a {
				udwTest.Ok(v == i)
			}
			sortFn(a)
			udwTest.Ok(sort.IntsAreSorted(a))
		}
		{
			const randomSize = 1000000
			a = make([]int, randomSize)
			for i := 0; i < randomSize; i++ {
				a[i] = rand.Intn(100)
			}
			sortFn(a)
			udwTest.Ok(sort.IntsAreSorted(a))
		}
		{
			const randomSize = 1024
			a = make([]int, randomSize)
			for i := 0; i < randomSize; i++ {
				a[i] = 0
			}
			sortFn(a)
			udwTest.Ok(sort.IntsAreSorted(a))
		}
		{
			var ints = []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
			sortFn(ints)
			udwTest.Ok(sort.IntsAreSorted(a))
		}
	}
	correctTestFn(SortInt)
	correctTestFn(func(a []int) {
		int_quickSort(a, 0, len(a), 0)
	})

	a := []int{4, 3, 1, 2, 5, 6}
	udwTest.BenchmarkWithRepeatNum(1<<10, func() {
		SortInt(a)
	})
	udwTest.BenchmarkWithRepeatNum(1<<10, func() {
		sort.Ints(a)
	})

	size := 1024 * 1024
	a = make([]int, size)
	for i := 0; i < size; i++ {
		a[i] = size - i
	}

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		SortInt(a)
	})
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		sort.Ints(a)
	})
}
