package udwStrings

import "strings"

func Padding(s string, expectTotalSize int) string {
	if len(s) >= expectTotalSize {
		return s
	}
	return s + strings.Repeat(" ", expectTotalSize-len(s))
}
