package udwJsonLib

import (
	"encoding/hex"
	"github.com/tachyon-protocol/udw/udwStrings"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func ReaderSetData(ctx *Context, data []byte) {
	ctx.readerData = data
	ctx.readerPos = 0
}

func ReadJsonString(ctx *Context) string {
	bs := readJsonStringToByteSlice(ctx)
	return string(bs)
}

func readJsonStringToByteSlice(ctx *Context) []byte {
	pos := ctx.readerPos
for1:
	for {
		if pos >= len(ctx.readerData) {
			panic(`[readJsonStringToWriter] need a " but get EOF pos ` + strconv.Itoa(pos))
		}
		c := ctx.readerData[pos]
		pos++
		switch c {
		case ' ', '\t', '\n', '\r':
		case '"':
			break for1
		case 'n':

			panic(`[readJsonStringToWriter] need a " but get ` + string(c) + ` pos ` + strconv.Itoa(pos) + ` ` + ComfirmGolangJsonLibBugTag)
		default:
			start := 0
			if pos-15 >= 0 {
				start = pos - 15
			}
			panic(`[readJsonStringToWriter] need a " but get ` + string(c) + ` pos ` + strconv.Itoa(pos) + `, content before 15 bytes:` + string(ctx.readerData[start:pos]))
		}
	}

	ctx.readerPos = pos
	startPos := pos
	for {
		if pos >= len(ctx.readerData) {
			break
		}
		c := ctx.readerData[pos]
		if c == '\\' {
			break
		} else if c < ' ' {

			panic(`[readJsonStringToWriter] invaild char1 [` + strconv.Itoa(int(c)) + `]`)
		}
		if c == '\\' || c < ' ' {
			break
		} else if c == '"' {
			pos++
			ctx.readerPos = pos
			buf := ctx.readerData[startPos : ctx.readerPos-1]
			return buf
		} else if c < utf8.RuneSelf {
			pos++
			continue
		} else {
			rr, size := utf8.DecodeRune(ctx.readerData[pos:])
			if rr == utf8.RuneError && size == 1 {

				panic(`[readJsonStringToWriter] ` + ComfirmGolangJsonLibBugTag + ` invalid utf8 code point1 [` + strconv.Itoa(int(c)) + `] pos [` + strconv.Itoa(pos) + `]`)
			}
			pos += size
		}
	}
	writerReset(ctx)
	for {
		c := ReaderReadByte(ctx)
		switch {
		case c == '"':
			return ctx.writerData[:ctx.writerPos]
		case c == '\\':
			b2 := ReaderReadByte(ctx)
			switch b2 {
			default:
				panic(`[readJsonStringToWriter] can not understand char "` + string(c) + `" after "\"`)
			case '"', '\\', '/', '\'':
				WriterWriteByte(ctx, b2)
			case 'b':
				WriterWriteByte(ctx, '\b')
			case 'f':
				WriterWriteByte(ctx, '\f')
			case 'n':
				WriterWriteByte(ctx, '\n')
			case 'r':
				WriterWriteByte(ctx, '\r')
			case 't':
				WriterWriteByte(ctx, '\t')
			case 'u':
				bHex := readerReadByteBySize(ctx, 4)
				r1Code, err := strconv.ParseUint(udwStrings.GetStringFromByteArrayNoAlloc(bHex), 16, 64)
				if err != nil {
					panic(`[readJsonStringToWriter] 1can not understand char "` + string(bHex) + `" after "\u"`)
				}
				r1 := rune(r1Code)
				if utf16.IsSurrogate(r1) {
					bHex := readerReadByteBySize(ctx, 6)
					if bHex[0] != '\\' || bHex[1] != 'u' {
						panic(`[readJsonStringToWriter] not another "\u" after a two uft16 code [` + string(bHex) + `] ` + ComfirmGolangJsonLibBugTag)
					}
					r2Code, err := strconv.ParseUint(udwStrings.GetStringFromByteArrayNoAlloc(bHex[2:]), 16, 64)
					if err != nil {
						panic(`[readJsonStringToWriter] 2can not understand char "` + string(bHex) + `" after "\u"`)
					}
					r2 := rune(r2Code)
					dec := utf16.DecodeRune(r1, r2)
					if dec == unicode.ReplacementChar {
						panic(`[readJsonStringToWriter] invaild unicode point1 ` + ComfirmGolangJsonLibBugTag)
					}
					buf := writerGetHeadBuffer(ctx, 4)
					writeSize := utf8.EncodeRune(buf, dec)
					writerAddPos(ctx, writeSize)
				}
				buf := writerGetHeadBuffer(ctx, 4)
				writeSize := utf8.EncodeRune(buf, r1)
				writerAddPos(ctx, writeSize)
			}

		case c < ' ':

			panic(`[readJsonStringToWriter] invaild char2 [` + strconv.Itoa(int(c)) + `]`)

		case c < utf8.RuneSelf:
			WriterWriteByte(ctx, c)
		default:

			rr, size := utf8.DecodeRune(ctx.readerData[ctx.readerPos-1:])
			if rr == utf8.RuneError && size == 1 {

				panic(`[readJsonStringToWriter] ` + ComfirmGolangJsonLibBugTag + ` invalid utf8 code point2 [` + strconv.Itoa(int(c)) + `] pos [` + strconv.Itoa(pos) + `]`)
			}
			WriterWriteByteList(ctx, ctx.readerData[ctx.readerPos-1:ctx.readerPos-1+size])
			ctx.readerPos += size - 1

		}
	}
}

func ReaderReadSpace(ctx *Context) {
	pos := ctx.readerPos
	for {
		if pos >= len(ctx.readerData) {
			return
		}
		c := ctx.readerData[pos]
		pos++
		switch c {
		case ' ', '\t', '\n', '\r':
		default:
			ctx.readerPos = pos - 1
			return
		}
	}
}

func (ctx *Context) ReaderReadSpaceAndIsEof() bool {
	pos := ctx.readerPos
	for {
		if pos >= len(ctx.readerData) {
			return true
		}
		c := ctx.readerData[pos]
		pos++
		switch c {
		case ' ', '\t', '\n', '\r':
		default:
			ctx.readerPos = pos - 1
			return false
		}
	}
}

func ReaderReadByte(ctx *Context) byte {
	if ctx.readerPos >= len(ctx.readerData) {
		panic("EOF " + string(ctx.readerData))
	}
	b := ctx.readerData[ctx.readerPos]
	ctx.readerPos++
	return b
}

func ReaderReadBack(ctx *Context, size int) {
	ctx.readerPos -= size
}
func ReaderReadPosAdd(ctx *Context, size int) {
	ctx.readerPos += size
}

func ReadJsonBool(ctx *Context) bool {
	ReaderReadSpace(ctx)
	buf := readerReadByteBySize(ctx, 4)
	switch {
	case buf[0] == 't' && buf[1] == 'r' && buf[2] == 'u' && buf[3] == 'e':
		return true
	case buf[0] == 'f' && buf[1] == 'a' && buf[2] == 'l' && buf[3] == 's':
		b := ReaderReadByte(ctx)
		if b == 'e' {
			return false
		}
	case buf[0] == 'n' && buf[1] == 'u' && buf[2] == 'l' && buf[3] == 'l':
		panic("[ReadJsonBool] fail [" + string(buf) + "] " + ComfirmGolangJsonLibBugTag)
	}
	panic("[ReadJsonBool] fail [" + string(buf) + "]")
}

func readNumberPartToCtx(ctx *Context) {
	writerReset(ctx)
for1:
	for {
		if ctx.readerPos >= len(ctx.readerData) {
			return
		}
		b := ctx.readerData[ctx.readerPos]
		ctx.readerPos++

		switch b {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.', 'e', 'E':
			WriterWriteByte(ctx, b)
		case ' ', '\n', '\t', '\r':
		default:
			ReaderReadBack(ctx, 1)
			break for1
		}
	}
}

func ReadJsonInt64(ctx *Context) int64 {
	readNumberPartToCtx(ctx)
	s := writerGetTmpString(ctx)
	if s == "" {
		panic("[ReadJsonInt64] read json with empty obj " + ComfirmGolangJsonLibBugTag)
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return int64(i)
}

func ReadJsonUint64(ctx *Context) uint64 {
	readNumberPartToCtx(ctx)
	s := writerGetTmpString(ctx)
	if s == "" {
		panic("[ReadJsonUint64] read json with empty obj " + ComfirmGolangJsonLibBugTag)
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint64(i)
}

func ReadJsonInt64FromString(ctx *Context) int64 {
	s := ReadJsonTmpString(ctx)
	if s == "" {
		panic("[ReadJsonInt64FromString] read json with empty obj " + ComfirmGolangJsonLibBugTag)
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return int64(i)
}

func ReadJsonUint64FromString(ctx *Context) uint64 {
	s := ReadJsonTmpString(ctx)
	if s == "" {
		panic("[ReadJsonUint64FromString] read json with empty obj " + ComfirmGolangJsonLibBugTag)
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint64(i)
}

func ReadJsonFloat64(ctx *Context) float64 {
	readNumberPartToCtx(ctx)
	s := writerGetTmpString(ctx)
	if s == "" {
		panic("[ReadJsonFloat64] read json with empty obj " + ComfirmGolangJsonLibBugTag + " " + strconv.Itoa(len(ctx.readerData)) + " " + hex.EncodeToString(ctx.readerData))
	}
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	if i == 0 {

		return 0
	}
	return i
}

func ReadJsonFloat32(ctx *Context) float32 {
	readNumberPartToCtx(ctx)
	s := writerGetTmpString(ctx)
	if s == "" {
		panic("[ReadJsonFloat32] read json with empty obj " + ComfirmGolangJsonLibBugTag)
	}
	i, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	if i == 0 {

		return 0
	}
	return float32(i)
}

func ReaderSkipMapValue(ctx *Context) {
	isInString := false
	isLastSplash := false
	bracketsLevel := 0
	for {
		b := ReaderReadByte(ctx)
		if isInString {
			switch b {
			case '"':
				if isLastSplash == false {
					isInString = false
				}
				isLastSplash = false
			case '\\':

				if isLastSplash == true {
					isLastSplash = false
				} else {
					isLastSplash = true
				}
			default:
				isLastSplash = false
			}
			continue
		}
		switch b {
		case '"':
			isInString = true
			isLastSplash = false
		case ',':
			if bracketsLevel == 0 {
				return
			}
		case '}', ']':
			if bracketsLevel == 0 {
				ReaderReadBack(ctx, 1)
				return
			}
			bracketsLevel--
		case '{', '[':
			bracketsLevel++
		}
	}
}

func ReaderReadMapValueByteList(ctx *Context) []byte {
	isInString := false
	isLastSplash := false
	bracketsLevel := 0
	startPos := ctx.readerPos
	for {
		b := ReaderReadByte(ctx)
		if isInString {
			switch b {
			case '"':
				if !isLastSplash {
					isInString = false
				}
				isLastSplash = false
			case '\\':

				if isLastSplash == true {
					isLastSplash = false
				} else {
					isLastSplash = true
				}
			default:
				isLastSplash = false
			}
			continue
		}
		switch b {
		case '"':
			isInString = true
			isLastSplash = false
		case ',':
			if bracketsLevel == 0 {
				return ctx.readerData[startPos : ctx.readerPos-1]
			}
		case '}', ']':
			if bracketsLevel == 0 {
				ReaderReadBack(ctx, 1)
				return ctx.readerData[startPos:ctx.readerPos]
			}
			bracketsLevel--
		case '{', '[':
			bracketsLevel++
		}
	}
}

func readerReadByteBySize(ctx *Context, size int) []byte {
	out := ctx.readerData[ctx.readerPos : ctx.readerPos+size]
	ctx.readerPos += size
	return out
}

func ReaderGetRemainByteSlice(ctx *Context) []byte {
	return ctx.readerData[ctx.readerPos:]
}

func MustReadJsonNull(ctx *Context) {
	buf := readerReadByteBySize(ctx, 4)
	if buf[0] == 'n' && buf[1] == 'u' && buf[2] == 'l' && buf[3] == 'l' {
		return
	}
	panic("[MustReadJsonNull] fail to read a null " + string(buf))
}
