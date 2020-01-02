package udwCryptoEncryptV3_test

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoEncryptV3"
	"github.com/tachyon-protocol/udw/udwCryptoTester"
	"github.com/tachyon-protocol/udw/udwIo"
	"testing"
)

func TestSymmetryConnOnce(ot *testing.T) {
	key1 := [32]byte{0xd8, 0x51, 0xea, 0x81, 0xb9, 0xe, 0xf, 0x2f, 0x8c, 0x85, 0x5f, 0xb6, 0x14, 0xb2}
	key := &key1
	block := udwCryptoEncryptV3.GetAesBlockFrom32Psk(key)

	thisEncrypt := func(data []byte) (output []byte) {
		bufW := udwBytes.BufWriter{}
		udwCryptoEncryptV3.MustSymmetryConnEncryptOnceWithBlock(data, block, &bufW)
		return bufW.GetBytes()
	}
	ReferenceDecrypt := func(data []byte) (output []byte, err error) {
		if len(data) < 20 {
			return nil, errors.New("len(data)<20")
		}
		reader := udwBytes.NewBufReader(data)
		conn1 := udwIo.StructWriterReaderCloser{
			Writer: udwIo.Nop,
			Reader: reader,
			Closer: udwIo.Nop,
		}
		conn2 := udwCryptoEncryptV3.NewSymmetryConnWithBlock(conn1, block)
		defer conn2.Close()
		output, err = udwIo.ReadAll(conn2)
		return output, err
	}
	udwCryptoTester.EncryptTesterL2(thisEncrypt, ReferenceDecrypt, 20)
}
