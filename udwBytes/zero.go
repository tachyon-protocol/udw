package udwBytes

func IsAllZero(b []byte) bool {
	for i := range b {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func FillZero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
