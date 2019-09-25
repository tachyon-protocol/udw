package udwStrings

func GetCountByteNumber(s string, b byte) int {
	number := 0
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			number++
		}
	}
	return number
}
