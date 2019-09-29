package udwTask

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimitThreadTaskManagerV2(t *testing.T) {
	tm := NewLimitThreadTaskManagerV2(1)
	result_chan := make(chan int)
	var result_locker sync.Mutex
	var result_var int
	f := func() {
		result_chan <- 1
		result_locker.Lock()
		result_var = 2
		result_locker.Unlock()
	}
	tm.AddFunc(f)
	udwTest.Equal(<-result_chan, 1, "not run result_chan not match")
	tm.AddFunc(f)
	udwTest.Equal(<-result_chan, 1, "not run result_chan not match")
	tm.Close()
	udwTest.Equal(result_var, 2, "not run result_var not match")

	tm = NewLimitThreadTaskManagerV2(1)
	tm.SetThreadNum(10)
	a := uint64(0)
	f = func() {
		atomic.AddUint64(&a, 1)
	}
	for i := 0; i < 100; i++ {
		tm.AddFunc(f)
	}
	tm.Close()
	udwTest.Equal(a, uint64(100))
}

func TestNewLimitThreadTaskManagerV2NoBuffer(t *testing.T) {
	tm := NewLimitThreadTaskManagerV2NoBuffer(1)
	start := time.Now()
	tm.AddFunc(func() {
		time.Sleep(10 * time.Millisecond)
	})
	udwTest.Ok(time.Since(start) < time.Millisecond)
	tm.AddFunc(func() {
		time.Sleep(20 * time.Millisecond)
	})
	dur := time.Since(start)
	udwTest.Ok(dur > 10*time.Millisecond)
	udwTest.Ok(dur < 30*time.Millisecond)
	tm.Close()
	udwTest.Ok(time.Since(start) > 30*time.Millisecond)

	start = time.Now()
	tm = NewLimitThreadTaskManagerV2NoBuffer(10)
	for i := 0; i < 100; i++ {
		tm.AddFunc(func() {
			time.Sleep(time.Millisecond * 10)
		})
	}
	tm.Close()
	udwTest.Ok(time.Since(start) >= 100*time.Millisecond)
	udwTest.Ok(time.Since(start) < 200*time.Millisecond)
}

func TestBenchNewLimitThreadTaskManagerV2NoBuffer(t *testing.T) {
	const runNumber = 1e5
	tasker := NewLimitThreadTaskManagerV2NoBuffer(10)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(runNumber)
		for i := 0; i < runNumber; i++ {
			tasker.AddFunc(func() {
			})
		}
		tasker.Close()
	})
}

func emptyLimitThreadTaskManagerV2() {
	tm := NewLimitThreadTaskManagerV2NoBuffer(10)
	tm.AddFunc(func() {

	})
}
