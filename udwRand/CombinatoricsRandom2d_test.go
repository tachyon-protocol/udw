package udwRand

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestCombinatoricsRandom2d(t *testing.T) {
	r := NewInt64SeedUdwRand(0)
	for testcaseId, testcase := range []struct {
		randomer    *CombinatoricsRandom2d
		retLen      int
		retANumList []int
		retBNumList []int
		retHasSolve bool
	}{
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{1, 2},
				BNumList: []int{1, 2},
				ValidCombine: [][]bool{
					{true, true},
					{true, true},
				},
			},
			retLen:      3,
			retANumList: []int{1, 2},
			retBNumList: []int{1, 2},
			retHasSolve: true,
		},
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 5},
				ValidCombine: [][]bool{
					{true, false},
					{false, true},
					{true, false},
				},
			},
			retLen:      9,
			retANumList: []int{2, 3, 4},
			retBNumList: []int{4, 5},
			retHasSolve: false,
		},
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 5},
				ValidCombine: [][]bool{
					{true, true},
					{false, true},
					{true, false},
				},
			},
			retLen:      9,
			retANumList: []int{2, 3, 4},
			retBNumList: []int{4, 5},
			retHasSolve: true,
		},
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{3, 5},
				ValidCombine: [][]bool{
					{true, true},
					{false, true},
					{true, false},
				},
			},
			retLen:      8,
			retANumList: []int{2, 3, 3},
			retBNumList: []int{3, 5},
			retHasSolve: true,
		},
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 4},
				ValidCombine: [][]bool{
					{true, true},
					{false, true},
					{true, false},
				},
			},
			retLen:      8,
			retANumList: nil,
			retBNumList: []int{4, 4},
			retHasSolve: true,
		},
		{
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{10, 10, 10},
				BNumList: []int{4, 4},
				ValidCombine: [][]bool{
					{true, true},
					{false, true},
					{true, false},
				},
			},
			retLen:      8,
			retANumList: nil,
			retBNumList: []int{4, 4},
			retHasSolve: true,
		},
	} {
		for i := 0; i < 10; i++ {
			randomer := testcase.randomer
			err := randomer.Random(r)
			if !testcase.retHasSolve {
				udwTest.Ok(err != nil)
				continue
			}
			udwTest.Equal(err, nil)
			udwTest.Equal(len(randomer.Output), testcase.retLen)
			ANumList := make([]int, len(randomer.ANumList))
			BNumList := make([]int, len(randomer.BNumList))
			for _, row := range randomer.Output {
				ANumList[row.X]++
				BNumList[row.Y]++
			}

			if testcase.retANumList != nil {
				udwTest.Equal(ANumList, testcase.retANumList)
			}
			if testcase.retBNumList != nil {
				udwTest.Equal(BNumList, testcase.retBNumList, "BNumList not correct testcaseId: %d", testcaseId)
			}
		}
	}
}
