package udwCompressTest

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestCompressor(compressor func(inb []byte) (outb []byte), decompressor func(inb []byte) (outb []byte)) {
	origin := []byte("123")
	ob := compressor([]byte("123"))
	output := decompressor(ob)
	udwTest.Equal(origin, output)

	for _, i := range []int{1, 10, 100, 1000, 1e4, 1e5} {
		for _, originGetter := range []func(i int) []byte{
			func(i int) []byte {
				return udwRand.MustCryptoRandBytes(i)
			},
			func(i int) []byte {
				return bytes.Repeat([]byte{31}, i)
			},
		} {
			origin := originGetter(i)
			origin2 := udwBytes.Clone(origin)
			ob = compressor(origin2)
			output := decompressor(ob)
			udwTest.Equal(origin, origin2)
			udwTest.Equal(origin, output)
			ob2 := udwBytes.Clone(ob)
			udwTest.AssertPanic(func() {
				decompressor(ob2[:len(ob2)-1])
			})

		}
	}

}
