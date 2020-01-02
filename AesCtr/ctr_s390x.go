package AesCtr

import (
	"crypto/cipher"
	"unsafe"
)

var _ ctrAble = (*aesCipherAsm)(nil)

func xorBytes(dst, a, b []byte) int

const streamBufferSize = 32 * BlockSize

type aesctr struct {
	block   *aesCipherAsm
	ctr     [2]uint64
	buffer  []byte
	storage [streamBufferSize]byte
}

func (c *aesCipherAsm) NewCTR(iv []byte) cipher.Stream {
	if len(iv) != BlockSize {
		panic("cipher.NewCTR: IV length must equal block size")
	}
	var ac aesctr
	ac.block = c
	ac.ctr[0] = *(*uint64)(unsafe.Pointer((&iv[0])))
	ac.ctr[1] = *(*uint64)(unsafe.Pointer((&iv[8])))
	ac.buffer = ac.storage[:0]
	return &ac
}

func (c *aesctr) refill() {

	c.buffer = c.storage[:streamBufferSize]
	c0, c1 := c.ctr[0], c.ctr[1]
	for i := 0; i < streamBufferSize; i += BlockSize {
		b0 := (*uint64)(unsafe.Pointer(&c.buffer[i]))
		b1 := (*uint64)(unsafe.Pointer(&c.buffer[i+BlockSize/2]))
		*b0, *b1 = c0, c1

		c1++
		if c1 == 0 {

			c0++
		}
	}
	c.ctr[0], c.ctr[1] = c0, c1

	cryptBlocks(c.block.function, &c.block.key[0], &c.buffer[0], &c.buffer[0], streamBufferSize)
}

func (c *aesctr) XORKeyStream(dst, src []byte) {
	for len(src) > 0 {
		if len(c.buffer) == 0 {
			c.refill()
		}
		n := xorBytes(dst, src, c.buffer)
		c.buffer = c.buffer[n:]
		src = src[n:]
		dst = dst[n:]
	}
}
