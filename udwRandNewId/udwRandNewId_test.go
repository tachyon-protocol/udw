package udwRandNewId_test

import (
	"github.com/tachyon-protocol/udw/udwRandNewId"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestUdwRandNewId(ot *testing.T) {
	id := udwRandNewId.NewId()
	udwTest.Equal(len(id), 24)
	idMap := map[string]struct{}{}
	num := int(1e5)
	for i := 0; i < num; i++ {
		id := udwRandNewId.NewId()
		idMap[id] = struct{}{}
	}
	udwTest.Equal(len(idMap), num)
}

func TestBenchNewId(ot *testing.T) {
	udwRandNewId.NewId()
	const runNumber = 1e4
	udwTest.BenchmarkWithRepeatNum(runNumber, func() {
		udwRandNewId.NewId()
	})
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(runNumber)
		wg := sync.WaitGroup{}
		wg.Add(4)
		for i := 0; i < 4; i++ {
			go func() {
				for i := 0; i < runNumber/4; i++ {
					udwRandNewId.NewId()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
	tasker := udwTask.New(10)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(runNumber)
		for i := 0; i < runNumber; i++ {
			tasker.AddFunc(func() {
				udwRandNewId.NewId()
			})
		}
		tasker.Close()
	})
}
