package udwGoTypeMarshal

import (
	"bytes"
	"unicode/utf8"
)

const lowerhex = "0123456789abcdef"

func WriteStringToGolang(s string) string {
	if canUseHereDoc(s) {
		return "`" + s + "`"
	} else {
		return WriteStringToGolangDoubleQuotation(s)
	}
}

func WriteStringToGolangDoubleQuotation(s string) string {
	outBuf := &bytes.Buffer{}
	WriteStringToGolangDoubleQuotationToBuf(s, outBuf)
	return outBuf.String()
}

func canUseHereDoc(s string) bool {
	l := len(s)
	i := 0
	for {
		if i >= l {
			return true
		}
		result, size := utf8.DecodeRuneInString(s[i:])
		if result == utf8.RuneError {
			return false
		}
		if result == '`' {
			return false
		}
		if result == 65279 {
			return false
		}
		if result == 0 {
			return false
		}
		if result == '\r' {
			return false
		}
		i += size
	}
}

func WriteStringToGolangToBuf(s string, _buf *bytes.Buffer) {
	if canUseHereDoc(s) {
		_buf.WriteByte('`')
		_buf.WriteString(s)
		_buf.WriteByte('`')
		return
	} else {
		WriteStringToGolangDoubleQuotationToBuf(s, _buf)
		return
	}
}

func WriteStringToGolangDoubleQuotationToBuf(s string, outBuf *bytes.Buffer) {
	outBuf.WriteByte('"')
	l := len(s)
	i := 0
	for {
		if i >= l {
			outBuf.WriteByte('"')
			return
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
		case '\x00':
			outBuf.WriteString(`\x00`)
		default:
			if b >= 0x20 && b <= 0x7e {

				outBuf.WriteByte(b)
			} else {
				result, size := utf8.DecodeRuneInString(s[i:])
				if result != utf8.RuneError && result != 65279 && result != 0 {
					outBuf.WriteRune(result)
					i += size - 1
				} else {
					outBuf.WriteString(`\x`)
					outBuf.WriteByte(lowerhex[s[i]>>4])
					outBuf.WriteByte(lowerhex[s[i]&0xF])
				}
			}
		}
		i++
	}
}
