package udwStrings

import "strings"

func StringFieldsAtPos(s string, pos int) string {
	sList := strings.Fields(s)
	if pos < len(sList) {
		return sList[pos]
	}
	return ""
}

func Fields(s string) []string {
	return strings.Fields(s)
}
