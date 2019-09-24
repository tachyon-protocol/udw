package udwStrings

import "bytes"

func ByteListIndexAll(b []byte, keyword []byte) (outList []int) {
	pos := 0
	for {
		if pos >= len(b) {
			break
		}
		thisPos := bytes.Index(b[pos:], keyword)
		if thisPos == -1 {
			break
		}
		outList = append(outList, pos+thisPos)
		pos = pos + thisPos + len(keyword)
	}
	return outList
}
