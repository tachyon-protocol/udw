package udwStrconv

import "strconv"

func FormatFloat64ToFInLen(f float64, showLen int) string {
	s1 := strconv.FormatFloat(f, 'f', 0, 64)
	if len(s1)+1 >= showLen {
		return s1
	}
	return strconv.FormatFloat(f, 'f', showLen-len(s1)-1, 64)
}
