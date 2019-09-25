package udwJsonLib

type Context struct {
	readerData []byte
	readerPos  int
	writerData []byte
	writerPos  int
}

func NewContext() *Context {
	return &Context{}
}

func NewContextFromBuffer(b []byte) *Context {
	ctx := &Context{
		readerData: b,
	}
	return ctx
}

func NewContextFromWriteBuffer(b []byte) *Context {
	ctx := &Context{
		writerData: b,
	}
	return ctx
}

func (ctx *Context) WriterReset() {
	writerReset(ctx)
}

func (ctx *Context) WriterBytes() []byte {
	return ctx.writerData[:ctx.writerPos]
}

func (ctx *Context) SetReader(data []byte) {
	ReaderSetData(ctx, data)
	writerReset(ctx)
}

func (ctx *Context) CopyWriterDataToReaderData() {
	if len(ctx.readerData) < len(ctx.writerData) {
		ctx.readerData = make([]byte, len(ctx.writerData))
	}
	copy(ctx.readerData, ctx.writerData)
	ctx.readerPos = 0
}

func (ctx *Context) ReaderReset() {
	ctx.readerPos = 0
}
