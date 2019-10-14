package udwFlate

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCompress/udwCompressTest"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTask"
	"github.com/tachyon-protocol/udw/udwTest"
	"runtime"
	"runtime/debug"
	"testing"
)

func TestCorrect(t *testing.T) {
	udwCompressTest.TestCompressor(FlateMustCompress, MustFlateUnCompress)
	{

		runOnceFn := func() {
			inb := make([]byte, 1024*1024)
			FlateUnCompress(inb)
		}
		debug.FreeOSMemory()
		memory1 := GetMemeoryUsage()
		for i := 0; i < 10; i++ {
			runOnceFn()
			debug.FreeOSMemory()
			memory2 := GetMemeoryUsage()
			udwTest.Ok((int64(memory2)-int64(memory1)) < 100*1024, udwStrconv.FormatGbAndIntFromInt64(int64(memory2-memory1)))
		}
	}
}

func GetMemeoryUsage() uint64 {
	memStat := &runtime.MemStats{}
	runtime.ReadMemStats(memStat)
	return memStat.Alloc
}

func TestTryCompress(t *testing.T) {

	testWithB := func(inB []byte) {
		outB := TryCompress(inB)
		udwTest.Ok(len(outB) <= len(inB)+1)
		inB2, err := TryUncompress(outB)
		udwTest.Equal(err, nil)
		udwTest.Equal(inB2, inB)

		bw1 := &udwBytes.BufWriter{}
		TryCompressToBufW(inB, bw1)
		outB2 := bw1.GetBytes()
		udwTest.Ok(len(outB2) <= len(inB)+1)

		inB2, err = TryUncompress(outB2)
		udwTest.Equal(err, nil)
		udwTest.Equal(inB2, inB)
	}
	const size = 1025
	for i := 0; i < size; i++ {
		inB := bytes.Repeat([]byte{0}, i)
		testWithB(inB)
	}
	for i := 10; i < size; i++ {
		inS := udwRand.MustCryptoRandToReadableAlpha(i)
		testWithB([]byte(inS))
	}
	const runNum = 1e3
	inB := bytes.Repeat([]byte{0}, 1025)
	bw1 := &udwBytes.BufWriter{}
	TryCompressToBufW(inB, bw1)
	compressed1k := bw1.GetBytesClone()

	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryCompressToBufW1k")
		for i := 0; i < num; i++ {
			bw1.Reset()
			TryCompressToBufW(inB, bw1)
		}
	})

	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryCompress1k")
		for i := 0; i < num; i++ {
			TryCompress(inB)
		}
	})

	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryUncompressToBufW1k")
		for i := 0; i < num; i++ {
			bw1.Reset()
			TryUncompressWithTmpBuf(compressed1k, bw1)
		}
	})

	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryUnCompress1k")
		for i := 0; i < num; i++ {
			TryUncompress(compressed1k)
		}
	})

	const threadNum = 8
	tasker := udwTask.New(threadNum)
	bufWritePool := udwBytes.BufWriterPool{}
	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryCompressToBufW1kMT")
		fn := func() {
			bw := bufWritePool.Get()
			TryCompressToBufW(inB, bw)
			bufWritePool.Put(bw)
		}
		for i := 0; i < num; i++ {
			tasker.AddFunc(fn)
		}
		tasker.Close()
	})

	tasker = udwTask.New(threadNum)
	udwTest.Benchmark(func() {
		const num = runNum
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1025)
		udwTest.BenchmarkSetName("TryUncompressToBufW1kMT")
		fn := func() {
			bw := bufWritePool.Get()
			TryUncompressWithTmpBuf(compressed1k, bw)
			bufWritePool.Put(bw)
		}
		for i := 0; i < num; i++ {
			tasker.AddFunc(fn)
		}
		tasker.Close()
	})
}
