package udwFile

import "bytes"

func MustFileContainsString(path string, toSearchString string) bool {
	content := MustReadFile(path)
	toSearchByte := []byte(toSearchString)
	return bytes.Contains(content, toSearchByte)
}
