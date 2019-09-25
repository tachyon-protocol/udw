package udwBytes

func Clone(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

func CutToMaxLen(b []byte, l int) []byte {
	if len(b) < l {
		return b
	}
	return b[:l]
}
