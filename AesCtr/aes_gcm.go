// +build amd64

package AesCtr

import (
	"crypto/cipher"
	"crypto/subtle"
	"errors"
)

func hasGCMAsm() bool

func aesEncBlock(dst, src *[16]byte, ks []uint32)

func gcmAesInit(productTable *[256]byte, ks []uint32)

func gcmAesData(productTable *[256]byte, data []byte, T *[16]byte)

func gcmAesEnc(productTable *[256]byte, dst, src []byte, ctr, T *[16]byte, ks []uint32)

func gcmAesDec(productTable *[256]byte, dst, src []byte, ctr, T *[16]byte, ks []uint32)

func gcmAesFinish(productTable *[256]byte, tagMask, T *[16]byte, pLen, dLen uint64)

const (
	gcmBlockSize         = 16
	gcmTagSize           = 16
	gcmStandardNonceSize = 12
)

var errOpen = errors.New("cipher: message authentication failed")

type aesCipherGCM struct {
	aesCipherAsm
}

var _ gcmAble = (*aesCipherGCM)(nil)

func (c *aesCipherGCM) NewGCM(nonceSize int) (cipher.AEAD, error) {
	g := &gcmAsm{ks: c.enc, nonceSize: nonceSize}
	gcmAesInit(&g.productTable, g.ks)
	return g, nil
}

type gcmAsm struct {
	ks []uint32

	productTable [256]byte

	nonceSize int
}

func (g *gcmAsm) NonceSize() int {
	return g.nonceSize
}

func (*gcmAsm) Overhead() int {
	return gcmTagSize
}

func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func (g *gcmAsm) Seal(dst, nonce, plaintext, data []byte) []byte {
	if len(nonce) != g.nonceSize {
		panic("cipher: incorrect nonce length given to GCM")
	}
	if uint64(len(plaintext)) > ((1<<32)-2)*BlockSize {
		panic("cipher: message too large for GCM")
	}

	var counter, tagMask [gcmBlockSize]byte

	if len(nonce) == gcmStandardNonceSize {

		copy(counter[:], nonce)
		counter[gcmBlockSize-1] = 1
	} else {

		gcmAesData(&g.productTable, nonce, &counter)
		gcmAesFinish(&g.productTable, &tagMask, &counter, uint64(len(nonce)), uint64(0))
	}

	aesEncBlock(&tagMask, &counter, g.ks)

	var tagOut [gcmTagSize]byte
	gcmAesData(&g.productTable, data, &tagOut)

	ret, out := sliceForAppend(dst, len(plaintext)+gcmTagSize)
	if len(plaintext) > 0 {
		gcmAesEnc(&g.productTable, out, plaintext, &counter, &tagOut, g.ks)
	}
	gcmAesFinish(&g.productTable, &tagMask, &tagOut, uint64(len(plaintext)), uint64(len(data)))
	copy(out[len(plaintext):], tagOut[:])

	return ret
}

func (g *gcmAsm) Open(dst, nonce, ciphertext, data []byte) ([]byte, error) {
	if len(nonce) != g.nonceSize {
		panic("cipher: incorrect nonce length given to GCM")
	}

	if len(ciphertext) < gcmTagSize {
		return nil, errOpen
	}
	if uint64(len(ciphertext)) > ((1<<32)-2)*BlockSize+gcmTagSize {
		return nil, errOpen
	}

	tag := ciphertext[len(ciphertext)-gcmTagSize:]
	ciphertext = ciphertext[:len(ciphertext)-gcmTagSize]

	var counter, tagMask [gcmBlockSize]byte

	if len(nonce) == gcmStandardNonceSize {

		copy(counter[:], nonce)
		counter[gcmBlockSize-1] = 1
	} else {

		gcmAesData(&g.productTable, nonce, &counter)
		gcmAesFinish(&g.productTable, &tagMask, &counter, uint64(len(nonce)), uint64(0))
	}

	aesEncBlock(&tagMask, &counter, g.ks)

	var expectedTag [gcmTagSize]byte
	gcmAesData(&g.productTable, data, &expectedTag)

	ret, out := sliceForAppend(dst, len(ciphertext))
	if len(ciphertext) > 0 {
		gcmAesDec(&g.productTable, out, ciphertext, &counter, &expectedTag, g.ks)
	}
	gcmAesFinish(&g.productTable, &tagMask, &expectedTag, uint64(len(ciphertext)), uint64(len(data)))

	if subtle.ConstantTimeCompare(expectedTag[:], tag) != 1 {
		for i := range out {
			out[i] = 0
		}
		return nil, errOpen
	}

	return ret, nil
}
