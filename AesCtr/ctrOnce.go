package AesCtr

import (
	"crypto/cipher"
)

type CtrForAesOnceXORKeyStreamRequest struct {
	AesBlock cipher.Block
	Dst      []byte
	Src      []byte
	Iv       []byte

	TmpBuf []byte
}

type ctrOnceInterface interface {
	ctrOnce(req CtrForAesOnceXORKeyStreamRequest)
}

func CtrForAesOnceXORKeyStream(req CtrForAesOnceXORKeyStreamRequest) {
	ctrOnce, ok := req.AesBlock.(ctrOnceInterface)
	if ok {
		ctrOnce.ctrOnce(req)
		return
	}
	if len(req.TmpBuf) < 32 {
		req.TmpBuf = make([]byte, 32)
	}
	const blockSize = 16
	ctr := req.TmpBuf[:16]
	blockOutBuf := req.TmpBuf[16:32]
	copy(ctr[:], req.Iv)
	pos := 0
	for {
		if pos >= len(req.Src) {
			return
		}
		req.AesBlock.Encrypt(blockOutBuf[:], ctr[:])
		for i := blockSize - 1; i >= 0; i-- {
			ctr[i]++
			if ctr[i] != 0 {
				break
			}
		}
		thisXorSize := len(req.Src) - pos
		if thisXorSize > blockSize {
			thisXorSize = blockSize
		}
		xorBytes(req.Dst[pos:pos+thisXorSize], req.Src[pos:pos+thisXorSize], blockOutBuf[:thisXorSize])

		pos += blockSize
	}
}
