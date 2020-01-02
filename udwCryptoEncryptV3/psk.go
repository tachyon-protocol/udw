package udwCryptoEncryptV3

import (
	"crypto/cipher"
	"crypto/sha512"
	"github.com/tachyon-protocol/udw/AesCtr"
	"strconv"
)

func Get32PskFromString(s string) *[32]byte {
	psk := sha512.Sum512([]byte(s))
	pskOut := [32]byte{}
	copy(pskOut[:], psk[:32])
	return &pskOut
}

func Get32PskFromSlice32(b []byte) *[32]byte {
	if len(b) != 32 {
		panic("2g6zve8mt5 " + strconv.Itoa(len(b)))
	}
	pskOut := [32]byte{}
	copy(pskOut[:], b[:32])
	return &pskOut
}

func GetAesBlockFrom32Psk(key *[32]byte) cipher.Block {
	key2 := *key
	block, err := AesCtr.NewCipher(key2[:])
	if err != nil {
		panic(err)
	}
	return block
}
