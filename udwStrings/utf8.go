package udwStrings

import (
	"unicode/utf8"
)

const (
	runeSelf = 0x80

	locb = 0x80
	hicb = 0xBF

	xx = 0xF1
	as = 0xF0
	s1 = 0x02
	s2 = 0x13
	s3 = 0x03
	s4 = 0x23
	s5 = 0x34
	s6 = 0x04
	s7 = 0x44
)

var first = [256]uint8{

	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as,

	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx,
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx,
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx,
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx,
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1,
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1,
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3,
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx,
}

type acceptRange struct {
	lo uint8
	hi uint8
}

var acceptRanges = [...]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}

func IsAllPrintableUtf8(p []byte) bool {

	n := len(p)
	for i := 0; i < n; {
		pi := p[i]
		if pi == 0 {

			return false
		}
		if pi < runeSelf {
			i++
			continue
		}
		x := first[pi]
		if x == xx {
			return false
		}
		size := int(x & 7)
		if i+size > n {
			return false
		}
		accept := acceptRanges[x>>4]
		if c := p[i+1]; c < accept.lo || accept.hi < c {
			return false
		} else if size == 2 {
		} else if c := p[i+2]; c < locb || hicb < c {
			return false
		} else if size == 3 {
		} else if c := p[i+3]; c < locb || hicb < c {
			return false
		}
		i += size
	}
	return true
}

func IsAllDecodeByUtf8(p []byte) bool {
	n := len(p)
	for i := 0; i < n; {
		c, size := utf8.DecodeRune(p[i:])
		if c == utf8.RuneError && size == 1 {
			return false
		}
		i += size
	}
	return true
}

func IsAllDecodeByUtf8String(p string) bool {

	return IsAllDecodeByUtf8([]byte(p))
}
