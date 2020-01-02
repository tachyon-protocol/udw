package udwRpcGoAndObjectiveCBuilder

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *builderCtx) getOcType(typ udwGoParser.Type) string {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		return "bool"
	case udwGoParser.String:
		return "NSString*"
	case udwGoParser.Int:
		return "size_t"
	case udwGoParser.Int64:
		return "int64_t"
	case udwGoParser.Float64:
		return "double"
	case udwGoParser.Float32:
		return "float"
	case udwGoParser.Slice:
		t := typ.(*udwGoParser.SliceType)
		switch t.Elem.GetKind() {
		case udwGoParser.Uint8:
			return "NSData*"
		default:
			return "NSArray<" + ctx.getOcTypeInNSArray(t.Elem) + ">*"

		}
	case udwGoParser.Map:
		t := typ.(*udwGoParser.MapType)
		return "NSDictionary<" + ctx.getOcTypeInNSArray(t.Key) + "," + ctx.getOcTypeInNSArray(t.Value) + ">*"
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" {
			if t.Name == "Time" {

				return "NSDate*"
			} else if t.Name == "Duration" {

				return "NSTimeInterval"
			}
		}
		if t.GetUnderType().GetKind() == udwGoParser.Struct {

			return t.Name + "*"
		}
		return ctx.getOcType(t.GetUnderType())
	case udwGoParser.Ptr:
		t := typ.(*udwGoParser.PointerType)
		t1 := t.Elem
		if t1.GetKind() == udwGoParser.Named {
			t2 := t1.(*udwGoParser.NamedType)
			t3 := t2.GetUnderType()
			if t3.GetKind() == udwGoParser.Struct {
				return t2.Name + "*"
			}
		}
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	default:
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	}
}

func (ctx *builderCtx) getOcTypeInNSArray(typ udwGoParser.Type) string {
	switch typ.GetKind() {
	case udwGoParser.Bool, udwGoParser.Int, udwGoParser.Int64, udwGoParser.Float64, udwGoParser.Float32:
		return "NSNumber* /*" + typ.GetKind().String() + "*/"
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" && t.Name == "Duration" {
			return "NSNumber* /*time.Duration*/"
		}
	}
	return ctx.getOcType(typ)
}

