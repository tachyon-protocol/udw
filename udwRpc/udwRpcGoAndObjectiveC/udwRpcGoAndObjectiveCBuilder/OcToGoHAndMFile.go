package udwRpcGoAndObjectiveCBuilder

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *builderCtx) ocToGoGenFnHAndMFile(fn *udwGoParser.FuncOrMethodDeclaration) {
	goFnName := "udwGaoc_go_" + fn.GetName()
	fnPrototypeContent := ctx.getOcFunctionPrototype(fn)
	ctx.hFileBuffer.WriteString(fnPrototypeContent + ";\n")
	ctx.mFileBuffer.WriteString(fnPrototypeContent + "{\n")
	if !fn.HasInOrOutParameter() {

		ctx.mFileBuffer.WriteString(`	` + goFnName + `();
	}
`)
		return
	}
	ctx.mFileBuffer.WriteString(`    udwGaocGoBuffer _buf1 = {};
	udwGaocGoBuffer* _buf = &_buf1;
    udwGaocl_c_init(_buf);
`)
	for _, para := range fn.GetInParameter() {
		ctx.genOcMarshal(&ctx.mFileBuffer, para.GetType(), para.GetName())
	}

	ctx.mFileBuffer.WriteString("\t" + goFnName + `(&(_buf->buf),&(_buf->cap));
	_buf->off = 0;
`)
	if len(fn.GetOutParameter()) == 1 {
		para := fn.GetOutParameter()[0]
		ctx.genOcUnmarshal(&ctx.mFileBuffer, para.GetType(), "_udwGaoc_read_"+para.GetName())
		ctx.mFileBuffer.WriteString("	udwGaocl_c_free(_buf);\n")
		ctx.mFileBuffer.WriteString("	return _udwGaoc_read_" + para.GetName() + ";\n")
	} else if len(fn.GetOutParameter()) == 0 {
		ctx.mFileBuffer.WriteString("	udwGaocl_c_free(_buf);\n")
		ctx.mFileBuffer.WriteString("	return ;\n")
	}
	ctx.mFileBuffer.WriteString("}\n")
}
