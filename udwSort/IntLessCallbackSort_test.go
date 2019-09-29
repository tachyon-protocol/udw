package udwSort_test

import (
	"github.com/tachyon-protocol/udw/udwSort"
	"github.com/tachyon-protocol/udw/udwTest"
	"reflect"
	"testing"
)

func TestIntLessCallbackSort(t *testing.T) {
	a := []int{2, 1, 3, 0}
	udwSort.IntLessCallbackSort(a, func(i int, j int) bool {
		return a[i] < a[j]
	})
	udwTest.Equal(a, []int{0, 1, 2, 3})
}

func TestReverseStringSort(t *testing.T) {
	dlist := []string{
		"1",
		"3",
		"2",
	}
	udwSort.ReverseStringSort(dlist)
	udwTest.Equal(dlist, []string{"3", "2", "1"})

}

func TestReverseStringList(t *testing.T) {
	origin := []string{
		`1`, `2`, `6`, `3`,
	}
	udwSort.ReverseStringList(origin)
	udwTest.Ok(reflect.DeepEqual(origin, []string{
		`3`, `6`, `2`, `1`,
	}))
	{
		dlist := []string{
			"1",
		}
		udwSort.ReverseStringList(dlist)
		udwTest.Equal(dlist, []string{"1"})
	}
	{
		dlist := []string{
			"2",
			"1",
		}
		udwSort.ReverseStringList(dlist)
		udwTest.Equal(dlist, []string{"1", "2"})
	}
	{
		dlist := []string{}
		udwSort.ReverseStringList(dlist)
		udwTest.Equal(dlist, []string{})
	}
	{
		dlist := []string{
			"2",
			"1",
			"3",
		}
		udwSort.ReverseStringList(dlist)
		udwTest.Equal(dlist, []string{"3", "1", "2"})
	}
}
