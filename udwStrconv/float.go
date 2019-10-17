package udwStrconv

import (
	"strconv"
)

func FormatFloat64ToFInLen(f float64, showLen int) string {
	s1 := strconv.FormatFloat(f, 'f', 0, 64)
	if showLen == 1 && s1 == "-0" {
		return "0"
	}
	if len(s1) >= showLen {
		return s1
	}

	if len(s1)+1 == showLen {
		if s1[0] == '-' {
			return "-0" + s1[1:]
		} else {
			return "0" + s1
		}
	}
	toSize := showLen - len(s1) - 1
	s2 := strconv.FormatFloat(f, 'f', toSize, 64)
	if len(s2) >= showLen {
		return s2
	}

	toSize = toSize + showLen - len(s2)
	return strconv.FormatFloat(f, 'f', toSize, 64)
}
