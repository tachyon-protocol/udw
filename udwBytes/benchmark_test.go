package udwBytes

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBenchWrite(b *testing.T) {
	s := []byte("foobarbaz")
	const num = int(1e4)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		var buf BufWriter
		for a := 0; a < num; a++ {
			for i := 0; i < 100; i++ {
				buf.Write(s)
			}
			buf.Reset()
		}
	})

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		var buf2 []byte
		for a := 0; a < num; a++ {
			for i := 0; i < 100; i++ {
				buf2 = append(buf2, s...)
			}
			buf2 = buf2[:0]
		}
	})

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		var buf2 bytes.Buffer
		for a := 0; a < num; a++ {
			for i := 0; i < 100; i++ {
				buf2.Write(s)
			}
			buf2.Reset()
		}
	})
}
