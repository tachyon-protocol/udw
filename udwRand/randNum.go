package udwRand

import (
	"crypto/rand"
	"encoding/binary"
)

func MustCryptoRandUint64NotZero() uint64 {
	return CryptoRandUint64NotZero()
}

func MustCryptoRandUint64() uint64 {
	var tmpBuf [8]byte
	_, err := rand.Read(tmpBuf[:])
	if err != nil {
		panic(err)
	}
	ret := binary.LittleEndian.Uint64(tmpBuf[:])
	return ret
}

func CryptoRandUint64NotZero() uint64 {
	for i := 0; i < 100; i++ {
		ret := MustCryptoRandUint64()
		if ret != 0 {
			return ret
		}
	}
	panic("[MustCryptoRandUint64NotZero] too many times loop,may be some bug.")
}

func MustCryptoRandUint16() uint16 {
	var tmpBuf [2]byte
	_, err := rand.Read(tmpBuf[:])
	if err != nil {
		panic(err)
	}
	ret := binary.LittleEndian.Uint16(tmpBuf[:])
	return ret
}
