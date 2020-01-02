package udwCryptoAesCtrV7

import (
	"github.com/tachyon-protocol/udw/udwCryptoTester"
	"testing"
)

func TestRecord16(ot *testing.T) {
	udwCryptoTester.EncryptTesterErrMsg(udwCryptoTester.EncryptTesterErrMsgRequest{
		Encrypt:     Encrypt16NoRandom,
		Decrypt:     Decrypt16NoRandom,
		MaxOverhead: 32,
		NoRandom:    true,
	})
	udwCryptoTester.EncryptTesterErrMsg(udwCryptoTester.EncryptTesterErrMsgRequest{
		Encrypt:     Encrypt32,
		Decrypt:     Decrypt32,
		MaxOverhead: 32,
	})
	type testCase struct {
		RandomOverHead  int
		VerifyOverHead  int
		NoRandom        bool
		RandomCheckSize int
	}
	for _, req := range []testCase{
		{
			RandomOverHead: 16,
			VerifyOverHead: 16,
		},
		{
			RandomOverHead:  2,
			VerifyOverHead:  6,
			RandomCheckSize: 10,
		},
		{
			RandomOverHead: 4,
			VerifyOverHead: 12,
		},
		{
			RandomOverHead: 128,
			VerifyOverHead: 64,
		},
		{
			RandomOverHead: 0,
			VerifyOverHead: 8,
			NoRandom:       true,
		},
	} {
		udwCryptoTester.EncryptTesterErrMsg(udwCryptoTester.EncryptTesterErrMsgRequest{
			Encrypt: func(key *[32]byte, data []byte) (output []byte) {
				return EncryptRecord(RecordRequest{
					In:             data,
					Psk:            key,
					RandomOverHead: req.RandomOverHead,
					VerifyOverHead: req.VerifyOverHead,
				})
			},
			Decrypt: func(key *[32]byte, data []byte) (output []byte, errMsg string) {
				return DecryptRecord(RecordRequest{
					In:             data,
					Psk:            key,
					RandomOverHead: req.RandomOverHead,
					VerifyOverHead: req.VerifyOverHead,
				})
			},
			MaxOverhead:     req.RandomOverHead + req.VerifyOverHead,
			NoRandom:        req.NoRandom,
			RandomCheckSize: req.RandomCheckSize,
		})
	}
}
