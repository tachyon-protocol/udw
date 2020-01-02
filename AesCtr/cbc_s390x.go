package AesCtr

import (
	"crypto/cipher"
)

var _ cbcEncAble = (*aesCipherAsm)(nil)
var _ cbcDecAble = (*aesCipherAsm)(nil)

type cbc struct {
	b  *aesCipherAsm
	c  code
	iv [BlockSize]byte
}

func (b *aesCipherAsm) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	var c cbc
	c.b = b
	c.c = b.function
	copy(c.iv[:], iv)
	return &c
}

func (b *aesCipherAsm) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	var c cbc
	c.b = b
	c.c = b.function + 128
	copy(c.iv[:], iv)
	return &c
}

func (x *cbc) BlockSize() int { return BlockSize }

func cryptBlocksChain(c code, iv, key, dst, src *byte, length int)

func (x *cbc) CryptBlocks(dst, src []byte) {
	if len(src)%BlockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	if len(src) > 0 {
		cryptBlocksChain(x.c, &x.iv[0], &x.b.key[0], &dst[0], &src[0], len(src))
	}
}

func (x *cbc) SetIV(iv []byte) {
	if len(iv) != BlockSize {
		panic("cipher: incorrect length IV")
	}
	copy(x.iv[:], iv)
}
