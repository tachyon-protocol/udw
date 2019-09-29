package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
	"time"
)

func TestMemoryObjectPool(ot *testing.T) {
	allocNum := 0
	allocNumLocker := sync.Mutex{}
	pool := NewMemoryObjectPool(MemoryObjectPoolRequest{
		MaxEntryNumber: 4,
		NewObjectFn: func() interface{} {
			allocNumLocker.Lock()
			allocNum++
			allocNumLocker.Unlock()
			return make([]byte, 1024*1024)
		},
	})
	hasRun := false
	for i := 0; i < 5; i++ {
		pool.Go(func(obj interface{}) {
			hasRun = true
			b := obj.([]byte)
			udwTest.Equal(len(b), 1024*1024)
		})
	}
	udwTest.Equal(hasRun, true)
	udwTest.Equal(allocNum, 1)

	startTime := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			pool.Go(func(obj interface{}) {
				time.Sleep(time.Millisecond * 10)
				b := obj.([]byte)
				udwTest.Equal(len(b), 1024*1024)
				wg.Done()
			})
		}()
	}
	wg.Wait()
	dur := time.Since(startTime)
	udwTest.Equal(allocNum, 4)
	udwTest.Ok(dur >= time.Millisecond*20)
	udwTest.Ok(dur < time.Millisecond*40)
}
