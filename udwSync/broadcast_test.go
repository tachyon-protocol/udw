package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBroadcast(ot *testing.T) {
	br := NewBroadcast()
	i := Int{}
	waitForIValueFn := func(waitForI int) {
		br.WaitWithCb(func() bool {
			return i.Get() >= waitForI
		})
	}
	go func() {
		i.Add(1)
		br.Broadcast()
		waitForIValueFn(2)
		i.Add(1)
		br.Broadcast()
	}()
	waitForIValueFn(1)
	i.Add(1)
	br.Broadcast()
	waitForIValueFn(3)
	udwTest.Equal(i.Get(), 3)
}

func TestBroadcast2(ot *testing.T) {
	br := NewBroadcast()
	i := Int{}
	waitForIValueFn := func(waitForI int) {
		var cv int
		for {
			if i.Get() >= waitForI {
				return
			}
			cv = br.WaitWithVersion(cv)
		}
	}
	go func() {
		i.Add(1)
		br.Broadcast()
		waitForIValueFn(2)
		i.Add(1)
		br.Broadcast()
	}()
	waitForIValueFn(1)
	i.Add(1)
	br.Broadcast()
	waitForIValueFn(3)
	udwTest.Equal(i.Get(), 3)
}

func TestBroadcastBench(ot *testing.T) {
	br := NewBroadcast()
	i := Int{}
	waitForIValueFn := func(waitForI int) {
		br.WaitWithCb(func() bool {
			return i.Get() >= waitForI
		})
	}
	udwTest.Benchmark(func() {
		const num = 10000
		udwTest.BenchmarkSetNum(num * 3)
		for l := 0; l < num; l++ {
			i.Set(0)
			go func() {
				i.Add(1)
				br.Broadcast()
				waitForIValueFn(2)
				i.Add(1)
				br.Broadcast()
			}()
			waitForIValueFn(1)
			i.Add(1)
			br.Broadcast()
			waitForIValueFn(3)
		}
	})
}
