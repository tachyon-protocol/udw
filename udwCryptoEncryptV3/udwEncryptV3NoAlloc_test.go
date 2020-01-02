package udwCryptoEncryptV3_test

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoEncryptV3"
	"github.com/tachyon-protocol/udw/udwCryptoTester"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestEncryptV3NoAlloc(ot *testing.T) {
	key := &[32]byte{0xd8, 0x51, 0xea, 0x81, 0xb9, 0xe, 0xf, 0x2f, 0x8c, 0x85, 0x5f, 0xb6, 0x14, 0xb2}
	obj := udwCryptoEncryptV3.NewUdwEncryterV3NoAlloc(key)
	udwCryptoTester.EncryptTesterL2(func(data []byte) (output []byte) {
		return obj.EncryptWithBuf(data, nil)
	}, func(data []byte) (output []byte, err error) {
		return obj.DecryptWithBuf(data, nil)
	}, 20)

	udwCryptoTester.EncryptTesterL2(func(data []byte) (output []byte) {
		return obj.EncryptWithBuf(data, nil)
	}, func(data []byte) (output []byte, err error) {
		return udwCryptoEncryptV3.DecryptV3(key, data)
	}, 20)

	udwCryptoTester.EncryptTesterL2(func(data []byte) (output []byte) {
		return udwCryptoEncryptV3.EncryptV3(key, data)
	}, func(data []byte) (output []byte, err error) {
		return obj.DecryptWithBuf(data, nil)
	}, 20)

	encryptBuf := make([]byte, 1000)
	decryptBuf := make([]byte, 1000)
	out, err := obj.DecryptWithBuf(obj.EncryptWithBuf([]byte("1"), encryptBuf), decryptBuf)
	udwTest.Equal(err, nil)
	udwTest.Equal(out, []byte("1"))
}

func TestAesMessageV6Encrypt(ot *testing.T) {
	key1 := [32]byte{0xd8, 0x51, 0xea, 0x81, 0xb9, 0xe, 0xf, 0x2f, 0x8c, 0x85, 0x5f, 0xb6, 0x14, 0xb2}
	key := &key1
	obj := udwCryptoEncryptV3.NewUdwEncryterV3NoAlloc(key)
	keyObj := udwCryptoEncryptV3.NewAesMessageV6KeyObject(key)

	thisEncrypt := func(data []byte) (output []byte) {
		bufW := udwBytes.BufWriter{}
		udwCryptoEncryptV3.AesMessageV6Encrypt(udwCryptoEncryptV3.AesMessageV6Request{
			KeyObject:  keyObj,
			Data:       data,
			OutputBufW: &bufW,
		})
		return bufW.GetBytes()
	}
	thisDecrypt := func(data []byte) (output []byte, err error) {
		bufW := udwBytes.BufWriter{}
		errMsg := udwCryptoEncryptV3.AesMessageV6Decrypt(udwCryptoEncryptV3.AesMessageV6Request{
			KeyObject:  keyObj,
			Data:       data,
			OutputBufW: &bufW,
		})

		return bufW.GetBytes(), udwErr.ErrorMsgToErr(errMsg)
	}
	ReferenceEecrypt := func(data []byte) (output []byte) {
		return obj.EncryptWithBuf(data, nil)
	}
	ReferenceDecrypt := func(data []byte) (output []byte, err error) {
		return obj.DecryptWithBuf(data, nil)
	}
	udwCryptoTester.EncryptTesterL2(thisEncrypt, ReferenceDecrypt, 20)
	udwCryptoTester.EncryptTesterL2(ReferenceEecrypt, thisDecrypt, 20)
	udwCryptoTester.EncryptTesterL2(thisEncrypt, thisDecrypt, 20)

	const benchSize = 1e4
	const dataSize = 1025
	bufW := udwBytes.BufWriter{}
	data := bytes.Repeat([]byte{55}, dataSize)
	benchThisEncryptRunning := func() {
		udwCryptoEncryptV3.AesMessageV6Encrypt(udwCryptoEncryptV3.AesMessageV6Request{
			KeyObject:  keyObj,
			Data:       data,
			OutputBufW: &bufW,
		})
	}
	benchThisEncryptRunning()
	thisEncryptResult := bufW.GetBytesClone()
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetName("thisEncrypt")
		udwTest.BenchmarkSetNum(benchSize)
		udwTest.BenchmarkSetBytePerRun(dataSize)
		for i := 0; i < benchSize; i++ {
			benchThisEncryptRunning()
		}
	})
	buf := make([]byte, dataSize*2)
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetName("ReferenceEncrypt")
		udwTest.BenchmarkSetNum(benchSize)
		udwTest.BenchmarkSetBytePerRun(dataSize)
		for i := 0; i < benchSize; i++ {
			obj.EncryptWithBuf(data, buf)
		}
	})
	benchThisDecrpt := func() {
		udwCryptoEncryptV3.AesMessageV6Decrypt(udwCryptoEncryptV3.AesMessageV6Request{
			KeyObject:  keyObj,
			Data:       thisEncryptResult,
			OutputBufW: &bufW,
		})
	}
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetName("thisDecrypt")
		udwTest.BenchmarkSetNum(benchSize)
		udwTest.BenchmarkSetBytePerRun(dataSize)
		for i := 0; i < benchSize; i++ {
			benchThisDecrpt()
		}
	})
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetName("ReferenceDecrypt")
		udwTest.BenchmarkSetNum(benchSize)
		udwTest.BenchmarkSetBytePerRun(dataSize)
		for i := 0; i < benchSize; i++ {
			obj.DecryptWithBuf(thisEncryptResult, buf)
		}
	})
}
