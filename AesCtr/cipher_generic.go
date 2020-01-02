// +build !amd64,!s390x,!ppc64le

package AesCtr

import (
	"crypto/cipher"
)

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherGeneric(key)
}

func expandKey(key []byte, enc, dec []uint32) {
	expandKeyGo(key, enc, dec)
}
