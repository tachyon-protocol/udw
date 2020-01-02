package udwRspLib

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestGoBuffer_ReadStringUTF16(t *testing.T) {
	_buf := &GoBuffer{}
	_buf.WriteStringUTF16("abc")
	_buf.ResetToRead()
	s := _buf.ReadStringUTF16()
	udwTest.Equal(s, "abc")

}
