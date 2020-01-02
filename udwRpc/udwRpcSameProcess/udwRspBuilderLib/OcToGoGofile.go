package udwRspBuilderLib

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *GoBuilderCtx) ToGoGenFnGoFile(fn *udwGoParser.FuncOrMethodDeclaration) {
	fnName := fn.GetName()
	goFnName := "udwGaoc_go_" + fn.GetName()
	ctx.GoFileBuffer.WriteString(`//export ` + goFnName + "\n")
	if !fn.HasInOrOutParameter() {

		ctx.GoFileBuffer.WriteString(`func ` + goFnName + `(){
	` + fnName + `()
}
`)
		return
	}
	ctx.GoFileContext.AddImportPath("github.com/tachyon-protocol/udw/udwRpc/udwRpcSameProcess/udwRspLib")
	ctx.GoFileContext.AddImportPath("unsafe")
	ctx.GoFileBuffer.WriteString(`func ` + goFnName + `(_p **C.uint8_t,_c *C.size_t){
	_buf:=udwRspLib.NewGoBufferFromC(uintptr(unsafe.Pointer(*_p)),int(*_c))
`)
	for _, para := range fn.GetInParameter() {
		ctx.GenGoUnmarshal(&ctx.GoFileBuffer, para.GetType(), "_udwGaoc_read_"+para.GetName())
	}
	writeCallInParaFn := func() {
		for _, para := range fn.GetInParameter() {
			ctx.GoFileBuffer.WriteString("_udwGaoc_read_" + para.GetName() + ",")
		}
	}
	if len(fn.GetOutParameter()) == 0 {
		ctx.GoFileBuffer.WriteString(fn.GetName() + "(")
		writeCallInParaFn()
		ctx.GoFileBuffer.WriteString(`)
}
`)

	} else if len(fn.GetOutParameter()) == 1 {
		para := fn.GetOutParameter()[0]
		ctx.GoFileBuffer.WriteString("_ret:=" + fn.GetName() + "(")
		writeCallInParaFn()
		ctx.GoFileBuffer.WriteString(`)
	_buf.ResetToWrite()
`)
		ctx.GenGoMarshal(&ctx.GoFileBuffer, para.GetType(), "_ret")
		ctx.GoFileBuffer.WriteString(`	_goP,_gocap := _buf.ToC()
	*_p = (*C.uint8_t)(unsafe.Pointer(_goP))
	*_c = (C.size_t)(_gocap)
}
	`)
	} else {
		panic("TODO len(outparameter)>0 " + ctx.CurrentProcessFnName)
	}

}
