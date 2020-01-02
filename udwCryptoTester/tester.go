package udwCryptoTester

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoEncryptV3"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
)

type EncryptTesterErrMsgRequest struct {
	Encrypt     func(key *[32]byte, data []byte) (output []byte)
	Decrypt     func(key *[32]byte, data []byte) (output []byte, errMsg string)
	DecryptErr  func(key *[32]byte, data []byte) (output []byte, err error)
	MaxOverhead int

	RandomCheckSize int
	NoRandom        bool
	NoCorrectVerify bool
}

func EncryptTesterErrMsg(req EncryptTesterErrMsgRequest) {
	if req.DecryptErr != nil && req.Decrypt == nil {
		req.Decrypt = func(key *[32]byte, data []byte) (output []byte, errMsg string) {
			output, err := req.DecryptErr(key, data)
			return output, udwErr.ErrorToMsg(err)
		}
	}
	key := &[32]byte{0xd8, 0x51, 0xea, 0x81, 0xb9, 0xe, 0xf, 0x2f, 0x8c, 0x85, 0x5f, 0xb6, 0x14, 0xb2}
	encryptTesterL2(encryptTesterL2Request{
		Encrypt: func(data []byte) (output []byte) {
			return req.Encrypt(key, data)
		},
		Decrypt: func(data []byte) (output []byte, errMsg string) {
			return req.Decrypt(key, data)
		},
		InReq: req,
	})
	if req.NoCorrectVerify == false {
		origin := []byte("jrez7dtu95yqxsq6gn8uxncqudcxc63cgetvp8h4684sekcxeb5jhtv4xzxpj8mj")
		k1 := udwCryptoEncryptV3.Get32PskFromString("123")
		k2 := udwCryptoEncryptV3.Get32PskFromString("456")
		ob := req.Encrypt(k1, origin)
		_, errMsg := req.Decrypt(k2, ob)
		udwTest.Ok(errMsg != "")
	}
}

type EncryptTesterErrMsgNoKeyRequest struct {
	Encrypt     func(data []byte) (output []byte)
	Decrypt     func(data []byte) (output []byte, errMsg string)
	MaxOverhead int
}

type encryptTesterL2Request struct {
	Encrypt func(data []byte) (output []byte)
	Decrypt func(data []byte) (output []byte, errMsg string)
	InReq   EncryptTesterErrMsgRequest
}

func encryptTesterL2(req encryptTesterL2Request) {
	encrypt := req.Encrypt
	decrypt := req.Decrypt
	maxOverhead := req.InReq.MaxOverhead
	dataList := [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("12"),
		[]byte("123"),
		[]byte("1234"),
		[]byte("12345"),
		[]byte("123456"),
		[]byte("1234567"),
		[]byte("12345678"),
		[]byte("123456789"),
		[]byte("1234567890"),
		[]byte("123456789012345"),
		[]byte("1234567890123456"),
		[]byte("12345678901234567"),
		bytes.Repeat([]byte("1234567890"), 100),
	}

	for _, origin := range dataList {
		ob := encrypt(origin)
		ret, errMsg := decrypt(ob)
		udwTest.Equal(errMsg, "", origin)
		udwTest.Equal(ret, origin)
	}

	if req.InReq.NoCorrectVerify == false {
		for _, origin := range dataList {
			output, errMsg := decrypt(origin)
			udwTest.Ok(errMsg != "", output, "no bad data find")
		}
	}
	origin := []byte("1234567890123456712345678901234567")
	resultMap := map[string]struct{}{}
	if req.InReq.RandomCheckSize == 0 {
		req.InReq.RandomCheckSize = 100
	}
	for i := 0; i < req.InReq.RandomCheckSize; i++ {
		result := encrypt(origin)
		resultMap[string(result)] = struct{}{}
	}
	if req.InReq.NoRandom == true {
		udwTest.Equal(len(resultMap), 1)
	} else {
		udwTest.Equal(len(resultMap), req.InReq.RandomCheckSize, maxOverhead)
	}

	if maxOverhead > 0 {
		for _, i := range []int{1, 10, 100, 1000, 10000} {
			ob := encrypt(udwRand.MustCryptoRandBytes(i))
			udwTest.Ok(len(ob)-i <= maxOverhead, i)
		}
	}

	ob := encrypt(origin)
	if req.InReq.NoCorrectVerify == false {
		for i := 0; i < len(ob); i++ {
			newOb := udwBytes.Clone(ob)
			newOb[i] = newOb[i] - 1
			_, errMsg := decrypt(newOb)
			udwTest.Ok(errMsg != "")
		}

		if len(ob) < 100 {
			for start := 0; start < len(ob)-1; start++ {
				for end := start + 1; end < len(ob); end++ {
					newOb := udwBytes.Clone(ob)
					newOb = newOb[start:end]
					_, errMsg := decrypt(newOb)
					udwTest.Ok(errMsg != "")
				}
			}
		} else {
			fmt.Println("encrypted data too large skip some cut check", len(ob))
			start := 0
			for end := start + 1; end < len(ob); end++ {
				newOb := udwBytes.Clone(ob)
				newOb = newOb[start:end]
				_, errMsg := decrypt(newOb)
				udwTest.Ok(errMsg != "")
			}
		}
	}

}
