package udwCryptoSha3

import (
	"encoding/hex"
	"io"
	"os"
)

func Sha3512ToHexString(content []byte) string {
	outBuf := Sum512(content)
	return hex.EncodeToString(outBuf[:])
}

func Sha3512ToHexString32(content []byte) string {
	v := Sha3512ToHexString(content)
	return v[:32]
}

func Sha3512First16(content []byte) []byte {
	outBuf := Sum512(content)
	return outBuf[:16]
}

func Sha3512ToHexStringFromString(content string) string {
	return Sha3512ToHexString([]byte(content))
}

func Sha3512File(path string) (string, error) {
	hash := New512()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func MustSha3512File(path string) (b []byte) {
	hash := New512()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(hash, f)
	if err != nil {
		panic(err)
	}
	return hash.Sum(nil)
}

func Sum512Slice(data []byte) (digest []byte) {
	h := New512()
	h.Write(data)
	return h.Sum(nil)
}
