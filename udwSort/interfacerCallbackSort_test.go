package udwSort_test

import (
	"github.com/tachyon-protocol/udw/udwSort"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

type testStruct struct {
	Name  string
	Score int
}

func TestInterfaceCallbackSort(t *testing.T) {
	sList := []string{"2", "1"}
	udwSort.InterfaceCallbackSort(sList, func(a string, b string) bool {
		return a < b
	})
	udwTest.Equal(sList[0], "1")
	udwTest.Equal(sList[1], "2")

	ssList := []testStruct{
		{Name: "a", Score: 2},
		{Name: "b", Score: 1},
	}
	udwSort.InterfaceCallbackSort(ssList, func(a testStruct, b testStruct) bool {
		return a.Score < b.Score
	})
	udwTest.Equal(ssList[0].Name, "b")
	udwTest.Equal(ssList[1].Name, "a")
}
