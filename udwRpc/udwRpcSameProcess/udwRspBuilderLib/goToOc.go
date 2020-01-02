package udwRspBuilderLib

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *GoBuilderCtx) FromGoGenFnGoFile(fn *udwGoParser.FuncOrMethodDeclaration) {
	cFnName := `udwGaoc_c_` + fn.GetName()
	if !fn.HasInOrOutParameter() {
		if ctx.Req.IsNoParameterFromGoDirectCall {
			cFnName = fn.GetName()
		}
		ctx.GoFileBuffer.WriteString(`func ` + fn.GetName() + `(){
	C.` + cFnName + `();
}
`)
		ctx.GoFileHBuffer.WriteString(`void ` + cFnName + `();` + "\n")
		return
	}
	ctx.GoFileContext.AddImportPath("github.com/tachyon-protocol/udw/udwRpc/udwRpcSameProcess/udwRspLib")
	ctx.GoFileContext.AddImportPath("unsafe")
	ctx.GoFileHBuffer.WriteString(ctx.GetFromGoCFnPrototypeContent(cFnName) + ";\n")
	ctx.GoFileBuffer.WriteString(`func ` + fn.GetName() + `(`)
	for _, para := range fn.GetInParameter() {
		typS := ctx.GoFileContext.MustWriteGoTypes(para.GetType())
		ctx.GoFileBuffer.WriteString(para.GetName() + " " + typS + ",")
	}
	ctx.GoFileBuffer.WriteString(`)(`)
	for _, para := range fn.GetOutParameter() {
		typS := ctx.GoFileContext.MustWriteGoTypes(para.GetType())
		ctx.GoFileBuffer.WriteString(para.GetName() + " " + typS + ",")
	}
	ctx.GoFileBuffer.WriteString(`){
	_buf := &udwRspLib.GoBuffer{}
`)
	for _, para := range fn.GetInParameter() {
		ctx.GenGoMarshal(&ctx.GoFileBuffer, para.GetType(), para.GetName())
	}
	ctx.GoFileBuffer.WriteString(`	_goP, _gocap := _buf.ToC()
	_ret:= C.` + cFnName + `((*C.uint8_t)(unsafe.Pointer(_goP)), C.size_t(_gocap))
	_buf.SetFromC(uintptr(unsafe.Pointer(_ret.buf)), int(_ret.cap))
	`)
	if len(fn.GetOutParameter()) == 0 {
		ctx.GoFileBuffer.WriteString(`	_buf.FreeFromGo()
	return
}
`)
	} else if len(fn.GetOutParameter()) == 1 {
		para := fn.GetOutParameter()[0]
		ctx.GenGoUnmarshal(&ctx.GoFileBuffer, para.GetType(), "_out")
		ctx.GoFileBuffer.WriteString(`	_buf.FreeFromGo()
	return _out
}
`)
	} else {
		panic("TODO len(fn.GetOutParameter())>1 fnName:" + ctx.CurrentProcessFnName)
	}
}
