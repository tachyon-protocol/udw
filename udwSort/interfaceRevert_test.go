package udwSort_test

import (
	"github.com/tachyon-protocol/udw/udwSort"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestInterfaceRevert(t *testing.T) {
	sList := []int{1, 2, 3}
	udwSort.InterfaceRevert(sList)
	udwTest.Equal(sList[0], 3)
	udwTest.Equal(sList[1], 2)
	udwTest.Equal(sList[2], 1)
	sList = []int{1, 2}
	udwSort.InterfaceRevert(sList)
	udwTest.Equal(sList[0], 2)
	udwTest.Equal(sList[1], 1)
}
