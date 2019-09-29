package udwRand

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustCryptoRand(ot *testing.T) {
	for _, f := range []func(length int) string{
		MustCryptoRandToHex,
		MustCryptoRandToReadableAlphaNum,
		MustCryptoRandToReadableAlphaNum,
	} {
		ret := f(15)
		fmt.Println(ret)
		udwTest.Equal(len(ret), 15)

		ret = f(1)
		udwTest.Equal(len(ret), 1)

		ret = f(20)
		udwTest.Equal(len(ret), 20)
	}
}
