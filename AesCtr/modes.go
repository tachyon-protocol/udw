package AesCtr

import (
	"crypto/cipher"
)

type gcmAble interface {
	NewGCM(size int) (cipher.AEAD, error)
}

type cbcEncAble interface {
	NewCBCEncrypter(iv []byte) cipher.BlockMode
}

type cbcDecAble interface {
	NewCBCDecrypter(iv []byte) cipher.BlockMode
}

type ctrAble interface {
	NewCTR(iv []byte) cipher.Stream
}
