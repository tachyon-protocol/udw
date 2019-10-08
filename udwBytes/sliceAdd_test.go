package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestSliceAdd(t *testing.T) {
	b := []byte{0, 0, 0}
	SliceAddLittleEndianFixLen(b, 1)
	udwTest.Equal(b, []byte{1, 0, 0})
	SliceAddLittleEndianFixLen(b, 255)
	udwTest.Equal(b, []byte{0, 1, 0})
	SliceAddLittleEndianFixLen(b, 256*256)
	udwTest.Equal(b, []byte{0, 1, 1})
	SliceAddLittleEndianFixLen(b, 256*256*256)
	udwTest.Equal(b, []byte{0, 1, 1})
	b = []byte{255, 255, 255}
	SliceAddLittleEndianFixLen(b, 1)
	udwTest.Equal(b, []byte{0, 0, 0})

	b = []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	udwTest.Benchmark(func() {
		const num = 1e7
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			SliceAddLittleEndianFixLen(b, 1)
		}
	})
	udwTest.Benchmark(func() {
		const num = 1e7
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			SliceAddLittleEndianFixLen(b, 4096)
		}
	})
}
