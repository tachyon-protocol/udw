package AesCtr

import (
	"crypto/cipher"
	"crypto/subtle"
	"errors"
)

type gcmCount [16]byte

func (x *gcmCount) inc() {

	n := uint32(x[15]) | uint32(x[14])<<8 | uint32(x[13])<<16 | uint32(x[12])<<24
	n += 1
	x[12] = byte(n >> 24)
	x[13] = byte(n >> 16)
	x[14] = byte(n >> 8)
	x[15] = byte(n)
}

func gcmLengths(len0, len1 uint64) [16]byte {
	return [16]byte{
		byte(len0 >> 56),
		byte(len0 >> 48),
		byte(len0 >> 40),
		byte(len0 >> 32),
		byte(len0 >> 24),
		byte(len0 >> 16),
		byte(len0 >> 8),
		byte(len0),
		byte(len1 >> 56),
		byte(len1 >> 48),
		byte(len1 >> 40),
		byte(len1 >> 32),
		byte(len1 >> 24),
		byte(len1 >> 16),
		byte(len1 >> 8),
		byte(len1),
	}
}

type gcmHashKey [16]byte

type gcmAsm struct {
	block     *aesCipherAsm
	hashKey   gcmHashKey
	nonceSize int
}

const (
	gcmBlockSize         = 16
	gcmTagSize           = 16
	gcmStandardNonceSize = 12
)

var errOpen = errors.New("cipher: message authentication failed")

var _ gcmAble = (*aesCipherAsm)(nil)

func (c *aesCipherAsm) NewGCM(nonceSize int) (cipher.AEAD, error) {
	var hk gcmHashKey
	c.Encrypt(hk[:], hk[:])
	g := &gcmAsm{
		block:     c,
		hashKey:   hk,
		nonceSize: nonceSize,
	}
	return g, nil
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

func ghash(key *gcmHashKey, hash *[16]byte, data []byte)

func (g *gcmAsm) paddedGHASH(hash *[16]byte, data []byte) {
	siz := len(data) &^ 0xf
	if siz > 0 {
		ghash(&g.hashKey, hash, data[:siz])
		data = data[siz:]
	}
	if len(data) > 0 {
		var s [16]byte
		copy(s[:], data)
		ghash(&g.hashKey, hash, s[:])
	}
}

func cryptBlocksGCM(fn code, key, dst, src, buf []byte, cnt *gcmCount)

func (g *gcmAsm) counterCrypt(dst, src []byte, cnt *gcmCount) {

	var ctrbuf, srcbuf [2048]byte
	for len(src) >= 16 {
		siz := len(src)
		if len(src) > len(ctrbuf) {
			siz = len(ctrbuf)
		}
		siz &^= 0xf
		copy(srcbuf[:], src[:siz])
		cryptBlocksGCM(g.block.function, g.block.key, dst[:siz], srcbuf[:siz], ctrbuf[:], cnt)
		src = src[siz:]
		dst = dst[siz:]
	}
	if len(src) > 0 {
		var x [16]byte
		g.block.Encrypt(x[:], cnt[:])
		for i := range src {
			dst[i] = src[i] ^ x[i]
		}
		cnt.inc()
	}
}

func (g *gcmAsm) deriveCounter(nonce []byte) gcmCount {

	var counter gcmCount
	if len(nonce) == gcmStandardNonceSize {
		copy(counter[:], nonce)
		counter[gcmBlockSize-1] = 1
	} else {
		var hash [16]byte
		g.paddedGHASH(&hash, nonce)
		lens := gcmLengths(0, uint64(len(nonce))*8)
		g.paddedGHASH(&hash, lens[:])
		copy(counter[:], hash[:])
	}
	return counter
}

func (g *gcmAsm) auth(out, ciphertext, additionalData []byte, tagMask *[gcmTagSize]byte) {
	var hash [16]byte
	g.paddedGHASH(&hash, additionalData)
	g.paddedGHASH(&hash, ciphertext)
	lens := gcmLengths(uint64(len(additionalData))*8, uint64(len(ciphertext))*8)
	g.paddedGHASH(&hash, lens[:])

	copy(out, hash[:])
	for i := range out {
		out[i] ^= tagMask[i]
	}
}

func (g *gcmAsm) Seal(dst, nonce, plaintext, data []byte) []byte {
	if len(nonce) != g.nonceSize {
		panic("cipher: incorrect nonce length given to GCM")
	}
	if uint64(len(plaintext)) > ((1<<32)-2)*BlockSize {
		panic("cipher: message too large for GCM")
	}

	ret, out := sliceForAppend(dst, len(plaintext)+gcmTagSize)

	counter := g.deriveCounter(nonce)

	var tagMask [gcmBlockSize]byte
	g.block.Encrypt(tagMask[:], counter[:])
	counter.inc()

	g.counterCrypt(out, plaintext, &counter)
	g.auth(out[len(plaintext):], out[:len(plaintext)], data, &tagMask)

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

	counter := g.deriveCounter(nonce)

	var tagMask [gcmBlockSize]byte
	g.block.Encrypt(tagMask[:], counter[:])
	counter.inc()

	var expectedTag [gcmTagSize]byte
	g.auth(expectedTag[:], ciphertext, data, &tagMask)

	ret, out := sliceForAppend(dst, len(ciphertext))

	if subtle.ConstantTimeCompare(expectedTag[:], tag) != 1 {

		for i := range out {
			out[i] = 0
		}
		return nil, errOpen
	}

	g.counterCrypt(out, ciphertext, &counter)
	return ret, nil
}
