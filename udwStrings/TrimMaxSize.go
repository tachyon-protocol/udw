package udwStrings

import (
	"unicode/utf8"
)

func TrimMaxSize(s string, maxSize int) string {
	if len(s) <= maxSize {
		return s
	}
	if len(s) == 0 {
		return ""
	}
	index := 0
	for {
		_, size := utf8.DecodeRuneInString(s[index:])
		if index+size > maxSize {
			return s[:index]
		}
		index += size
	}
}

func TrimMaxSizeByte(s string, maxSize int) string {
	if len(s) < maxSize {
		return s
	}
	return s[:maxSize]
}

func TrimMaxSizeWithSubfix(s string, maxSize int, subfix string) string {
	if len(s) <= maxSize {
		return s
	}
	if len(s) == 0 {
		return ""
	}
	withSubfixMaxSize := maxSize - len(subfix)
	if withSubfixMaxSize <= 0 {

		return subfix[:maxSize]
	}
	index := 0
	for {
		_, size := utf8.DecodeRuneInString(s[index:])
		if index+size > withSubfixMaxSize {
			return s[:index] + subfix
		}
		index += size
	}
}
