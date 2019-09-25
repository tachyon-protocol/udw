package udwStrings

import "unicode/utf8"

func IsSafePrintByte(s []byte) bool {
	for i := range s {
		switch s[i] {
		case '\x00':
			return false
		}
	}
	return utf8.Valid(s)
}
