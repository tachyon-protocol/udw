package udwStrings

func AsciiToLower(str string) string {
	bs := []byte(str)
	for i, b := range bs {
		if 'A' <= b && b <= 'Z' {
			bs[i] += 'a' - 'A'
		}
	}
	return string(bs)
}

func IsAllAscii(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if (b >= 0x20 && b <= 0x7e || b == '\n' || b == '\r' || b == '\t') == false {
			return false
		}
	}
	return true
}

func IsAllAlphabet(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 65 &&
			r <= 90) ||
			(r >= 97 &&
				r <= 122)) {
			return false
		}
	}
	return true
}

func IsAllNum(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func IsAllAlphabetNum(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 65 &&
			r <= 90) ||
			(r >= 97 &&
				r <= 122) ||
			(r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func IsGoodFileName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'A' &&
			r <= 'Z') ||
			(r >= 'a' &&
				r <= 'z') ||
			(r >= '0' && r <= '9') ||
			(r == '.') || (r == '-') || (r == '_')) {
			return false
		}
	}
	return true
}
