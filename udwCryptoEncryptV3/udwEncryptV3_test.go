package udwCryptoEncryptV3_test

import (
	"github.com/tachyon-protocol/udw/udwCryptoEncryptV3"
	"github.com/tachyon-protocol/udw/udwCryptoTester"
	"testing"
)

func TestEncryptV3(ot *testing.T) {
	udwCryptoTester.EncryptTesterErrMsg(udwCryptoTester.EncryptTesterErrMsgRequest{
		Encrypt:         udwCryptoEncryptV3.EncryptV3,
		DecryptErr:      udwCryptoEncryptV3.DecryptV3,
		MaxOverhead:     20,
		NoCorrectVerify: true,
	})
}
