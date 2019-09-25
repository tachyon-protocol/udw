package udwCryptoSha3

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwRand"
	"math/rand"
)

func AlphaNumByString(inString string, returnLen int) string {
	if returnLen == 0 {
		return ""
	}
	hash := make([]byte, returnLen*2)
	ShakeSum256(hash, []byte(inString))
	return udwRand.EncodeReadableAlphaNumForRand(hash)
}

func GetSha3ShakeSum256RandomGenerater(inString string) func(returnLen int) []byte {
	return func(returnLen int) []byte {
		if returnLen == 0 {
			return nil
		}
		hash := make([]byte, returnLen)
		ShakeSum256(hash, []byte(inString))
		return hash
	}
}

func NewGoRandByString(inString string) *rand.Rand {
	hasher := NewShake256()
	hasher.Write([]byte(inString))
	rander := rand.New(sha3RandSource{
		hasher: hasher,
	})
	return rander
}

type sha3RandSource struct {
	hasher ShakeHash
}

func (source sha3RandSource) Int63() int64 {
	buf := make([]byte, 8)
	source.hasher.Read(buf)
	out := binary.LittleEndian.Uint64(buf)
	out = out >> 1
	return int64(out)
}

func (source sha3RandSource) Seed(seed int64) {

}

func IntnByString(inString string, n int) int {
	rander := NewGoRandByString(inString)
	return rander.Intn(n)
}

func Float64ByString(inString string) float64 {
	rander := NewGoRandByString(inString)
	return rander.Float64()
}
