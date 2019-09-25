package udwRand

const ReadableAlphaNumMap = "23456789abcdefghjkmnpqrstuvwxyz"

func MustCryptoRandToReadableAlphaNum(length int) string {
	return MustCryptoRandFromByteList(length, ReadableAlphaNumMap)
}

func MustCryptoRandToReadableAlphaNumNoAlloc(length int, tmpBuf []byte) []byte {
	return MustCryptoRandFromByteListNoAlloc(length, ReadableAlphaNumMap, tmpBuf)
}

func EncodeReadableAlphaNumForRand(b []byte) string {
	outBytes := make([]byte, len(b)/2)
	mapLen := len(ReadableAlphaNumMap)
	for i := 0; i < len(outBytes); i++ {
		outBytes[i] = ReadableAlphaNumMap[(int(b[2*i])*256+int(b[2*i+1]))%(mapLen)]
	}
	return string(outBytes)
}

const readableAlphaMap = "abcdefghjkmnpqrstuvwxyz"

func MustCryptoRandToReadableAlpha(length int) string {
	return MustCryptoRandFromByteList(length, readableAlphaMap)
}

func EncodeReadableAlphaNumForRandNoAlloc(b []byte, tmpBuf []byte) string {
	if len(tmpBuf) < len(b)/2 {
		tmpBuf = make([]byte, len(b)/2)
	} else {
		tmpBuf = tmpBuf[:len(b)/2]
	}
	mapLen := len(ReadableAlphaNumMap)
	for i := 0; i < len(tmpBuf); i++ {
		tmpBuf[i] = ReadableAlphaNumMap[(int(b[2*i])*256+int(b[2*i+1]))%(mapLen)]
	}
	return string(tmpBuf)
}
