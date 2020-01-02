package udwCryptoTester

import "github.com/tachyon-protocol/udw/udwErr"

func EncryptTester(encrypt func(key *[32]byte, data []byte) (output []byte),
	decrypt func(key *[32]byte, data []byte) (output []byte, err error),
	maxOverhead int) {
	EncryptTesterErrMsg(EncryptTesterErrMsgRequest{
		Encrypt: encrypt,
		Decrypt: func(key *[32]byte, data []byte) (output []byte, errMsg string) {
			output, err := decrypt(key, data)
			return output, udwErr.ErrorToMsg(err)
		},
		MaxOverhead: maxOverhead,
	})
}

func EncryptTesterL2(Encrypt func(data []byte) (output []byte), Decrypt func(data []byte) (output []byte, err error), MaxOverhead int) {
	encryptTesterErrMsgNoKey(EncryptTesterErrMsgNoKeyRequest{
		Encrypt: Encrypt,
		Decrypt: func(data []byte) (output []byte, errMsg string) {
			output, err := Decrypt(data)
			return output, udwErr.ErrorToMsg(err)
		},
		MaxOverhead: MaxOverhead,
	})
}

func encryptTesterErrMsgNoKey(req EncryptTesterErrMsgNoKeyRequest) {
	encryptTesterL2(encryptTesterL2Request{
		Encrypt: req.Encrypt,
		Decrypt: req.Decrypt,
		InReq: EncryptTesterErrMsgRequest{
			MaxOverhead: req.MaxOverhead,
		},
	})
}
