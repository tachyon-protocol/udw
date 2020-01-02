package udwCryptoAesCtrV7

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoEncryptV3"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestDecryptRecord(t *testing.T) {
	psk := udwCryptoEncryptV3.Get32PskFromString("123")
	testFn := func(in []byte) {
		encrypted := EncryptRecord(RecordRequest{
			In:             in,
			Psk:            psk,
			RandomOverHead: 2,
			VerifyOverHead: 6,
		})
		udwTest.Equal(len(encrypted), len(in)+8)
		out, errMsg := DecryptRecord(RecordRequest{
			In:             encrypted,
			Psk:            psk,
			RandomOverHead: 2,
			VerifyOverHead: 6,
		})
		udwTest.Equal(errMsg, "")
		udwTest.Equal(out, in)
		for i := 0; i < len(encrypted); i++ {
			e2 := udwBytes.Clone(encrypted)
			e2[i] = e2[i] + 1
			_, errMsg := DecryptRecord(RecordRequest{
				In:             e2,
				Psk:            psk,
				RandomOverHead: 2,
				VerifyOverHead: 6,
			})
			udwTest.Ok(errMsg != "")
		}
	}
	testFn([]byte{1})
	for _, size := range []int{
		10, 100, 1000,
	} {
		testFn(udwRand.MustCryptoRandBytes(size))
	}
}
