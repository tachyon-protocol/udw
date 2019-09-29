package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestPool(ot *testing.T) {

	pool := BufWriterPool{}
	for i := 0; i < 10; i++ {
		bw := pool.Get()
		udwTest.Equal(len(bw.GetBytes()), 0)
		bw.AddPos(i)
		udwTest.Equal(len(bw.GetBytes()), i)
		pool.Put(bw)
	}

	udwTest.Benchmark(func() {
		const num = 1e6
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			bw := pool.Get()
			bw.AddPos(4096)
			pool.Put(bw)
		}
	})

	udwTest.Benchmark(func() {
		const num = 1e6
		const perThread = num / 10
		udwTest.BenchmarkSetNum(num)
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				for i := 0; i < perThread; i++ {
					bw := pool.Get()
					bw.AddPos(4096)
					pool.Put(bw)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})

	udwTest.Benchmark(func() {
		const num = 1e5
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			bw := pool.Get()
			bw.AddPos(4096)
		}
	})
}
