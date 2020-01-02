package udwCryptoEncryptV3

import (
	"crypto/cipher"
)

type ctrForAes struct {
	b       cipher.Block
	ctr     [16]byte
	outBuf  [512]byte
	out     []byte
	outUsed int
}

func (x *ctrForAes) refill(srcSize int) {
	remain := len(x.out) - x.outUsed
	copy(x.out[:], x.out[x.outUsed:])
	x.out = x.out[:cap(x.out)]
	bs := x.b.BlockSize()
	targetSize := len(x.out) - bs
	if targetSize > srcSize+bs {
		targetSize = srcSize + bs
	}
	for remain <= targetSize {
		x.b.Encrypt(x.out[remain:], x.ctr[:])
		remain += bs

		for i := 16 - 1; i >= 0; i-- {
			x.ctr[i]++
			if x.ctr[i] != 0 {
				break
			}
		}
	}
	x.out = x.out[:remain]
	x.outUsed = 0
}

func (x *ctrForAes) XORKeyStream(dst, src []byte, iv []byte) {
	copy(x.ctr[:], iv)
	x.outUsed = 0
	x.out = x.out[0:0]
	for len(src) > 0 {
		if x.outUsed >= len(x.out)-x.b.BlockSize() {
			x.refill(len(src))
		}
		n := XorBytes(dst, src, x.out[x.outUsed:])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

type ctrForAesXORKeyStreamRequest struct {
	AesBlock cipher.Block
	Dst      []byte
	Src      []byte
	Iv       []byte

	tmpBuf []byte
}

func ctrForAesXORKeyStream(req ctrForAesXORKeyStreamRequest) {
	const blockSize = 16
	ctr := req.tmpBuf[:16]
	blockOutBuf := req.tmpBuf[16:32]
	copy(ctr[:], req.Iv)
	pos := 0
	for {
		if pos >= len(req.Src) {
			return
		}
		req.AesBlock.Encrypt(blockOutBuf[:], ctr[:])
		for i := 16 - 1; i >= 0; i-- {
			ctr[i]++
			if ctr[i] != 0 {
				break
			}
		}
		thisXorSize := len(req.Src) - pos
		if thisXorSize > blockSize {
			thisXorSize = blockSize
		}

		XorBytes(req.Dst[pos:], req.Src[pos:], blockOutBuf[:thisXorSize])
		pos += blockSize
	}
}
