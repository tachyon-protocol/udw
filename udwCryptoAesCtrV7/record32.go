package udwCryptoAesCtrV7

import (
	"encoding/base64"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
)

func Encrypt32(psk *[32]byte, in []byte) (out []byte) {
	return EncryptRecord(RecordRequest{
		In:             in,
		Psk:            psk,
		RandomOverHead: 16,
		VerifyOverHead: 16,
	})
}

func Decrypt32(psk *[32]byte, in []byte) (out []byte, errMsg string) {
	return DecryptRecord(RecordRequest{
		In:             in,
		Psk:            psk,
		RandomOverHead: 16,
		VerifyOverHead: 16,
	})
}

func Encrypt32ToBase64String(pskS string, in []byte) (outB64 string) {
	psk := Get32PskSha3FromString(pskS)
	out := Encrypt32(psk, in)
	return base64.RawURLEncoding.EncodeToString(out)
}

func Decrypt32FromBase64String(pskS string, inB64 string) (out []byte, ok bool) {
	in, err := base64.RawURLEncoding.DecodeString(inB64)
	if err != nil {
		return nil, false
	}
	psk := Get32PskSha3FromString(pskS)
	out, errMsg := Decrypt32(psk, in)
	if errMsg != "" {
		return nil, false
	}
	return out, true
}

func Get32PskSha3FromString(s string) *[32]byte {
	psk := udwCryptoSha3.Sum512([]byte(s))
	pskOut := [32]byte{}
	copy(pskOut[:], psk[:32])
	return &pskOut
}
