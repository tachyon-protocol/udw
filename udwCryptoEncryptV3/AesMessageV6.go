package udwCryptoEncryptV3

import (
	"bytes"
	"crypto/cipher"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwRand/udwRandFast"
)

type AesMessageV6KeyObject struct {
	block cipher.Block
}

func NewAesMessageV6KeyObject(key *[32]byte) AesMessageV6KeyObject {
	return AesMessageV6KeyObject{
		block: GetAesBlockFrom32Psk(key),
	}
}

type AesMessageV6Request struct {
	KeyObject  AesMessageV6KeyObject
	Data       []byte
	OutputBufW *udwBytes.BufWriter
}

func AesMessageV6Encrypt(req AesMessageV6Request) {
	afterCtrSize := len(req.Data)
	bufSize := afterCtrSize + 16 + 4 + 144
	if req.OutputBufW == nil {
		panic("[AesMessageV6Encrypt] req.OutputBufW==nil")
	}
	req.OutputBufW.Reset()
	buf := req.OutputBufW.GetHeadBuffer(bufSize)
	udwRandFast.MustRead(buf[:16])

	copy(buf[16:len(req.Data)+16], req.Data)
	copy(buf[len(req.Data)+16:len(req.Data)+16+4], magicCode4[:])

	AesCtr.CtrForAesOnceXORKeyStream(AesCtr.CtrForAesOnceXORKeyStreamRequest{
		AesBlock: req.KeyObject.block,
		Dst:      buf[16:],
		Src:      buf[16:],
		Iv:       buf[:16],
		TmpBuf:   buf[len(req.Data)+16+4:],
	})
	req.OutputBufW.AddPos(afterCtrSize + 16 + 4)
}

func AesMessageV6Decrypt(req AesMessageV6Request) (errMsg string) {
	if len(req.Data) < 20 {
		return "[AesMessageV6Decrypt] input data too small"
	}
	req.OutputBufW.Reset()
	buf := req.OutputBufW.GetHeadBuffer(len(req.Data) - 16 + 144)
	AesCtr.CtrForAesOnceXORKeyStream(AesCtr.CtrForAesOnceXORKeyStreamRequest{
		AesBlock: req.KeyObject.block,
		Dst:      buf[:len(req.Data)-16],
		Src:      req.Data[16:],
		Iv:       req.Data[:16],
		TmpBuf:   buf[len(req.Data)-16:],
	})
	beforeCbcSize := len(req.Data) - 16 - 4
	if !bytes.Equal(magicCode4[:], buf[beforeCbcSize:beforeCbcSize+4]) {
		return "[AesMessageV6Decrypt] magicCode not match mtjqaqzgm9"
	}
	req.OutputBufW.AddPos(beforeCbcSize)
	return ""
}
