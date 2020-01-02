package udwCryptoEncryptV3

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/tachyon-protocol/udw/udwBase64"
	"github.com/tachyon-protocol/udw/udwRand"
)

var magicCode4 = [4]byte{0xa7, 0x97, 0x6d, 0x15}

func EncryptV3(key *[32]byte, data []byte) (output []byte) {
	Iv := udwRand.MustCryptoRandBytes(16)
	block, err := aes.NewCipher((*key)[:])
	if err != nil {
		panic(err)
	}
	afterCbcSize := len(data)
	output = make([]byte, afterCbcSize+16+4)
	copy(output[:16], Iv)
	copy(output[16:len(data)+16], data)
	copy(output[len(data)+16:len(data)+16+4], magicCode4[:])
	ctr := cipher.NewCTR(block, Iv)
	ctr.XORKeyStream(output[16:], output[16:])
	return output
}

func DecryptV3(key *[32]byte, data []byte) (output []byte, err error) {

	if len(data) < 20 {
		return nil, errors.New("[udwCrypto.DecryptV3] input data too small")
	}
	aseKey := key[:]
	Iv := data[:16]
	block, err := aes.NewCipher(aseKey)
	if err != nil {
		return nil, err
	}
	output = make([]byte, len(data)-16)
	ctr := cipher.NewCTR(block, Iv)
	ctr.XORKeyStream(output, data[16:])
	beforeCbcSize := len(data) - 16 - 4
	if !bytes.Equal(magicCode4[:], output[beforeCbcSize:beforeCbcSize+4]) {
		return nil, errors.New("[udwCrypto.DecryptV3] magicCode not match mtjqaqzgm9")
	}
	output = output[:beforeCbcSize]

	return output, nil
}

func EncryptV3WithBase64(key *[32]byte, input []byte) (output string) {
	outputB := EncryptV3(key, input)
	output = udwBase64.EncodeByteToStringV2(outputB)
	return output
}

func DecryptV3WithBase64(key *[32]byte, s string) (output []byte, err error) {
	b1, err := udwBase64.DecodeStringToByteV2(s)
	if err != nil {
		return nil, err
	}
	return DecryptV3(key, b1)
}
