package udwDebug_test

import (
	"encoding/hex"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestHexDumpWithAddrOffset(t *testing.T) {
	udwTest.Equal(udwDebug.HexDumpWithAddrOffset([]byte("1234567890123456"), 0), `00000000  31 32 33 34 35 36 37 38  39 30 31 32 33 34 35 36  |1234567890123456|
`)
	udwTest.Equal(udwDebug.HexDumpWithAddrOffset([]byte("1234567890123456"), 0), hex.Dump([]byte("1234567890123456")))
	for i := 1; i < 48; i++ {
		b := udwRand.MustCryptoRandBytes(i)
		udwTest.Equal(udwDebug.HexDumpWithAddrOffset(b, 0), hex.Dump(b))
	}
}
