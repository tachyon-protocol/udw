package udwBase64

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestStdDecode(ot *testing.T) {
	shouldB := bytes.Repeat([]byte{0}, 10)

	udwTest.Equal(MustStdBase64DecodeStringToByte("AAAAAAAAAAAAAA=="), shouldB)
	udwTest.Equal(MustStdBase64DecodeStringToByte("A\n\r\t AAAAAAAAAAAAA=="), shouldB)
}
