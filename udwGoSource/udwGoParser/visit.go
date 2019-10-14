package udwGoParser

func VisitNamedStruct(req VisitNamedStructRequest) {
	ctx := visitNamedStructCtx{
		visitedMap: map[string]bool{},
		req:        req,
	}
	if req.Type == nil {
		panic("[VisitNamedStruct] req.Type==nil")
	}
	visitNamedStructL1(ctx, req.Type)
}

type VisitNamedStructRequest struct {
	Type                Type
	Fn                  func(namedTyp *NamedType)
	VisitNotExportField bool
}

type visitNamedStructCtx struct {
	req        VisitNamedStructRequest
	visitedMap map[string]bool
}

func visitNamedStructL1(ctx visitNamedStructCtx, typ Type) {
	switch typ.GetKind() {
	case Named:
		t := typ.(*NamedType)
		ts := t.String()

		if ctx.visitedMap[ts] {
			return
		}
		ctx.visitedMap[ts] = true
		if t.GetUnderType().GetKind() == Struct {
			structType := t.GetUnderType().(*StructType)
			ctx.req.Fn(t)
			for _, f := range structType.Field {
				if ctx.req.VisitNotExportField == false && IsNameGoExport(f.Name) == false {
					continue
				}
				visitNamedStructL1(ctx, f.Elem)
			}
		}
	case DefinedFunc, Method:
		thisFn := typ.(*FuncOrMethodDeclaration)
		for _, para := range thisFn.GetInParameter() {
			visitNamedStructL1(ctx, para.GetType())
		}
		for _, para := range thisFn.GetOutParameter() {
			visitNamedStructL1(ctx, para.GetType())
		}
	case Slice:
		t := typ.(*SliceType)
		visitNamedStructL1(ctx, t.Elem)
	case Array:
		t := typ.(*ArrayType)
		visitNamedStructL1(ctx, t.Elem)
	case Map:
		t := typ.(*MapType)
		visitNamedStructL1(ctx, t.Key)
		visitNamedStructL1(ctx, t.Value)
	case Ptr:
		t := typ.(*PointerType)
		visitNamedStructL1(ctx, t.Elem)
	}
}
