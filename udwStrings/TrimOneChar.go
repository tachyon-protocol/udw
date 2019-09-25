package udwStrings

func TrimOneChar(s string, c byte) string {
	if len(s) == 0 {
		return s
	}
	startPos := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			startPos = i + 1
		} else {
			break
		}
	}
	if startPos >= len(s) {
		return ""
	}
	endPos := len(s)
	for i := len(s) - 1; i >= startPos; i-- {
		if s[i] == '.' {
			endPos = i
		} else {
			break
		}
	}

	return s[startPos:endPos]
}
