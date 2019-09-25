package udwStrings

func GetMonospaceWidth(str string) int {
	count := 0
	for _, r := range str {
		count++
		if r >= 0x2E80 && r <= 0x2EFF {
			count++
		}
		if r >= 0x3000 && r <= 0x303F {
			count++
		}
		if r >= 0x31C0 && r <= 0x31EF {
			count++
		}
		if r >= 0x3200 && r <= 0x32FF {
			count++
		}
		if r >= 0x3300 && r <= 0x33FF {
			count++
		}
		if r >= 0x3400 && r <= 0x4DBF {
			count++
		}
		if r >= 0xF900 && r <= 0xFAFF {
			count++
		}

		if r >= 0x4E00 && r <= 0x9FFF {
			count++
		}
		if r >= 0xFE30 && r <= 0xFE4F {
			count++
		}

	}
	return count
}
