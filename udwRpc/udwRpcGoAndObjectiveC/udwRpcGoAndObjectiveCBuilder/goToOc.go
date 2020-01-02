package udwRpcGoAndObjectiveCBuilder

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *builderCtx) goToOcGenFnHAndMFile(fn *udwGoParser.FuncOrMethodDeclaration) {
	if !fn.HasInOrOutParameter() {

		return
	}
	cFnName := `udwGaoc_c_` + fn.GetName()

	ctx.mFileBuffer.WriteString(ctx.getOcFunctionPrototype(fn) + ";\n")
	ctx.mFileBuffer.WriteString(ctx.goBuilderCtx.GetFromGoCFnPrototypeContent(cFnName) + `{
    udwGaocGoBuffer _buf1 = {};
    udwGaocGoBuffer* _buf = &_buf1;
    _buf->buf = _pIn;
    _buf->cap = _cIn;
`)
	callExprContentBuf := bytes.Buffer{}
	for i, para := range fn.GetInParameter() {
		ctx.genOcUnmarshal(&ctx.mFileBuffer, para.GetType(), "_udwGaoc_read_"+para.GetName())
		callExprContentBuf.WriteString("_udwGaoc_read_" + para.GetName())
		if i != len(fn.GetInParameter())-1 {
			callExprContentBuf.WriteString(",")
		}
	}
	if len(fn.GetOutParameter()) == 0 {
		ctx.mFileBuffer.WriteString("	" + fn.GetName() + "(")
		ctx.mFileBuffer.Write(callExprContentBuf.Bytes())
		ctx.mFileBuffer.WriteString(`);
	udwGaoc_c_PAndC _out = {_buf->buf,_buf->cap};
	return _out;
}
`)
	} else if len(fn.GetOutParameter()) == 1 {
		para := fn.GetOutParameter()[0]
		ctx.mFileBuffer.WriteString("	" + ctx.getOcType(para.GetType()) + " _ret = " + fn.GetName() + "(")
		ctx.mFileBuffer.Write(callExprContentBuf.Bytes())
		ctx.mFileBuffer.WriteString(`);
	_buf->off = 0;
`)
		ctx.genOcMarshal(&ctx.mFileBuffer, para.GetType(), "_ret")
		ctx.mFileBuffer.WriteString(`	udwGaoc_c_PAndC _out = {_buf->buf,_buf->cap};
	return _out;
}
`)
	} else {
		panic("TODO len(fn.GetOutParameter())>1 fnName:" + ctx.getCurrentProcessFnName())
	}
}