func (ctx *builderCtx) genOcMarshal(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString("\t" + `udwGaocl_c_writeBool(_buf,` + varName + `);` + "\n")
	case udwGoParser.String:
		_buf.WriteString("\t" + `udwGaocl_c_writeString(_buf,` + varName + `);` + "\n")
	case udwGoParser.Int:
		_buf.WriteString("\t" + `udwGaocl_c_writeInt(_buf,` + varName + `);` + "\n")
	case udwGoParser.Int64:
		_buf.WriteString("\t" + `udwGaocl_c_writeInt64(_buf,` + varName + `);` + "\n")
	case udwGoParser.Float64:
		_buf.WriteString("\t" + `udwGaocl_c_writeFloat64(_buf,` + varName + `);` + "\n")
	case udwGoParser.Float32:
		_buf.WriteString("\t" + `udwGaocl_c_writeFloat32(_buf,` + varName + `);` + "\n")
	case udwGoParser.Slice:
		t := typ.(*udwGoParser.SliceType)
		switch t.Elem.GetKind() {
		case udwGoParser.Uint8:
			_buf.WriteString("\t" + `udwGaocl_c_writeByteSlice(_buf,` + varName + `);` + "\n")
		default:
			var1Name := ctx.getNextVarString()
			_buf.WriteString("\t" + `udwGaocl_c_writeInt(_buf,(size_t)[` + varName + ` count]);
    for (` + ctx.getOcTypeInNSArray(t.Elem) + ` ` + var1Name + ` in ` + varName + `){` + "\n")
			ctx.genOcMarshalInNSArray(_buf, t.Elem, var1Name)
			_buf.WriteString(`	}
`)
		}
	case udwGoParser.Map:
		t := typ.(*udwGoParser.MapType)
		var1Name := ctx.getNextVarString()
		_buf.WriteString("\t" + `udwGaocl_c_writeInt(_buf,(size_t)[` + varName + ` count]);
    for (` + ctx.getOcTypeInNSArray(t.Key) + ` ` + var1Name + ` in ` + varName + `){` + "\n")
		ctx.genOcMarshalInNSArray(_buf, t.Key, var1Name)
		var2Name := ctx.getNextVarString()
		_buf.WriteString(`		` + ctx.getOcTypeInNSArray(t.Value) + " " + var2Name + "=" + varName + "[" + var1Name + "];\n")
		ctx.genOcMarshalInNSArray(_buf, t.Value, var2Name)
		_buf.WriteString(`	}
`)
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" {
			if t.Name == "Time" {
				var1Name := ctx.getNextVarString()
				_buf.WriteString(`double ` + var1Name + ` = (double)[` + varName + ` timeIntervalSince1970];` + "\n")
				_buf.WriteString(`udwGaocl_c_writeInt64(_buf,(int64_t)` + var1Name + `);` + "\n")
				_buf.WriteString(`udwGaocl_c_writeInt64(_buf,(int64_t)((` + var1Name + `-(double)((int64_t)` + var1Name + `))*1e9) );` + "\n")
				return
			} else if t.Name == "Duration" {

				_buf.WriteString("\t" + `udwGaocl_c_writeFloat64(_buf,(double)` + varName + `);` + "\n")
				return
			}
		}
		if t.GetUnderType().GetKind() == udwGoParser.Struct {
			structType := t.GetUnderType().(*udwGoParser.StructType)
			ctx.genOcMarshalNamedStruct(_buf, t.Name, structType, varName)
			return
		}
		ctx.genOcMarshal(_buf, t.GetUnderType(), varName)
	case udwGoParser.Ptr:
		t := typ.(*udwGoParser.PointerType)
		t1 := t.Elem
		if t1.GetKind() == udwGoParser.Named {
			t2 := t1.(*udwGoParser.NamedType)
			t3 := t2.GetUnderType()
			if t3.GetKind() == udwGoParser.Struct {
				structType := t3.(*udwGoParser.StructType)
				_buf.WriteString("\t" + `udwGaocl_c_writeBool(_buf,` + varName + `!=nil);` + "\n")
				_buf.WriteString(`	if (` + varName + `!=nil){` + "\n")
				ctx.genOcMarshalNamedStruct(_buf, t2.Name, structType, varName)
				_buf.WriteString("	}\n")
				return
			}
		}
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	default:
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	}
}

func (ctx *builderCtx) genOcMarshalInNSArray(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString("\t" + `udwGaocl_c_writeBool(_buf,(bool)[` + varName + ` boolValue]);` + "\n")
		return
	case udwGoParser.Int:
		_buf.WriteString("\t" + `udwGaocl_c_writeInt(_buf,(size_t)[` + varName + ` integerValue]);` + "\n")
		return
	case udwGoParser.Int64:
		_buf.WriteString("\t" + `udwGaocl_c_writeInt64(_buf,(int64_t)[` + varName + ` longLongValue]);` + "\n")
		return
	case udwGoParser.Float64:
		_buf.WriteString("\t" + `udwGaocl_c_writeFloat64(_buf,(double)[` + varName + ` doubleValue]);` + "\n")
		return
	case udwGoParser.Float32:
		_buf.WriteString("\t" + `udwGaocl_c_writeFloat32(_buf,(float)[` + varName + ` floatValue]);` + "\n")
		return
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" && t.Name == "Duration" {

			_buf.WriteString("\t" + `udwGaocl_c_writeFloat64(_buf,(double)[` + varName + ` doubleValue]);` + "\n")
			return
		}
	}
	ctx.genOcMarshal(_buf, typ, varName)
	return
}

func (ctx *builderCtx) genOcMarshalNamedStruct(_buf *bytes.Buffer, name string, structType *udwGoParser.StructType, varName string) {
	marshalFnName := `_udwGaoc_Marshal_` + name
	_buf.WriteString(marshalFnName + `(_buf,` + varName + `)` + ";\n")
	if ctx.seenMarhshalNameMap[name] {
		return
	}
	ctx.seenMarhshalNameMap[name] = true
	thisFuncBuffer := &bytes.Buffer{}
	thisFuncBuffer.WriteString("void " + marshalFnName + "(udwGaocGoBuffer* _buf," + name + "* _var){\n")
	for _, f := range structType.Field {
		ctx.genOcMarshal(thisFuncBuffer, f.Elem, "_var."+f.Name)
	}
	thisFuncBuffer.WriteString("}\n")
	ctx.mFuncFileBuffer.Write(thisFuncBuffer.Bytes())
	return
}

