package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestBenchBroadcastWaitGroup(ot *testing.T) {
	const num = 100
	var wgList []*sync.WaitGroup
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		wgList = make([]*sync.WaitGroup, num)
		for i := 0; i < num; i++ {
			thisWg := sync.WaitGroup{}
			thisWg.Add(1)
			wgList[i] = &thisWg
			go func() {
				thisWg.Wait()
			}()
		}
	})
	time.Sleep(time.Millisecond * 10)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			wgList[i].Done()
		}
	})
}

func TestBenchBroadcastChannel(ot *testing.T) {
	const num = 100
	var wgList []chan struct{}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		wgList = make([]chan struct{}, num)
		for i := 0; i < num; i++ {
			thisWg := make(chan struct{})
			wgList[i] = thisWg
			go func() {
				<-thisWg
			}()
		}
	})

	time.Sleep(time.Millisecond * 10)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			close(wgList[i])

		}
	})
}

func TestBenchBroadcastChannel2(ot *testing.T) {
	const num = 300
	var wgList []bool
	broadcastChannel := make(chan struct{})
	broadcastChannelLocker := sync.Mutex{}
	wgListLocker := sync.Mutex{}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		wgList = make([]bool, num)
		for i := 0; i < num; i++ {
			i := i
			go func() {
				for {
					broadcastChannelLocker.Lock()
					thisBroadcastChannel := broadcastChannel
					broadcastChannelLocker.Unlock()
					<-thisBroadcastChannel
					wgListLocker.Lock()
					if wgList[i] {
						wgListLocker.Unlock()
						return
					}
				}
			}()
		}
	})
	runtime.GC()
	time.Sleep(time.Millisecond * 10)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		wgListLocker.Lock()
		for i := 0; i < num; i++ {
			wgList[i] = true
		}
		wgListLocker.Unlock()
		broadcastChannelLocker.Lock()
		thisBroadcastChannel := broadcastChannel
		close(thisBroadcastChannel)
		broadcastChannel = make(chan struct{})
		broadcastChannelLocker.Unlock()
	})
}
