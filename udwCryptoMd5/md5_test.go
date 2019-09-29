package udwCryptoMd5

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"testing"
)

func TestMd5Hash(ot *testing.T) {
	udwTest.Equal(Md5Hex([]byte{}), "d41d8cd98f00b204e9800998ecf8427e")
	udwTest.Equal(Md5Hex([]byte{0}), "93b885adfe0da089cdf634904fd59f71")

	const benNum = 1e4
	for _, size := range []int{
		10, 100, 1000,
	} {
		inputBuf := bytes.Repeat([]byte{0}, size)
		udwTest.Benchmark(func() {
			udwTest.BenchmarkSetNum(benNum)
			udwTest.BenchmarkSetBytePerRun(size)
			udwTest.BenchmarkSetName("md5Hex" + strconv.Itoa(size))
			for i := 0; i < benNum; i++ {
				Md5Hex(inputBuf)
			}
		})
	}
}
