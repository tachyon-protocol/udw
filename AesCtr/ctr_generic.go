// +build !amd64

package AesCtr

import (
	"crypto/cipher"
)

func PoolGetAesCtr(block cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCTR(block, iv)
}

func PoolPutAesCtr(ctrObj cipher.Stream) {
}
