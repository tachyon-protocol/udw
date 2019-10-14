package udwGoTypeMarshal

import (
	"bytes"
	"strconv"
)

func WriteByteSlice(buf []byte) string {
	_buf := bytes.Buffer{}
	_buf.WriteString("[]byte{")
	for i, b := range buf {
		_buf.WriteString("0x")
		s := strconv.FormatInt(int64(b), 16)
		if len(s) == 1 {
			s = "0" + s
		}
		_buf.WriteString(s)
		_buf.WriteString(",")
		if i%16 == 15 {
			_buf.WriteString("\n")
		}
	}
	_buf.WriteString("}")
	return _buf.String()
}
