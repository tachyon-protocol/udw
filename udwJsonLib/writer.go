package udwJsonLib

import (
	"encoding/base64"
	"math"
	"strconv"
	"unicode/utf8"
)

func WriteJsonString(ctx *Context, s string) {
	WriterWriteByte(ctx, '"')
	start := 0
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			if getHtmlSafeSetV2(b) {
				i++
				continue
			}
			if start < i {
				WriterWriteString(ctx, s[start:i])
			}
			switch b {
			case '\\', '"':
				WriterWriteByte(ctx, '\\')
				WriterWriteByte(ctx, b)
			case '\n':
				WriterWriteByte(ctx, '\\')
				WriterWriteByte(ctx, 'n')
			case '\r':
				WriterWriteByte(ctx, '\\')
				WriterWriteByte(ctx, 'r')
			case '\t':
				WriterWriteByte(ctx, '\\')
				WriterWriteByte(ctx, 't')
			default:

				WriterWriteString(ctx, `\u00`)
				WriterWriteByte(ctx, hextable[b>>4])
				WriterWriteByte(ctx, hextable[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			panic("[WriteJsonString] invalid utf8 string b [" + strconv.Itoa(int(b)) + "] pos [" + strconv.Itoa(i) + "]")

		}

		if c == '\u2028' || c == '\u2029' {
			if start < i {
				WriterWriteString(ctx, s[start:i])
			}
			WriterWriteString(ctx, `\u202`)
			WriterWriteByte(ctx, hextable[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		WriterWriteString(ctx, s[start:])
	}
	WriterWriteByte(ctx, '"')
	return
}

func WriteJsonInt64(ctx *Context, i int64) {
	writerTryRealloc(ctx, 21)
	outByte := strconv.AppendInt(writerGetHeadBuffer(ctx, 0), i, 10)
	writerAddPos(ctx, len(outByte))
	return
}

func WriteJsonUint64(ctx *Context, i uint64) {
	writerTryRealloc(ctx, 21)
	outByte := strconv.AppendUint(writerGetHeadBuffer(ctx, 0), i, 10)
	writerAddPos(ctx, len(outByte))
	return
}

func WriteJsonInt64AsString(ctx *Context, i int64) {
	writerTryRealloc(ctx, 23)
	WriterWriteByte(ctx, '"')
	outByte := strconv.AppendInt(writerGetHeadBuffer(ctx, 0), i, 10)
	writerAddPos(ctx, len(outByte))
	WriterWriteByte(ctx, '"')
	return
}

func WriteJsonUint64AsString(ctx *Context, i uint64) {
	writerTryRealloc(ctx, 23)
	WriterWriteByte(ctx, '"')
	outByte := strconv.AppendUint(writerGetHeadBuffer(ctx, 0), i, 10)
	writerAddPos(ctx, len(outByte))
	WriterWriteByte(ctx, '"')
	return
}

func WriteJsonFloat64(ctx *Context, f float64) {
	if f == 0 {

		WriterWriteByte(ctx, '0')
		return
	}
	if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, -1) {
		panic("[WriteJsonFloat64] nsqm7mcm7t can not encode special float64 value to json [" +
			strconv.FormatFloat(f, 'g', -1, 64) + "] [" +
			string(writerGetLastString(ctx, 100)) + "]")
	}
	writerTryRealloc(ctx, 32)
	outByte := strconv.AppendFloat(writerGetHeadBuffer(ctx, 0), f, 'g', -1, 64)
	writerAddPos(ctx, len(outByte))
	return
}

func WriteJsonFloat32(ctx *Context, f float32) {
	if f == 0 {

		WriterWriteByte(ctx, '0')
		return
	}
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 1) || math.IsInf(float64(f), -1) {
		panic("[WriteJsonFloat64] vwuca6b4y8 can not encode special float32 value to json [" +
			strconv.FormatFloat(float64(f), 'g', -1, 64) + "] [" +
			string(writerGetLastString(ctx, 100)) + "]")
	}
	writerTryRealloc(ctx, 32)
	outByte := strconv.AppendFloat(writerGetHeadBuffer(ctx, 0), float64(f), 'g', -1, 32)
	writerAddPos(ctx, len(outByte))
	return
}

func WriteJsonByteSlice(ctx *Context, buf []byte) {
	if buf == nil {
		WriterWriteString(ctx, "null")
		return
	}
	encodeSize := base64.StdEncoding.EncodedLen(len(buf))
	writerTryRealloc(ctx, encodeSize+2)
	WriterWriteByte(ctx, '"')
	if len(buf) > 0 {

		base64.StdEncoding.Encode(ctx.writerData[ctx.writerPos:ctx.writerPos+encodeSize], buf)
		ctx.writerPos += encodeSize
	}
	WriterWriteByte(ctx, '"')
}

func WriteIfEmptyStruct(ctx *Context) bool {
	return ctx.writerData[ctx.writerPos-1] == '}' && ctx.writerData[ctx.writerPos-2] == '{'
}

func WritePosAdd(ctx *Context, offset int) {
	ctx.writerPos += offset
}
