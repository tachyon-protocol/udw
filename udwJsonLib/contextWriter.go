package udwJsonLib

import (
	"github.com/tachyon-protocol/udw/udwStrings"
)

func writerReset(ctx *Context) {
	ctx.writerPos = 0
}

func writerGetTmpString(ctx *Context) string {
	return udwStrings.GetStringFromByteArrayNoAlloc(ctx.writerData[:ctx.writerPos])
}

func WriterWriteByte(ctx *Context, b byte) {
	writerTryRealloc(ctx, 1)
	ctx.writerData[ctx.writerPos] = b
	ctx.writerPos++
}

func WriterWriteByteList(ctx *Context, s []byte) {
	writerTryRealloc(ctx, len(s))
	copy(ctx.writerData[ctx.writerPos:], s)
	ctx.writerPos += len(s)
}

func WriterWriteString(ctx *Context, s string) {
	writerTryRealloc(ctx, len(s))
	copy(ctx.writerData[ctx.writerPos:], s)
	ctx.writerPos += len(s)
}

func writerGetHeadBuffer(ctx *Context, size int) []byte {
	writerTryRealloc(ctx, size)
	return ctx.writerData[ctx.writerPos : ctx.writerPos+size]
}

func writerAddPos(ctx *Context, size int) {
	ctx.writerPos += size
}

func writerTryRealloc(ctx *Context, toAddSize int) {
	needSize := ctx.writerPos + toAddSize
	if len(ctx.writerData) >= needSize {
		return
	}
	needToAllocSize := len(ctx.writerData)*2 + toAddSize
	allocSize := needToAllocSize
	if allocSize <= 64 {
		allocSize = 64

	}

	newBuf := make([]byte, allocSize)
	copy(newBuf, ctx.writerData[:ctx.writerPos])
	ctx.writerData = newBuf
	return
}

func writerGetLastString(ctx *Context, size int) []byte {
	if ctx.writerPos >= size {
		return ctx.writerData[ctx.writerPos-size : ctx.writerPos]
	}
	return ctx.writerData[0:ctx.writerPos]
}
