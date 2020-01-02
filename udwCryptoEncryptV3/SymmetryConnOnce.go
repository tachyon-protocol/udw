package udwCryptoEncryptV3

import (
	"crypto/cipher"
	"crypto/rand"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
)

func MustSymmetryConnEncryptOnceWithBlock(src []byte, block cipher.Block, bufW *udwBytes.BufWriter) {
	buf := bufW.GetHeadBuffer(len(src) + 20 + 144)
	_, err := io.ReadFull(rand.Reader, buf[:16])
	if err != nil {
		panic(err)
	}
	copy(buf[16:20], gMagicBuf)
	copy(buf[20:], src)
	AesCtr.CtrForAesOnceXORKeyStream(AesCtr.CtrForAesOnceXORKeyStreamRequest{
		AesBlock: block,
		Dst:      buf[16:],
		Src:      buf[16:],
		Iv:       buf[:16],
		TmpBuf:   buf[len(src)+20:],
	})
	bufW.AddPos(len(src) + 20)
	return
}
