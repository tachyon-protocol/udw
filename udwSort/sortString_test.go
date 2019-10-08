package udwSort

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sort"
	"strconv"
	"testing"
)

func TestSortString(t *testing.T) {
	correctTestFn := func(sortFn func(a []string)) {
		a := []string{"4", "3", "1", "2", "5", "6"}
		sortFn(a)
		for i, v := range a {
			udwTest.Ok(v == strconv.Itoa(i+1))
		}
		udwTest.Ok(sort.StringsAreSorted(a))
		size := 1024
		a = make([]string, size)
		for i := 0; i < size; i++ {
			a[i] = strconv.Itoa(size - i)
		}
		sortFn(a)
		udwTest.Ok(sort.StringsAreSorted(a))
	}
	correctTestFn(SortString)
	correctTestFn(func(a []string) {
		string_quickSort(a, 0, len(a), 0)
	})

	a := []string{"4", "3", "1", "2", "5", "6"}
	SortString(a)
	for i, v := range a {
		udwTest.Ok(v == strconv.Itoa(i+1))
	}
	udwTest.BenchmarkWithRepeatNum(1<<10, func() {
		SortString(a)
	})
	udwTest.BenchmarkWithRepeatNum(1<<10, func() {
		sort.Strings(a)
	})

	size := 1024 * 102
	a = make([]string, size)
	for i := 0; i < size; i++ {
		a[i] = strconv.Itoa(size - i)
	}

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		SortString(a)
	})
	for i := 0; i < size; i++ {
		a[i] = strconv.Itoa(size - i)
	}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(size)
		sort.Strings(a)
	})
}
