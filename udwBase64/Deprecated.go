package udwBase64

import "encoding/base64"

func MustStdBase64DecodeString(s string) (out []byte) {
	return MustStdBase64DecodeStringToByte(s)
}

func StdBase64Decode(s []byte) (out []byte, err error) {
	return StdBase64DecodeByteToByte(s)
}

func MustBase64EncodeStringToString(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}
