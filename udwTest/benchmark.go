package udwTest

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

func Benchmark(fn func()) {
	gBenchmarkCtx.runLocker.Lock()
	defer gBenchmarkCtx.runLocker.Unlock()
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.num = 1
	gBenchmarkCtx.sizePerRun = 0
	gBenchmarkCtx.startCgoCall = runtime.NumCgoCall()
	gBenchmarkCtx.dataLocker.Unlock()

	var memstats1 runtime.MemStats
	var memstats2 runtime.MemStats
	runtime.ReadMemStats(&memstats1)
	startTime := time.Now()
	fn()
	dur := time.Since(startTime)
	runtime.ReadMemStats(&memstats2)
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.endCgoCall = runtime.NumCgoCall()
	gBenchmarkCtx.dur = dur
	gBenchmarkCtx.allocNum = memstats2.Mallocs - memstats1.Mallocs
	gBenchmarkCtx.allocSize = memstats2.TotalAlloc - memstats1.TotalAlloc
	resultS := benchmarkResultString__NOLOCK()
	gBenchmarkCtx.name = ""
	gBenchmarkCtx.dataLocker.Unlock()
	fmt.Println(resultS)

}

func BenchmarkWithRepeatNum(num int, fn func()) {
	Benchmark(func() {
		BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			fn()
		}
	})
}

func BenchmarkWithThread(threadNum int, perThreadNumber int, fn func()) {
	wg := &sync.WaitGroup{}
	wg.Add(threadNum)
	Benchmark(func() {
		BenchmarkSetNum(threadNum * perThreadNumber)
		for i := 0; i < threadNum; i++ {
			go func() {
				fn()
				wg.Done()
			}()
		}
		wg.Wait()
	})

}

type benchmarkCtx struct {
	num             int
	sizePerRun      int
	dur             time.Duration
	name            string
	namePadding     int
	allocNum        uint64
	allocSize       uint64
	runLocker       sync.Mutex
	dataLocker      sync.Mutex
	startCgoCall    int64
	endCgoCall      int64
	namePaddingSize int
}

var gBenchmarkCtx benchmarkCtx

func BenchmarkSetNum(num int) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.num = num
	gBenchmarkCtx.dataLocker.Unlock()
}

func BenchmarkSetBytePerRun(sizePerRun int) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.sizePerRun = sizePerRun
	gBenchmarkCtx.dataLocker.Unlock()
}

func BenchmarkSetName(name string) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.name = name
	gBenchmarkCtx.dataLocker.Unlock()
}

func BenchmarkSetNamePadding(namePadding int) {
	gBenchmarkCtx.dataLocker.Lock()
	gBenchmarkCtx.namePadding = namePadding
	gBenchmarkCtx.dataLocker.Unlock()
}

func benchmarkResultString__NOLOCK() string {
	buf := bytes.Buffer{}
	buf.WriteString("Benchmark ")
	if gBenchmarkCtx.name != "" {
		buf.WriteString(gBenchmarkCtx.name)
		if gBenchmarkCtx.namePadding <= 1 {
			gBenchmarkCtx.namePadding = 8
		}
		if len(gBenchmarkCtx.name) <= gBenchmarkCtx.namePadding-1 {
			buf.WriteString(strings.Repeat(" ", gBenchmarkCtx.namePadding-1-len(gBenchmarkCtx.name)))
		}
		buf.WriteString(" ")
	}
	buf.WriteString("[")
	buf.WriteString(durationFormatFloat64Ns(float64(gBenchmarkCtx.dur) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	buf.WriteString("[")
	buf.WriteString(gbFromFloat64(float64(gBenchmarkCtx.num) / float64(gBenchmarkCtx.dur) * 1e9))
	buf.WriteString("op/s] ")

	buf.WriteString("duration:[")
	buf.WriteString(durationFormatFloat64Ns(float64(gBenchmarkCtx.dur)))
	buf.WriteString("] ")

	buf.WriteString("allocNum:[")
	buf.WriteString(gbFromFloat64(float64(gBenchmarkCtx.allocNum) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	buf.WriteString("allocSize:[")
	buf.WriteString(gbFromFloat64(float64(gBenchmarkCtx.allocSize) / float64(gBenchmarkCtx.num)))
	buf.WriteString("/op] ")

	if gBenchmarkCtx.sizePerRun > 0 {
		buf.WriteString("bandwidth:[")
		buf.WriteString(gbFromFloat64(float64(gBenchmarkCtx.num*gBenchmarkCtx.sizePerRun) / float64(gBenchmarkCtx.dur) * 1e9))
		buf.WriteString("/s] ")
	}
	if gBenchmarkCtx.endCgoCall-gBenchmarkCtx.startCgoCall > 0 {
		buf.WriteString("cgoNum:[")
		buf.WriteString(gbFromFloat64(float64(gBenchmarkCtx.endCgoCall-gBenchmarkCtx.startCgoCall) / float64(gBenchmarkCtx.num)))
		buf.WriteString("/op] ")
	}
	return buf.String()
}
