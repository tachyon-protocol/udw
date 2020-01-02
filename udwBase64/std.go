package udwBase64

import "encoding/base64"

func MustStdBase64EncodeByteToByte(s []byte) (out []byte) {
	out = make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(out, s)
	return out
}

func MustStdBase64EncodeByteToString(s []byte) (out string) {
	return string(MustStdBase64EncodeByteToByte(s))
}

func MustStdBase64DecodeStringToByte(s string) (out []byte) {
	out, err := StdBase64DecodeByteToByte([]byte(s))
	if err != nil {
		panic(err)
	}
	return
}

func StdBase64DecodeByteToByte(s []byte) (out []byte, err error) {
	out = make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	s1 := make([]byte, len(s))
	s1i := 0
	for i := range s {
		c := s[i]
		switch c {
		case '\n', '\r', '\t', ' ':
		default:
			s1[s1i] = c
			s1i++
		}
	}
	nw, err := base64.StdEncoding.Decode(out, s1[:s1i])
	if err != nil {
		return nil, err
	}
	return out[:nw], nil
}