func (ctx *builderCtx) genOcUnmarshal(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString("\t" + `bool ` + varName + `=udwGaocl_c_readBool(_buf);` + "\n")
		return
	case udwGoParser.String:
		_buf.WriteString("\t" + `NSString* ` + varName + `=udwGaocl_c_readString(_buf);` + "\n")
		return
	case udwGoParser.Int:
		_buf.WriteString("\t" + `size_t ` + varName + `=udwGaocl_c_readInt(_buf);` + "\n")
		return
	case udwGoParser.Int64:
		_buf.WriteString("\t" + `size_t ` + varName + `=udwGaocl_c_readInt64(_buf);` + "\n")
		return
	case udwGoParser.Float64:
		_buf.WriteString("\t" + `double ` + varName + `=udwGaocl_c_readFloat64(_buf);` + "\n")
		return
	case udwGoParser.Float32:
		_buf.WriteString("\t" + `float ` + varName + `=udwGaocl_c_readFloat32(_buf);` + "\n")
		return
	case udwGoParser.Slice:
		t := typ.(*udwGoParser.SliceType)
		switch t.Elem.GetKind() {
		case udwGoParser.Uint8:
			_buf.WriteString("\t" + `NSData* ` + varName + `=udwGaocl_c_readByteSlice(_buf);` + "\n")
			return
		default:
			var1Name := ctx.getNextVarString()
			var2Name := ctx.getNextVarString()
			_buf.WriteString("\t" + `size_t ` + var1Name + `_len = udwGaocl_c_readInt(_buf);
    NSMutableArray<` + ctx.getOcTypeInNSArray(t.Elem) + `>* ` + varName + ` = [[NSMutableArray alloc] initWithCapacity:` + var1Name + `_len];
    for (int i=0;i<` + var1Name + `_len;i++){` + "\n")
			ctx.genOcUnmarshalInNSArray(_buf, t.Elem, var2Name)
			_buf.WriteString(`		[` + varName + " addObject:" + var2Name + "];\n")
			_buf.WriteString(`	}
`)
		}
		return
	case udwGoParser.Map:
		t := typ.(*udwGoParser.MapType)
		var1Name := ctx.getNextVarString()
		var2Name := ctx.getNextVarString()
		var3Name := ctx.getNextVarString()
		_buf.WriteString("\t" + `size_t ` + var1Name + ` = udwGaocl_c_readInt(_buf);
    ` + "NSMutableDictionary<" + ctx.getOcTypeInNSArray(t.Key) + "," + ctx.getOcTypeInNSArray(t.Value) + ">* " + varName +
			` = [[NSMutableDictionary alloc] initWithCapacity:` + var1Name + `];
    for (int i=0;i<` + var1Name + `;i++){` + "\n")
		ctx.genOcUnmarshalInNSArray(_buf, t.Key, var2Name)
		ctx.genOcUnmarshalInNSArray(_buf, t.Value, var3Name)
		_buf.WriteString(`		` + varName + `[` + var2Name + `]=` + var3Name + ";\n")
		_buf.WriteString(`	}
`)
		return
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" {
			if t.Name == "Time" {
				var1Name := ctx.getNextVarString()
				var2Name := ctx.getNextVarString()
				_buf.WriteString("\t" + `int64_t ` + var1Name + `=udwGaocl_c_readInt64(_buf);` + "\n")
				_buf.WriteString("\t" + `int64_t ` + var2Name + `=udwGaocl_c_readInt64(_buf);` + "\n")

				_buf.WriteString("\t" + "NSDate* " + varName +
					`=[NSDate dateWithTimeIntervalSince1970:(NSTimeInterval)((double)` + var1Name + `+((double)` + var2Name + `)/1e9)];` + "\n")
				return
			} else if t.Name == "Duration" {

				_buf.WriteString("\t" + "NSTimeInterval " + varName + `=(NSTimeInterval)udwGaocl_c_readFloat64(_buf);` + "\n")
				return
			}
		}
		if t.GetUnderType().GetKind() == udwGoParser.Struct {
			_buf.WriteString("	" + t.Name + "* " + varName + "=[" + t.Name + " new];\n")
			structType := t.GetUnderType().(*udwGoParser.StructType)
			ctx.genOcUnmarshalNamedStruct(_buf, t.Name, structType, varName)
			return
		}
		ctx.genOcUnmarshal(_buf, t.GetUnderType(), varName)
		return
	case udwGoParser.Ptr:
		t := typ.(*udwGoParser.PointerType)
		t1 := t.Elem
		if t1.GetKind() == udwGoParser.Named {
			t2 := t1.(*udwGoParser.NamedType)
			t3 := t2.GetUnderType()
			if t3.GetKind() == udwGoParser.Struct {
				var1Name := ctx.getNextVarString()
				_buf.WriteString("	bool " + var1Name + ` = udwGaocl_c_readBool(_buf);
	` + t2.Name + "* " + varName + ` = nil;
	if (` + var1Name + `){
`)
				var2Name := ctx.getNextVarString()
				ctx.genOcUnmarshal(_buf, t1, var2Name)
				_buf.WriteString(`		` + varName + ` = ` + var2Name + `;
	}
`)
				return
			}
		}
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	default:
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.getCurrentProcessFnName()))
	}
}

