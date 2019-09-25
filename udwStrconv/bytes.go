package udwStrconv

import "errors"

func ByteSliceParseUint(buf []byte) (int, error) {
	v, n, err := parseUintBuf(buf)
	if n != len(buf) {
		return -1, errors.New("UnexpectedTrailingChar")
	}
	return v, err
}

func parseUintBuf(b []byte) (int, int, error) {
	n := len(b)
	if n == 0 {
		return -1, 0, errors.New("EmptyInt")
	}
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return -1, i, errors.New("UnexpectedFirstChar")
			}
			return v, i, nil
		}
		if i >= 18 {
			return -1, i, errors.New("TooLongInt")
		}
		v = 10*v + int(k)
	}
	return v, n, nil
}
