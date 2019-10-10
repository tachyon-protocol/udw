package udwJsonLib

import (
	"encoding/base64"
)

func MustReadJsonByteSlice(ctx *Context) (outB []byte) {
	b := ReaderReadByte(ctx)
	switch b {
	case 'n':
		ReaderReadBack(ctx, 1)
		MustReadJsonNull(ctx)
		return nil
	case '[':
		outB = []byte{}
		for {
			b = ReaderReadByte(ctx)
			if b == ',' || b == ' ' || b == '\t' || b == '\n' || b == '\r' {
				continue
			} else if b == ']' {
				break
			}
			ReaderReadBack(ctx, 1)
			_var1 := byte(ReadJsonUint64(ctx))
			outB = append(outB, _var1)
		}
		return outB
	case '"':

		startPos := ctx.readerPos
		endPos := startPos
		for {
			b := ReaderReadByte(ctx)
			if b == '"' {
				endPos = ctx.readerPos - 1
				break
			}
		}
		if endPos == startPos {
			return []byte{}
		}
		dbuf := make([]byte, base64.StdEncoding.DecodedLen(endPos-startPos))
		n, err := base64.StdEncoding.Decode(dbuf, ctx.readerData[startPos:endPos])
		if err != nil {
			panic(err)
		}
		return dbuf[:n]
	default:
		panic(`MustReadJsonByteSlice need a [ or null or " but get [` + string(b) + "]")
	}
}
