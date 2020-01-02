package udwRspBuilderLib

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *GoBuilderCtx) ToGoGenFnGoFileForJava(fn *udwGoParser.FuncOrMethodDeclaration) {

	goFnName := "UdwGaoc_go_" + fn.GetName()

	if !fn.HasInOrOutParameter() {

		ctx.GoFileBuffer.WriteString(`func ` + goFnName + `(_pp uintptr,_c uintptr){
	` + fn.GetName() + `()
}
`)
		return
	}
	ctx.GoFileContext.AddImportPath("github.com/tachyon-protocol/udw/udwRpc/udwRpcSameProcess/udwRspLib")
	ctx.GoFileContext.AddImportPath("unsafe")
	ctx.GoFileBuffer.WriteString(`func ` + goFnName + `(_pp uintptr,_c uintptr){
	_p:=*(*uintptr)(unsafe.Pointer(_pp))
	_buf:=udwRspLib.NewGoBufferFromC(uintptr(unsafe.Pointer(_p)),int(*(*int)(unsafe.Pointer(_c))))
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
	*(*uintptr)(unsafe.Pointer(_pp)) = uintptr(unsafe.Pointer(_goP))
	*(*int)(unsafe.Pointer(_c)) = (int)(_gocap)
}
	`)
	} else {
		panic("TODO len(outparameter)>0 " + ctx.CurrentProcessFnName)
	}

}
