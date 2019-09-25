package udwStrings

import (
	"strings"
)

func StringAfterFirstSubString(s string, subString string) string {
	pos := strings.Index(s, subString)
	if pos == -1 {
		return ""
	}
	return s[pos+len(subString):]
}

func StringBeforeFirstSubString(s string, subString string) string {
	pos := strings.Index(s, subString)
	if pos == -1 {
		return ""
	}
	return s[:pos]
}

func StringAfterLastSubString(s string, subString string) string {
	pos := strings.LastIndex(s, subString)
	if pos == -1 {
		return ""
	}
	return s[pos+len(subString):]
}

func StringBeforeLastSubString(s string, subString string) string {
	pos := strings.LastIndex(s, subString)
	if pos == -1 {
		return ""
	}
	return s[:pos]
}

func StringSplitToTwoPart(s string, sep string) (string, string) {
	return StringBeforeFirstSubString(s, sep), StringAfterFirstSubString(s, sep)
}

func StringBeforeFirstSubStringOrInput(s string, subString string) string {
	pos := strings.Index(s, subString)
	if pos == -1 {
		return s
	}
	return s[:pos]
}
