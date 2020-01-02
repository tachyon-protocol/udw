package AesCtr

import (
	"crypto/cipher"
	"github.com/tachyon-protocol/udw/AesCtr/internal/cipherhw"
)

type code int

const (
	aes128 code = 18
	aes192      = 19
	aes256      = 20
)

type aesCipherAsm struct {
	function code
	key      []byte
	storage  [256]byte
}

func cryptBlocks(c code, key, dst, src *byte, length int)

var useAsm = cipherhw.AESGCMSupport()

func newCipher(key []byte) (cipher.Block, error) {
	if !useAsm {
		return newCipherGeneric(key)
	}

	var function code
	switch len(key) {
	case 128 / 8:
		function = aes128
	case 192 / 8:
		function = aes192
	case 256 / 8:
		function = aes256
	default:
		return nil, KeySizeError(len(key))
	}

	var c aesCipherAsm
	c.function = function
	c.key = c.storage[:len(key)]
	copy(c.key, key)
	return &c, nil
}

func (c *aesCipherAsm) BlockSize() int { return BlockSize }

func (c *aesCipherAsm) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < BlockSize {
		panic("crypto/aes: output not full block")
	}
	cryptBlocks(c.function, &c.key[0], &dst[0], &src[0], BlockSize)
}

func (c *aesCipherAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < BlockSize {
		panic("crypto/aes: output not full block")
	}

	cryptBlocks(c.function+128, &c.key[0], &dst[0], &src[0], BlockSize)
}

func expandKey(key []byte, enc, dec []uint32) {
	expandKeyGo(key, enc, dec)
}
