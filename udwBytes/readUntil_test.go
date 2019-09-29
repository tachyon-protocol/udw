package udwBytes

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestReadUtil(t *testing.T) {
	var resp ReadUtilResponse
	toReadBuf := make([]byte, 1024)
	resp = ReadUtil(ReadUtilRequest{
		R:         bytes.NewReader([]byte("123")),
		Sub:       []byte("2"),
		ToReadBuf: toReadBuf,
	})
	udwTest.Equal(resp.TotalReadLen, 3)
	udwTest.Equal(resp.ErrMsg, "")
	udwTest.Equal(resp.SubStringStartPos, 1)
}

func TestReadBufToZeroByte(t *testing.T) {
	udwTest.Equal(ReadCStringFromBufToByte(nil), nil)
	udwTest.Equal(ReadCStringFromBufToByte([]byte{}), []byte{})
	udwTest.Equal(ReadCStringFromBufToByte([]byte{0}), []byte{})
	udwTest.Equal(ReadCStringFromBufToByte([]byte{1, 0}), []byte{1})
	udwTest.Equal(ReadCStringFromBufToByte([]byte("abc\x00\x00")), []byte("abc"))
	udwTest.Equal(ReadCStringFromBuf([]byte("abc\x00\x00")), "abc")
	udwTest.Equal(ReadCStringFromBuf([]byte("abc\x00abcd\x00")), "abc")
}
