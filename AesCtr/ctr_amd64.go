package AesCtr

import (
	"crypto/cipher"
	"sync"
	"unsafe"
)

func xorBytes(dst, a, b []byte) int

func fillEightBlocks(nr int, xk uintptr, dst uintptr, counter uintptr)

const streamBufferSize = 32 * BlockSize

type aesctr struct {
	block   *aesCipherAsm
	nr      int
	ctr     [BlockSize]byte
	buffer  []byte
	storage [streamBufferSize]byte
}

func (c *aesCipherAsm) NewCTR(iv []byte) cipher.Stream {
	if len(iv) != BlockSize {
		panic("cipher.NewCTR: IV length must equal block size")
	}
	var ac aesctr
	ac.block = c
	ac.nr = len(c.enc)/4 - 1
	copy(ac.ctr[:], iv)
	ac.buffer = ac.storage[:0]
	return &ac
}
func aesCtrRefill(c *aesctr) {

	c.buffer = c.storage[:streamBufferSize]
	encPtr := uintptr(unsafe.Pointer(&c.block.enc[0]))
	bufPtr := uintptr(unsafe.Pointer(&c.buffer[0]))
	ctrPtr := uintptr(unsafe.Pointer(&c.ctr[0]))

	fillEightBlocks(c.nr, encPtr, bufPtr, ctrPtr)
	bufPtr += 128
	fillEightBlocks(c.nr, encPtr, bufPtr, ctrPtr)
	bufPtr += 128
	fillEightBlocks(c.nr, encPtr, bufPtr, ctrPtr)
	bufPtr += 128
	fillEightBlocks(c.nr, encPtr, bufPtr, ctrPtr)
}

type ctrResetAble interface {
	ctrReset(block cipher.Block, iv []byte) cipher.Stream
}

func (ac *aesctr) ctrReset(block cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != BlockSize {
		return cipher.NewCTR(block, iv)
	}
	switch obj1 := block.(type) {
	case *aesCipherAsm:
		ac.block = obj1
	case *aesCipherGCM:
		ac.block = &obj1.aesCipherAsm
	default:
		return cipher.NewCTR(block, iv)
	}
	ac.nr = len(ac.block.enc)/4 - 1
	copy(ac.ctr[:], iv)
	ac.buffer = ac.storage[:0]
	return ac
}
func (c *aesctr) XORKeyStream(dst, src []byte) {
	if len(src) > 0 {

		_ = dst[len(src)-1]
	}
	for len(src) > 0 {
		if len(c.buffer) == 0 {
			aesCtrRefill(c)
		}
		n := xorBytes(dst, src, c.buffer)
		c.buffer = c.buffer[n:]
		src = src[n:]
		dst = dst[n:]
	}
}

func (c *aesCipherAsm) ctrOnce(req CtrForAesOnceXORKeyStreamRequest) {
	if len(req.TmpBuf) < 144 {
		req.TmpBuf = make([]byte, 144)
	}
	ctr := req.TmpBuf[:16]
	blockOutBuf := req.TmpBuf[16 : 16+128]
	copy(ctr[:], req.Iv)
	pos := 0
	nr := len(c.enc)/4 - 1
	for {
		if pos >= len(req.Src) {
			return
		}
		fillEightBlocks(nr, uintptr(unsafe.Pointer(&c.enc[0])), uintptr(unsafe.Pointer(&blockOutBuf[0])), uintptr(unsafe.Pointer(&ctr[0])))
		thisXorSize := len(req.Src) - pos
		if thisXorSize > 128 {
			thisXorSize = 128
		}
		xorBytes(req.Dst[pos:pos+thisXorSize], req.Src[pos:pos+thisXorSize], blockOutBuf[:thisXorSize])
		pos += thisXorSize
	}
}

var gAesCtrPool = sync.Pool{}

func PoolGetAesCtr(block cipher.Block, iv []byte) cipher.Stream {
	stream := gAesCtrPool.Get()
	if stream == nil {
		return cipher.NewCTR(block, iv)
	}
	obj1, ok := stream.(ctrResetAble)
	if ok {
		return obj1.ctrReset(block, iv)
	}
	return cipher.NewCTR(block, iv)
}

func PoolPutAesCtr(ctrObj cipher.Stream) {
	if ctrObj == nil {
		return
	}

	_, ok := ctrObj.(*aesctr)
	if !ok {
		return
	}
	gAesCtrPool.Put(ctrObj)
}
