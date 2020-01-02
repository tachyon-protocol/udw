package udwBase64

import (
	"encoding/base64"
	"strings"
)

func EncodeByteToStringV2(input []byte) string {
	out := base64.URLEncoding.EncodeToString([]byte(input))
	return strings.Replace(out, "=", "", -1)
}

func DecodeStringToByteV2(input string) (b []byte, err error) {
	if len(input)%4 != 0 {
		input += strings.Repeat("=", 4-len(input)%4)
	}
	return base64.URLEncoding.DecodeString(input)
}
