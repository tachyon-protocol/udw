package udwQueryOnlyUrl

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"strings"
)

type QueryUrlObj struct {
	ProtocolName string
	Query        [][2]string
}

func (obj *QueryUrlObj) Marshal() string {
	bufW := udwBytes.BufWriter{}
	bufW.WriteString_(obj.ProtocolName)
	bufW.WriteString_("://?")
	if len(obj.Query) > 0 {
		for i, pair := range obj.Query {
			if i != 0 {
				bufW.WriteString_("&")
			}
			bufW.WriteString_(UrlvEncode(pair[0]))
			bufW.WriteString_("=")
			bufW.WriteString_(UrlvEncode(pair[1]))
		}
	}
	return bufW.GetString()
}
func (obj *QueryUrlObj) String() string {
	if obj == nil {
		return "<nil>"
	}
	return obj.Marshal()
}
func (obj *QueryUrlObj) AddKv(key string, value string) {
	obj.Query = append(obj.Query, [2]string{key, value})
}
func (obj *QueryUrlObj) GetFirstValueByKey(key string) string {
	for _, pair := range obj.Query {
		if pair[0] == key {
			return pair[1]
		}
	}
	return ""
}

func ParseQueryUrlObj(url string) *QueryUrlObj {
	obj := &QueryUrlObj{}
	index := strings.Index(url, "://?")
	if index == -1 || index == 0 {
		return nil
	}
	obj.ProtocolName = url[:index]
	queryS := url[index+4:]
	lastPos := 0
	lastK := ""
	lastV := ""
	finishVFn := func(i int) bool {
		if lastK == "" {
			return false
		}
		if lastV != "" {
			return false
		}
		v := queryS[lastPos:i]
		lastPos = i + 1
		v1, ok := urlDecode(v)
		if !ok {
			return false
		}
		lastV = v1
		obj.AddKv(lastK, lastV)
		lastK = ""
		return true
	}
	for i := 0; i < len(queryS); i++ {
		c := queryS[i]
		if c == '=' {
			if lastK != "" {
				return nil
			}
			v := queryS[lastPos:i]
			lastPos = i + 1
			v1, ok := urlDecode(v)
			if !ok {
				return nil
			}
			lastK = v1
			lastV = ""
			continue
		} else if c == '&' {
			ok := finishVFn(i)
			if !ok {
				return nil
			}
			continue
		}
	}
	ok := finishVFn(len(queryS))
	if !ok {
		return nil
	}
	return obj
}

const hexTable = "0123456789ABCDEF"

func UrlvEncode(s string) string {
	afterLen := len(s)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isPlainInUrlQuery(c) {
		} else {
			afterLen += 2
		}
	}
	if afterLen == len(s) {
		return s
	}
	out := make([]byte, afterLen)
	outPos := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isPlainInUrlQuery(c) {
			out[outPos] = c
			outPos++
		} else {
			out[outPos] = '%'
			out[outPos+1] = hexTable[c>>4]
			out[outPos+2] = hexTable[c&15]
			outPos += 3
		}
	}
	return string(out)
}

func urlDecode(s string) (out string, ok bool) {
	n := 0
	for i := 0; i < len(s); {
		c := s[i]
		if c == '%' {
			n++
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				return "", false
			}
			i += 3
		} else if isPlainInUrlQuery(c) {
			i++
		} else {
			return "", false
		}
	}

	if n == 0 {
		return s, true
	}

	bufW := udwBytes.BufWriter{}
	bufW.TryGrow(len(s) - 2*n)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '%':
			bufW.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
			i += 2
		default:
			bufW.WriteByte(s[i])
		}
	}
	return bufW.GetString(), true
}

func isPlainInUrlQuery(c byte) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') || c == '-' || c == '.' || c == '_'
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
