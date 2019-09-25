package udwJsonLib

import "sync"

var gJsonLibPool = sync.Pool{}

func PoolGetContextWithWriteBuf() *Context {
	var ctx *Context
	ctxI := gJsonLibPool.Get()
	if ctxI != nil {
		ctx = ctxI.(*Context)
		ctx.readerPos = 0
		ctx.writerPos = 0
	} else {
		ctx = NewContext()
	}
	return ctx
}

func PoolGetContextForRead(readData []byte) *Context {
	ctx := PoolGetContextWithWriteBuf()
	ctx.readerData = readData
	return ctx
}

func PoolPutContextWithWriteBuf(ctx *Context) {
	if ctx == nil {
		return
	}
	ctx.readerData = nil
	gJsonLibPool.Put(ctx)
}
