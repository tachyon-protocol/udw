package udwHex

import (
	"encoding/hex"
	"strings"
)

func UpperEncodeBytesToString(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

func EncodeBytesToString(b []byte) string {
	return hex.EncodeToString(b)
}

func EncodeStringToString(s string) string {
	return hex.EncodeToString([]byte(s))
}

func DecodeStringToString(s string) (string, error) {
	b, err := hex.DecodeString(s)
	return string(b), err
}

func MustDecodeStringToString(s string) string {
	b, err := DecodeStringToString(s)
	if err != nil {
		panic(err)
	}
	return b
}
func MustDecodeStringToByteArray(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

const HexTableLower = "0123456789abcdef"

func FromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func IsHexChar(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func MustDecodeStringToByteArrayIgnoreNonHex(s string) (b []byte) {
	s2 := make([]byte, len(s))
	s2Pos := 0
	for i := 0; i < len(s); i++ {
		if IsHexChar(s[i]) {
			s2[s2Pos] = s[i]
			s2Pos++
		}
	}
	s2 = s2[:s2Pos]
	b = make([]byte, hex.DecodedLen(len(s2)))
	_, err := hex.Decode(b, s2)
	if err != nil {
		panic(err)
	}
	return b
}
