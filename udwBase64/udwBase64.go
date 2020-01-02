package udwBase64

import "encoding/base64"

func Base64EncodeStringToString(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func Base64EncodeByteToString(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

func MustBase64DecodeStringToString(input string) string {
	output, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return string(output)
}

func Base64DecodeStringToByte(input string) (b []byte, err error) {
	return base64.URLEncoding.DecodeString(input)
}

func MustUrlNoPaddingEncode(input []byte) []byte {
	return []byte(base64.RawURLEncoding.EncodeToString(input))
}
