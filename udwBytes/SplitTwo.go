package udwBytes

import "bytes"

func SplitTwoBetweenFirst(s []byte, subString []byte) (before []byte, after []byte) {
	pos := bytes.Index(s, subString)
	if pos == -1 {
		return s, nil
	}
	return s[:pos], s[pos+len(subString):]
}

func BytesBeforeFirstSub(s []byte, subString []byte) (b []byte) {
	pos := bytes.Index(s, subString)
	if pos == -1 {
		return nil
	}
	return s[:pos]
}
