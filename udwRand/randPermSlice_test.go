package udwRand_test

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"testing"
)

func TestMustPermInterfaceSlice(t *testing.T) {
	{
		inS := []int{}
		outS := udwRand.MustPermInterfaceSlice(inS).([]int)
		udwTest.Equal(outS, []int{})
	}
	{
		inS := []int{1}
		outS := udwRand.MustPermInterfaceSlice(inS).([]int)
		udwTest.Equal(outS, []int{1})
	}
	{
		inS := []int{1, 2, 3}
		outS := udwRand.MustPermInterfaceSlice(inS).([]int)
		udwTest.Equal(len(outS), 3)
	}
	{
		inS := []int{1, 2, 3}
		udwTest.Equal(len(inS), 3)

		for i := 0; i < 100; i++ {
			inS = []int{1, 2, 3}
			udwRand.MustPermInterfaceSliceInPlace(inS)

		}
	}
}

func TestPermNoAlloc(t *testing.T) {
	for i := 0; i < 10; i++ {
		inS := make([]int, i)
		udwRand.PermNoAlloc(inS)
	}
	r := udwRand.NewInt64SeedUdwRand(0)
	repeatCheck := map[string]struct{}{}
	_buf := udwBytes.BufWriter{}
	arrayToFn := func(inS []int) string {
		_buf.Reset()
		for _, v := range inS {
			_buf.WriteString_(strconv.Itoa(v))
			_buf.WriteString_(",")
		}
		return _buf.GetString()
	}
	for i := 0; i < 10; i++ {
		inS := make([]int, 10)
		r.PermNoAllocNoLock(inS)
		repeatCheck[arrayToFn(inS)] = struct{}{}
	}
	udwTest.Equal(len(repeatCheck), 10)
	{
		repeatCheck := map[string]struct{}{}
		for i := 0; i < 1000; i++ {
			inS := make([]int, 5)
			r.PermNoAllocNoLock(inS)
			repeatCheck[arrayToFn(inS)] = struct{}{}
		}
		udwTest.Equal(len(repeatCheck), 120)
	}
}

func TestUdwRand_ShuffleIntArrayNoAllocNoLock(t *testing.T) {
	r := udwRand.NewInt64SeedUdwRand(0)
	posList2 := make([]int, 7)
	for i := 0; i < 7; i++ {
		posList2[i] = i
	}
	seenSet := map[string]struct{}{}
	for i := 0; i < 2000; i++ {
		r.ShuffleIntArrayNoAllocNoLock(posList2, 3)
		seenSet[udwJson.MustMarshalToString(posList2[:3])] = struct{}{}
	}
	udwTest.Equal(len(seenSet), 210)
}
