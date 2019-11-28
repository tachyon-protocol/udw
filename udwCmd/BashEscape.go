package udwCmd

import (
	"bytes"
	"strings"
)

func BashEscape(inS string) (outS string) {

	return "'" + strings.Replace(inS, "'", `'\''`, -1) + "'"
}

func BashEscapeSlice(in []string) (outS string) {
	_buf := &bytes.Buffer{}
	for _, s := range in {
		_buf.WriteString(BashEscape(s))
		_buf.WriteByte(' ')
	}
	return _buf.String()
}
