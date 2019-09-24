package udwRand

import (
	"crypto/rand"
	"encoding/hex"
)

func MustCryptoRandBytes(length int) []byte {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func MustCryptoRandBytesWithBuf(buf []byte) {
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
}

func MustCryptoRandToHex(length int) string {
	readLen := length/2 + length%2
	readBuf := make([]byte, readLen)
	_, err := rand.Read(readBuf[:readLen])
	if err != nil {
		panic(err)
	}
	outBuf := make([]byte, readLen*2)
	hex.Encode(outBuf, readBuf[:readLen])
	return string(outBuf[:length])
}

func MustCryptoRandToAlphaNum(length int) string {
	return MustCryptoRandToReadableAlphaNum(length)

}

const numMap = "0123456789"

func MustCryptoRandToNum(length int) string {
	return MustCryptoRandFromByteList(length, numMap)
}

func MustCryptoRandToNumWithRandomContent(randomContent []byte) string {
	return MustCryptoRandFromRandomContent(randomContent, numMap)
}

func MustCryptoRandFromByteList(length int, list string) string {

	return string(MustCryptoRandFromByteListNoAlloc(length, list, nil))
}

func MustCryptoRandFromByteListNoAlloc(length int, list string, tmpBuf []byte) []byte {
	if length*3 > len(tmpBuf) {
		tmpBuf = make([]byte, length*3)
	}
	var bytes = tmpBuf[length : 3*length]
	var outBytes = tmpBuf[0:length]
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	mapLen := len(list)
	for i := 0; i < length; i++ {
		outBytes[i] = list[(int(bytes[2*i])*256+int(bytes[2*i+1]))%(mapLen)]
	}
	return outBytes
}

func MustCryptoRandFromRandomContent(randomContent []byte, list string) string {
	outBytes := make([]byte, len(randomContent)/2)
	mapLen := len(list)
	for i := 0; i < len(outBytes); i++ {
		outBytes[i] = list[(int(randomContent[2*i])*256+int(randomContent[2*i+1]))%(mapLen)]
	}
	return string(outBytes)
}
