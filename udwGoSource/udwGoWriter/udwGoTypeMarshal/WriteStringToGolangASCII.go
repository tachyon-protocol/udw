package udwGoTypeMarshal

import (
	"bytes"
)

func WriteStringToGolangASCII(s string) string {
	if canUseHereDocASCII(s) {
		return "`" + s + "`"
	} else {
		return WriteStringToGolangDoubleQuotationASCII(s)
	}
}

func WriteStringToGolangDoubleQuotationASCII(s string) string {
	outBuf := &bytes.Buffer{}
	outBuf.WriteByte('"')
	l := len(s)
	i := 0
	for {
		if i >= l {
			outBuf.WriteByte('"')
			return outBuf.String()
		}
		b := s[i]
		switch b {
		case '\a':
			outBuf.WriteString(`\a`)
		case '\b':
			outBuf.WriteString(`\b`)
		case '\f':
			outBuf.WriteString(`\f`)
		case '\n':
			outBuf.WriteString(`\n`)
		case '\r':
			outBuf.WriteString(`\r`)
		case '\t':
			outBuf.WriteString(`\t`)
		case '\v':
			outBuf.WriteString(`\v`)
		case '\\':
			outBuf.WriteString(`\\`)
		case '"':
			outBuf.WriteString(`\"`)
		default:
			if b >= 0x20 && b <= 0x7e {

				outBuf.WriteByte(b)
			} else {
				outBuf.WriteString(`\x`)
				outBuf.WriteByte(lowerhex[s[i]>>4])
				outBuf.WriteByte(lowerhex[s[i]&0xF])
			}
		}
		i++
	}
}

func canUseHereDocASCII(s string) bool {
	l := len(s)
	i := 0
	for {
		if i >= l {
			return true
		}
		b := s[i]
		if !isASCIISafe(b) {
			return false
		}
		if b == '`' {
			return false
		}
		i += 1
	}
}

func isASCIISafe(b byte) bool {
	return b >= 0x20 && b <= 0x7e
}
