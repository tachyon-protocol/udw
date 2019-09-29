package udwRand

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestIntBetween(t *testing.T) {
	udwTest.Equal(IntBetween(0, 0), 0)
	hasView := [2]bool{}
	for i := 0; i < 100; i++ {
		ret := IntBetween(0, 1)
		udwTest.Ok(ret == 0 || ret == 1)
		hasView[ret] = true
	}
	udwTest.Equal(hasView[0], true)
	udwTest.Equal(hasView[1], true)
}

func TestLock(t *testing.T) {
	r := MustNewCryptSeedUdwRand()
	r.Float64()
	r.Int()
	r.Intn(1)
	r.Int63Between(0, 10)
	r.IntBetween(0, 10)
	r.Float64Between(0, 10)
	r.TimeDurationBetween(0, 10)
	r.Float32()
	r.MulitChoice(10, 5)
	r.MulitChoiceOriginOrder(10, 5)
	r.Perm(10)
	r.HappendBaseOnPossibility(0.5)
	r.ChoiceFromIntSlice([]int{1, 2})
	r.PermIntSlice([]int{1, 2})

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			IntBetween(0, 100)
			Float64()
			Intn(1)
			Perm(100)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIntSliceBetween(t *testing.T) {
	Intn(1)
	const repeat = 5 << 10
	const max = 2 << 10
	indexArray := make([]int, max)
	doRandom := false
	outputArray := make([]int, 7)
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		if !doRandom {
			PermNoAlloc(indexArray)
			doRandom = true
		}
		PermFromCacheNoAlloc(indexArray, 1, 1500, outputArray)
	})
	fmt.Println(outputArray)

	var v1Result []int
	udwTest.BenchmarkWithRepeatNum(repeat, func() {
		v1Result = IntSliceBetween(1, 1500, 7)
	})
	fmt.Println(v1Result)
}