func (ctx *builderCtx) genOcUnmarshalInNSArray(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithBool:(BOOL)udwGaocl_c_readBool(_buf)];` + "\n")
		return
	case udwGoParser.Int:
		_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithInteger:(NSInteger)udwGaocl_c_readInt(_buf)];` + "\n")
		return
	case udwGoParser.Int64:
		_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithLongLong:(long long)udwGaocl_c_readInt64(_buf)];` + "\n")
		return
	case udwGoParser.Float64:
		_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithDouble:(double)udwGaocl_c_readFloat64(_buf)];` + "\n")
		return
	case udwGoParser.Float32:
		_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithFloat:(float)udwGaocl_c_readFloat32(_buf)];` + "\n")
		return
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" && t.Name == "Duration" {

			_buf.WriteString("\t" + `NSNumber* ` + varName + `=[NSNumber numberWithDouble:(double)udwGaocl_c_readFloat64(_buf)];` + "\n")
			return
		}
	}
	ctx.genOcUnmarshal(_buf, typ, varName)
	return
}

func (ctx *builderCtx) genOcUnmarshalNamedStruct(_buf *bytes.Buffer, name string, structType *udwGoParser.StructType, varName string) {
	marshalFnName := `_udwGaoc_Unmarshal_` + name
	_buf.WriteString(`		` + marshalFnName + `(_buf,` + varName + `)` + ";\n")
	if ctx.seenUnmarshalNameMap[name] {
		return
	}
	ctx.seenUnmarshalNameMap[name] = true
	thisFuncBuffer := &bytes.Buffer{}
	thisFuncBuffer.WriteString("void " + marshalFnName + "(udwGaocGoBuffer* _buf," + name + "* _var){\n")
	for _, f := range structType.Field {
		var1Name := ctx.getNextVarString()
		ctx.genOcUnmarshal(thisFuncBuffer, f.Elem, var1Name)
		thisFuncBuffer.WriteString("	_var." + f.Name + "=" + var1Name + ";\n")
	}
	thisFuncBuffer.WriteString("}\n")
	ctx.mFuncFileBuffer.Write(thisFuncBuffer.Bytes())
	return
}

func (ctx *builderCtx) getOcFunctionPrototype(fn *udwGoParser.FuncOrMethodDeclaration) string {
	_buf := bytes.Buffer{}
	returnTypeName := ""
	if len(fn.GetOutParameter()) == 0 {
		returnTypeName = "void"
	} else if len(fn.GetOutParameter()) == 1 {
		returnTypeName = ctx.getOcType(fn.GetOutParameter()[0].GetType())
	} else {
		panic("TODO len(fn.GetOutParameter())>1 fnName:" + ctx.getCurrentProcessFnName())
	}
	_buf.WriteString(returnTypeName + " " + fn.GetName() + "(")
	for i, para := range fn.GetInParameter() {
		_buf.WriteString(ctx.getOcType(para.GetType()) + " " + para.GetName())
		if i != len(fn.GetInParameter())-1 {
			_buf.WriteString(",")
		}
	}
	_buf.WriteString(")")
	return _buf.String()
}
