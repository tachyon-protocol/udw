package udwCryptoEncryptV3

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/tachyon-protocol/udw/udwRand"
	"sync"
)

func NewUdwEncryterV3NoAlloc(key *[32]byte) *EncryptV3Obj {
	block, err := aes.NewCipher((*key)[:])
	if err != nil {
		panic(err)
	}
	obj := &EncryptV3Obj{}
	obj.ctr.b = block
	obj.ctr.out = obj.ctr.outBuf[:]
	return obj
}

func NewUdwEncryterV3NoAllocV2(aesBlock cipher.Block) *EncryptV3Obj {
	obj := &EncryptV3Obj{}
	obj.ctr.b = aesBlock
	obj.ctr.out = obj.ctr.outBuf[:]
	return obj
}

type EncryptV3Obj struct {
	ctr    ctrForAes
	locker sync.Mutex
}

func (obj *EncryptV3Obj) GetAesBlock() cipher.Block {
	return obj.ctr.b
}

func (obj *EncryptV3Obj) EncryptWithBuf(data []byte, buf []byte) []byte {
	obj.locker.Lock()
	afterCbcSize := len(data)
	bufSize := afterCbcSize + 16 + 4
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	} else {
		buf = buf[:bufSize]
	}
	udwRand.MustCryptoRandBytesWithBuf(buf[:16])
	copy(buf[16:len(data)+16], data)
	copy(buf[len(data)+16:len(data)+16+4], magicCode4[:])
	obj.ctr.XORKeyStream(buf[16:], buf[16:], buf[:16])
	obj.locker.Unlock()
	return buf
}

func (obj *EncryptV3Obj) DecryptWithBuf(data []byte, buf []byte) (output []byte, err error) {
	obj.locker.Lock()

	if len(data) < 20 {
		obj.locker.Unlock()
		return nil, errors.New("[udwCrypto.DecryptV3] input data too small")
	}
	Iv := data[:16]
	bufSize := len(data) - 16
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	} else {
		buf = buf[:bufSize]
	}
	obj.ctr.XORKeyStream(buf, data[16:], Iv)
	beforeCbcSize := len(data) - 16 - 4
	if !bytes.Equal(magicCode4[:], buf[beforeCbcSize:beforeCbcSize+4]) {
		obj.locker.Unlock()
		return nil, errors.New("[udwCrypto.DecryptV3] magicCode not match ta7jwe3swh")
	}
	buf = buf[:beforeCbcSize]

	obj.locker.Unlock()
	return buf, nil
}

func (obj *EncryptV3Obj) Clone() *EncryptV3Obj {
	returnObj := &EncryptV3Obj{}
	returnObj.ctr.b = obj.ctr.b
	returnObj.ctr.out = returnObj.ctr.outBuf[:]
	return returnObj
}
