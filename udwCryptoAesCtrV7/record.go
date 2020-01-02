package udwCryptoAesCtrV7

import (
	"crypto/hmac"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"github.com/tachyon-protocol/udw/udwRand"
)

type RecordRequest struct {
	In  []byte
	Psk *[32]byte

	RandomOverHead int
	VerifyOverHead int
}

func EncryptRecord(req RecordRequest) (out []byte) {
	if req.RandomOverHead < 0 {
		panic("fhvkukae9p")
	}
	if req.VerifyOverHead <= 0 {
		panic("5nvfnvpgye")
	}
	if req.VerifyOverHead > 64 {
		panic("pew57g3fdd")
	}
	in2 := append(req.In, udwRand.MustCryptoRandBytes(req.RandomOverHead)...)
	mac := hmac.New(udwCryptoSha3.New512, (*req.Psk)[:])
	mac.Write(in2)
	headerAll := mac.Sum(nil)
	iv := udwCryptoSha3.Sum512(headerAll[:req.VerifyOverHead])
	block, err := AesCtr.NewCipher((*req.Psk)[:])
	if err != nil {
		panic(err)
	}
	out = make([]byte, len(req.In)+req.RandomOverHead+req.VerifyOverHead)
	copy(out[:req.VerifyOverHead], headerAll[:req.VerifyOverHead])
	AesCtr.CtrForAesOnceXORKeyStream(AesCtr.CtrForAesOnceXORKeyStreamRequest{
		AesBlock: block,
		Dst:      out[req.VerifyOverHead:],
		Src:      in2,
		Iv:       iv[:16],
	})
	return out
}

func DecryptRecord(req RecordRequest) (out []byte, errMsg string) {
	if len(req.In) < req.VerifyOverHead+req.RandomOverHead {
		return nil, "ygt885j2dr"
	}
	iv1 := req.In[:req.VerifyOverHead]
	iv2 := udwCryptoSha3.Sum512(iv1)
	block, err := AesCtr.NewCipher((*req.Psk)[:])
	if err != nil {
		panic(err)
	}
	out = make([]byte, len(req.In)-req.VerifyOverHead)
	AesCtr.CtrForAesOnceXORKeyStream(AesCtr.CtrForAesOnceXORKeyStreamRequest{
		AesBlock: block,
		Dst:      out,
		Src:      req.In[req.VerifyOverHead:],
		Iv:       iv2[:16],
	})
	mac := hmac.New(udwCryptoSha3.New512, (*req.Psk)[:])
	mac.Write(out)
	headerAll := mac.Sum(nil)
	if hmac.Equal(headerAll[:req.VerifyOverHead], iv1) == false {
		return nil, "888esbt26t"
	}
	return out[:len(out)-req.RandomOverHead], ""
}
