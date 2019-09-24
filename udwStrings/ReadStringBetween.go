package udwStrings

import (
	"bytes"
	"fmt"
	"strings"
)

func ReadFirstStringBetweenWithByteList(content []byte, startPart []byte, endPart []byte) []byte {
	indexStart := bytes.Index(content, startPart)
	if indexStart == -1 {
		return nil
	}
	indexStart = indexStart + len(startPart)
	indexEnd := bytes.Index(content[indexStart:], endPart)
	if indexEnd == -1 {
		return nil
	}
	indexEnd = indexStart + indexEnd
	return content[indexStart:indexEnd]
}

func MustStringBetweenFirstSubString(s string, sub1 string, sub2 string) string {
	pos1 := strings.Index(s, sub1)
	if pos1 == -1 {
		panic(fmt.Errorf("[MustStringBetweenFirstSubString] can not found sub1[%s] from [%s]", sub1, s))
	}
	pos2 := strings.Index(s[pos1+len(sub1):], sub2)
	if pos2 == -1 {
		panic(fmt.Errorf("[MustStringBetweenFirstSubString] can not found sub2[%s] from [%s]", sub2, s[pos1+len(sub1):]))
	}
	return s[pos1+len(sub1) : pos1+len(sub1)+pos2]
}

func StringBetweenFirstSubStringIgnoreNotExist(s string, sub1 string, sub2 string) string {
	pos1 := strings.Index(s, sub1)
	if pos1 == -1 {
		return ""
	}
	pos2 := strings.Index(s[pos1+len(sub1):], sub2)
	if pos2 == -1 {
		return ""
	}
	return s[pos1+len(sub1) : pos1+len(sub1)+pos2]
}
