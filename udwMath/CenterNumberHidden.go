package udwMath

func HiddenCenterString(info string) string {
	str := info
	if len(str) < 3 {
		return info
	}
	bytes := []byte(str)
	length := len(bytes)/3 + 1
	for i, _ := range bytes {
		if i > (len(bytes)-length)/2 && i <= (len(bytes)+length)/2 && bytes[i] != '@' {
			bytes[i] = '*'
		}
	}
	return string(bytes)
}
