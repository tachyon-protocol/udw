package udwStrings

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwMap"
	"strings"
	"unicode"
)

type SliceExistChecker map[string]bool

func (c SliceExistChecker) Add(s string) bool {
	ret, ok := c[s]
	if !ok {
		return false
	}
	if ret == false {
		c[s] = true
	}
	return true
}

func (c SliceExistChecker) FindNotExistOne() (NotExist string) {
	for s, ret := range c {
		if ret == false {
			return s
		}
	}
	return ""
}

func NewSliceExistChecker(slice ...string) SliceExistChecker {
	out := SliceExistChecker{}
	for _, s := range slice {
		out[s] = false
	}
	return out
}

func FirstLetterToUpper(s string) string {
	b := []byte(s)
	if len(b) == 0 {
		return s
	}
	if b[0] > unicode.MaxASCII {
		return s
	}
	firstLetter := bytes.ToUpper(b[0:1])
	b = append(firstLetter, b[1:]...)
	return string(b)
}

func MapStringBoolToSortedSlice(m map[string]bool) []string {
	return udwMap.MapStringBoolToSortedSlice(m)
}

func LastTwoPartSplit(originS string, splitS string) (p1 string, p2 string, ok bool) {
	part := strings.Split(originS, splitS)
	if len(part) < 2 {
		return "", "", false
	}
	return strings.Join(part[:len(part)-1], splitS), part[len(part)-1], true
}

func LineDataToSlice(lineData string) []string {
	part := strings.Split(lineData, "\n")
	out := []string{}
	for _, s := range part {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}
func SubStr(s string, from int, to int) string {
	rs := []rune(s)
	rl := len(rs)
	if to == 0 {
		to = rl
	}
	if to < 0 {
		to = rl + to
	}
	if to > rl {
		to = rl
	}
	if to < 0 {
		to = 0
	}
	return string(rs[from:to])
}

func RemoveWhiteSpace(in string) string {
	s1 := strings.Replace(in, "\n", "", -1)
	s2 := strings.Replace(s1, "\r", "", -1)
	s3 := strings.Replace(s2, "\t", "", -1)
	s4 := strings.Replace(s3, " ", "", -1)
	return s4
}
