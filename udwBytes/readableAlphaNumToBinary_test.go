package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestReadableAlphaNumToBinaryOrOriginToWriter(ot *testing.T) {
	m := map[string]struct{}{}
	const num = int(1e3)
	for i := 0; i < num; i++ {
		id := udwRand.MustCryptoRandToReadableAlphaNum(24)
		m[id] = struct{}{}
		b := ReadableAlphaNumToBinaryOrOriginToSlice(id)
		m[string(b)] = struct{}{}
	}
	udwTest.Equal(len(m), int(2*num))

	bufW := &BufWriter{}
	id := udwRand.MustCryptoRandToReadableAlphaNum(24)
	udwTest.Benchmark(func() {
		const num = 1e3
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			bufW.Reset()
			ReadableAlphaNumToBinaryOrOriginToWriter(id, bufW)
		}
	})

}

func TestReadableAlphaNumFromBinary(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := udwRand.MustCryptoRandToReadableAlphaNum(24)
		b := ReadableAlphaNumToBinaryOrOriginToSlice(id)
		id2 := ReadableAlphaNumFromBinary(b, 24)
		udwTest.Equal(id, id2)
	}
	{
		id := udwRand.MustCryptoRandToReadableAlphaNum(24)
		b := ReadableAlphaNumToBinaryOrOriginToSlice(id)
		id2 := ReadableAlphaNumFromBinary(b, 30)
		udwTest.Equal(len(id2), 30)
	}
	{
		id := "3" + udwRand.MustCryptoRandToReadableAlphaNum(24)
		b := ReadableAlphaNumToBinaryOrOriginToSlice(id)
		id2 := ReadableAlphaNumFromBinary(b, 20)
		udwTest.Equal(len(id2), 25)
	}

	id := udwRand.MustCryptoRandToReadableAlphaNum(24)
	b := ReadableAlphaNumToBinaryOrOriginToSlice(id)
	udwTest.Benchmark(func() {
		const num = 1e3
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			ReadableAlphaNumFromBinary(b, 24)
		}
	})

}
