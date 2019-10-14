package udwGoTypeMarshal

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwMap"
)

func WriteGoConst(pairList []udwMap.KeyValuePair) string {
	_buf := bytes.Buffer{}
	_buf.WriteString(`const (
`)
	for _, pair := range pairList {
		_buf.WriteString(pair.Key + " = " + WriteStringToGolang(pair.Value) + "\n")
	}
	_buf.WriteString(`)
`)
	return _buf.String()
}

func WriteGoConstByMap(m map[string]string) string {
	pairList := udwMap.MapStringStringToKeyValuePairListAes(m)
	return WriteGoConst(pairList)
}
