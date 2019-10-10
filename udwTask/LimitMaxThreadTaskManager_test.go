package udwTask

import (
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	const tick = time.Millisecond * 50
	{
		startTime := time.Now()
		tasker := New(1)
		var runNumber udwSync.Int
		tasker.AddFunc(func() {
			time.Sleep(tick)
			runNumber.Add(1)
		})
		dur := time.Since(startTime)
		udwTest.Ok(dur < tick, dur.String())
		tasker.Close()
		udwTest.Equal(runNumber.Get(), 1)
		dur = time.Since(startTime)
		udwTest.Ok(dur > tick, dur.String())
		udwTest.Ok(dur < tick*2, dur.String())
	}
	{
		startTime := time.Now()
		tasker := New(10)
		tasker.Close()
		udwTest.Ok(time.Since(startTime) < tick)
	}
	{
		startTime := time.Now()
		tasker := New(10)
		var runNumber udwSync.Int
		for i := 0; i < 11; i++ {
			tasker.AddFunc(func() {
				time.Sleep(tick)
				runNumber.Add(1)
			})
		}
		udwTest.Ok(time.Since(startTime) > tick)
		udwTest.Ok(time.Since(startTime) < tick*2)
		tasker.Close()
		udwTest.Equal(runNumber.Get(), 11)
		udwTest.Ok(time.Since(startTime) > tick*2)
		udwTest.Ok(time.Since(startTime) < tick*3)
	}
	{
		startTime := time.Now()
		tasker := New(10)
		var runNumber udwSync.Int
		for i := 0; i < 21; i++ {
			tasker.AddFunc(func() {
				time.Sleep(tick)
				runNumber.Add(1)
			})
		}
		udwTest.Ok(time.Since(startTime) > tick*2)
		udwTest.Ok(time.Since(startTime) < tick*3)
		tasker.WaitAndNotClose()
		udwTest.Equal(runNumber.Get(), 21)
		udwTest.Ok(time.Since(startTime) > tick*3)
		udwTest.Ok(time.Since(startTime) < tick*4)
		for i := 0; i < 11; i++ {
			tasker.AddFunc(func() {
				time.Sleep(tick)
				runNumber.Add(1)
			})
		}
		udwTest.Ok(time.Since(startTime) > tick*4)
		udwTest.Ok(time.Since(startTime) < tick*5)
		tasker.Close()
		dur := time.Since(startTime)
		udwTest.Ok(dur > tick*5, dur)
		udwTest.Ok(dur < tick*6, dur)
	}
	{
		startTime := time.Now()
		tasker := New(-1)
		var runNumber udwSync.Int
		tasker.AddFunc(func() {
			time.Sleep(tick)
			runNumber.Add(1)
		})
		time.Sleep(tick)
		tasker.Close()
		udwTest.Ok(time.Since(startTime) > tick*2)
		udwTest.Ok(time.Since(startTime) < tick*3)
		udwTest.Equal(runNumber.Get(), 1)
	}
	{
		startTime := time.Now()
		tasker := New(1)
		var runNumber udwSync.Int
		tasker.AddFunc(func() {
			time.Sleep(tick)
			runNumber.Add(1)
		})
		time.Sleep(tick)
		tasker.Close()
		udwTest.Ok(time.Since(startTime) > tick)
		udwTest.Ok(time.Since(startTime) < tick*2)
		udwTest.Equal(runNumber.Get(), 1)
	}
	{
		startTime := time.Now()
		tasker := New(-1)
		var runNumber udwSync.Int
		tasker.AddFuncSync(func() {
			time.Sleep(tick)
			runNumber.Add(1)
		})
		tasker.Close()
		udwTest.Ok(time.Since(startTime) > tick)
		udwTest.Ok(time.Since(startTime) < tick*2)
		udwTest.Equal(runNumber.Get(), 1)
	}
	{
		startTime := time.Now()
		tasker := New(1)
		var runNumber udwSync.Int
		tasker.AddFuncSync(func() {
			time.Sleep(tick)
			runNumber.Add(1)
		})
		tasker.Close()
		udwTest.Ok(time.Since(startTime) > tick)
		udwTest.Ok(time.Since(startTime) < tick*2)
		udwTest.Equal(runNumber.Get(), 1)
	}
}

func TestBenchNew(t *testing.T) {
	udwTest.BenchmarkWithRepeatNum(1000, func() {
		tasker := New(10)
		tasker.Close()
	})
	fn := func() {
		tasker := New(10)
		for i := 0; i < 1000; i++ {
			tasker.AddFunc(func() {})
		}
		tasker.Close()
	}
	udwTest.BenchmarkWithRepeatNum(1000, fn)

	fn = func() {
		tasker := New(10)
		for i := 0; i < 5; i++ {
			tasker.AddFunc(func() {})
		}
		tasker.Close()
	}
	udwTest.BenchmarkWithRepeatNum(1000, fn)

}
