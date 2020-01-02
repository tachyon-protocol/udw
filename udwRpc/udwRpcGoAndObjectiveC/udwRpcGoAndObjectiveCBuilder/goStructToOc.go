package udwRpcGoAndObjectiveCBuilder

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *builderCtx) genGoStructToOc(visitFn func(f func(fnDef *udwGoParser.FuncOrMethodDeclaration))) {
	nameMap := map[string]bool{}
	visitFn(func(fn *udwGoParser.FuncOrMethodDeclaration) {
		udwGoParser.VisitNamedStruct(udwGoParser.VisitNamedStructRequest{
			Type:                fn,
			VisitNotExportField: true,
			Fn: func(namedTyp *udwGoParser.NamedType) {
				if nameMap[namedTyp.Name] {
					return
				}
				nameMap[namedTyp.Name] = true
				if namedTyp.PkgImportPath == "time" {
					return
				}
				ctx.hFileBuffer.WriteString("@class " + namedTyp.Name + ";\n")
			},
		})
	})
	nameMap = map[string]bool{}
	visitFn(func(fn *udwGoParser.FuncOrMethodDeclaration) {
		udwGoParser.VisitNamedStruct(udwGoParser.VisitNamedStructRequest{
			Type:                fn,
			VisitNotExportField: true,
			Fn: func(namedTyp *udwGoParser.NamedType) {
				if nameMap[namedTyp.Name] {
					return
				}
				nameMap[namedTyp.Name] = true
				if namedTyp.PkgImportPath == "time" {
					return
				}
				ctx.hFileBuffer.WriteString("@interface " + namedTyp.Name + ":NSObject\n")
				for _, f := range namedTyp.GetUnderType().(*udwGoParser.StructType).Field {
					ctx.hFileBuffer.WriteString("@property " + ctx.getOcType(f.Elem) + " " + f.Name + ";\n")
				}
				ctx.hFileBuffer.WriteString(`@end

`)
				ctx.mFileBuffer.WriteString(`@implementation ` + namedTyp.Name + `
@end
`)
			},
		})
	})
}
