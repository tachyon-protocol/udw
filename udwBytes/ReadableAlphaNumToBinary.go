package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwStrings"
	"math/big"
)

var alphaNumRevertMap = [...]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x00, 0x11, 0x12, 0x00, 0x13, 0x14, 0x00, 0x15, 0x16,
	0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}

func ReadableAlphaNumToBinaryOrOriginToWriter(id string, bufWriter *BufWriter) {
	if len(id) == 0 {
		return
	}
	const minS = 50
	const maxS = 122
	getIndexFromMapFn := func(b byte) byte {
		if b < minS || b > maxS {
			return 0
		}
		return alphaNumRevertMap[b-minS]
	}
	buf := bufWriter.GetHeadBuffer(len(id))
	for i := 0; i < len(id); i++ {
		bs := byte(id[i])
		index := getIndexFromMapFn(bs)
		if index == 0 {
			bufWriter.WriteString(id)
			return
		}
		index = index - 1
		thisB := byte(0)
		if index < 10 {
			thisB = index + '0'
		} else {
			thisB = index - 10 + 'a'
		}
		buf[i] = thisB
	}
	s := udwStrings.GetStringFromByteArrayNoAlloc(buf)
	i := big.Int{}
	i.SetString(s, 31)

	bufWriter.Write(i.Bytes())
}

func ReadableAlphaNumToBinaryOrOriginToSlice(id string) []byte {
	bufW := &BufWriter{}
	ReadableAlphaNumToBinaryOrOriginToWriter(id, bufW)
	return bufW.GetBytes()
}

func ReadableAlphaNumFromBinary(b []byte, originLen int) string {
	bi := big.Int{}
	bi.SetBytes(b)
	bs := bi.Text(31)
	bufLen := len(bs)
	if originLen > bufLen {
		bufLen = originLen
	}
	diff := bufLen - len(bs)
	buf := make([]byte, bufLen)
	for i := 0; i < diff; i++ {
		buf[i] = '2'
	}
	for i := diff; i < bufLen; i++ {
		b := bs[i-diff]
		index := uint8(0)
		if b >= '0' && b <= '9' {
			index = b - '0'
		} else {
			index = b - 'a' + 10
		}
		b2 := readableAlphaNumMap[index]
		buf[i] = b2
	}
	return string(buf)
}

const readableAlphaNumMap = "23456789abcdefghjkmnpqrstuvwxyz"
