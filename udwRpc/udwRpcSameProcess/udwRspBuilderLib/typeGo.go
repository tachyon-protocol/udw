package udwRspBuilderLib

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
)

func (ctx *GoBuilderCtx) GenGoMarshal(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString(`_buf.WriteBool(` + varName + `)` + "\n")
	case udwGoParser.String:
		if ctx.Req.IsStringUTF16 {
			_buf.WriteString(`_buf.WriteStringUTF16(` + varName + `)` + "\n")
		} else {
			_buf.WriteString(`_buf.WriteStringUTF8(` + varName + `)` + "\n")
		}
	case udwGoParser.Int:
		_buf.WriteString(`_buf.WriteInt(` + varName + `)` + "\n")
	case udwGoParser.Uint8:
		_buf.WriteString(`_buf.WriteUint8(` + varName + `)` + "\n")
	case udwGoParser.Uint16:
		_buf.WriteString(`_buf.WriteUint16(` + varName + `)` + "\n")
	case udwGoParser.Uint32:
		_buf.WriteString(`_buf.WriteUint32(` + varName + `)` + "\n")
	case udwGoParser.Uint64:
		_buf.WriteString(`_buf.WriteUint64(` + varName + `)` + "\n")
	case udwGoParser.Int64:
		_buf.WriteString(`_buf.WriteInt64(` + varName + `)` + "\n")
	case udwGoParser.Float32:
		_buf.WriteString(`_buf.WriteFloat32(` + varName + `)` + "\n")
	case udwGoParser.Float64:
		_buf.WriteString(`_buf.WriteFloat64(` + varName + `)` + "\n")
	case udwGoParser.Slice:
		t := typ.(*udwGoParser.SliceType)
		if t.Elem.GetKind() == udwGoParser.Uint8 {
			_buf.WriteString(`_buf.WriteByteSlice(` + varName + `)` + "\n")
			return
		}
		_buf.WriteString(`_buf.WriteInt(len(` + varName + `))
	for _,s:=range ` + varName + `{` + "\n")
		ctx.GenGoMarshal(_buf, t.Elem, "s")
		_buf.WriteString(`	}
`)
	case udwGoParser.Map:
		t := typ.(*udwGoParser.MapType)
		_buf.WriteString(`_buf.WriteInt(len(` + varName + `))
	for k,v:=range ` + varName + `{` + "\n")
		ctx.GenGoMarshal(_buf, t.Key, "k")
		ctx.GenGoMarshal(_buf, t.Value, "v")
		_buf.WriteString(`	}
`)
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" {
			if t.Name == "Time" {

				_buf.WriteString(`_buf.WriteInt64(int64(` + varName + `.Unix()))` + "\n")
				_buf.WriteString(`_buf.WriteInt64(int64(` + varName + `.Nanosecond()))` + "\n")
				return
			} else if t.Name == "Duration" {

				_buf.WriteString(`_buf.WriteFloat64(` + varName + `.Seconds())` + "\n")
				return
			}
		}
		if t.GetUnderType().GetKind() == udwGoParser.Struct {
			marshalFnName := `_udwRsp_Marshal_` + t.Name
			_buf.WriteString(marshalFnName + `(_buf,&` + varName + `)` + "\n")
			ts := t.String()
			if ctx.seenMarshalNamedStructMap[ts] {
				return
			}
			ctx.seenMarshalNamedStructMap[ts] = true
			structType := t.GetUnderType().(*udwGoParser.StructType)
			thisFuncBuffer := &bytes.Buffer{}
			thisFuncBuffer.WriteString("func " + marshalFnName + "(_buf *udwRspLib.GoBuffer,_var *" + ctx.GoFileContext.MustWriteGoTypes(typ) + "){\n")
			for _, f := range structType.Field {
				ctx.GenGoMarshal(thisFuncBuffer, f.Elem, "_var."+f.Name)
			}
			thisFuncBuffer.WriteString("}\n")
			ctx.GoFileFuncBuffer.Write(thisFuncBuffer.Bytes())
			return
		}
		var1Name := ctx.GetNextVarString()
		_buf.WriteString("	" + var1Name + ":=" + ctx.GoFileContext.MustWriteGoTypes(t.GetUnderType()) + "(" + varName + ")\n")
		ctx.GenGoMarshal(_buf, t.GetUnderType(), var1Name)
	case udwGoParser.Ptr:

		t := typ.(*udwGoParser.PointerType)
		_buf.WriteString(`	_buf.WriteBool(` + varName + `!=nil);
	if ` + varName + `!=nil{
`)
		var1Name := ctx.GetNextVarString()
		_buf.WriteString("	" + var1Name + ":=*" + varName + "\n")
		ctx.GenGoMarshal(_buf, t.Elem, var1Name)
		_buf.WriteString(`}
`)

	default:
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.CurrentProcessFnName))
	}
}

func (ctx *GoBuilderCtx) GenGoUnmarshal(_buf *bytes.Buffer, typ udwGoParser.Type, varName string) {
	switch typ.GetKind() {
	case udwGoParser.Bool:
		_buf.WriteString(varName + `:=_buf.ReadBool()` + "\n")
	case udwGoParser.String:
		if ctx.Req.IsStringUTF16 {
			_buf.WriteString(varName + `:=_buf.ReadStringUTF16()` + "\n")
		} else {
			_buf.WriteString(varName + `:=_buf.ReadStringUTF8()` + "\n")
		}
	case udwGoParser.Int:
		_buf.WriteString(varName + `:=_buf.ReadInt()` + "\n")
	case udwGoParser.Int64:
		_buf.WriteString(varName + `:=_buf.ReadInt64()` + "\n")
	case udwGoParser.Float32:
		_buf.WriteString(varName + `:=_buf.ReadFloat32()` + "\n")
	case udwGoParser.Float64:
		_buf.WriteString(varName + `:=_buf.ReadFloat64()` + "\n")
	case udwGoParser.Slice:
		t := typ.(*udwGoParser.SliceType)
		if t.Elem.GetKind() == udwGoParser.Uint8 {
			_buf.WriteString(varName + `:=_buf.ReadByteSlice()` + "\n")
			return
		}
		goTypeString := ctx.GoFileContext.MustWriteGoTypes(typ)
		var1Name := ctx.GetNextVarString() + "_sLen"
		var2Name := ctx.GetNextVarString()

		_buf.WriteString(var1Name + `:=_buf.ReadInt()
	` + varName + `:=make(` + goTypeString + `,0,` + var1Name + `)
	for i:=0;i<` + var1Name + `;i++{` + "\n")
		ctx.GenGoUnmarshal(_buf, t.Elem, var2Name)
		_buf.WriteString(`		` + varName + ` = append(` + varName + `,` + var2Name + `)
	}
	`)
	case udwGoParser.Map:
		t := typ.(*udwGoParser.MapType)
		goTypeString := ctx.GoFileContext.MustWriteGoTypes(typ)
		var1Name := ctx.GetNextVarString() + "_sLen"
		varkName := ctx.GetNextVarString()
		varvName := ctx.GetNextVarString()

		_buf.WriteString(var1Name + `:=_buf.ReadInt()
	` + varName + `:=make(` + goTypeString + `,` + var1Name + `)
	for i:=0;i<` + var1Name + `;i++{` + "\n")
		ctx.GenGoUnmarshal(_buf, t.Key, varkName)
		ctx.GenGoUnmarshal(_buf, t.Value, varvName)
		_buf.WriteString(`		` + varName + `[` + varkName + `]=` + varvName + `
	}
	`)
	case udwGoParser.Named:
		t := typ.(*udwGoParser.NamedType)
		if t.PkgImportPath == "time" {
			if t.Name == "Time" {
				ctx.GoFileContext.AddImportPath("time")
				var1Name := ctx.GetNextVarString()
				var2Name := ctx.GetNextVarString()
				_buf.WriteString(var1Name + `:=_buf.ReadInt64()
				` + var2Name + `:=_buf.ReadInt64()
`)

				_buf.WriteString(varName + `:=time.Unix(` + var1Name + `,` + var2Name + `).In(time.UTC)` + "\n")
				return
			} else if t.Name == "Duration" {
				ctx.GoFileContext.AddImportPath("time")

				_buf.WriteString(varName + `:=time.Duration(_buf.ReadFloat64()*1e9)` + "\n")
				return
			}
		}
		if t.GetUnderType().GetKind() == udwGoParser.Struct {
			marshalFnName := `_udwRsp_Unmarshal_` + t.Name
			_buf.WriteString("	" + varName + ":=" + ctx.GoFileContext.MustWriteGoTypes(t) + "{}\n")
			_buf.WriteString(marshalFnName + `(_buf,&` + varName + `)` + "\n")
			ts := t.String()
			if ctx.seenUnmarshalNamedStructMap[ts] {
				return
			}
			ctx.seenUnmarshalNamedStructMap[ts] = true
			thisFuncBuffer := &bytes.Buffer{}
			thisFuncBuffer.WriteString("func " + marshalFnName + "(_buf *udwRspLib.GoBuffer,_var *" + ctx.GoFileContext.MustWriteGoTypes(typ) + "){\n")
			structType := t.GetUnderType().(*udwGoParser.StructType)
			for _, f := range structType.Field {
				var1Name := ctx.GetNextVarString()
				ctx.GenGoUnmarshal(thisFuncBuffer, f.Elem, var1Name)
				thisFuncBuffer.WriteString("	_var." + f.Name + "=" + var1Name + "\n")
			}
			thisFuncBuffer.WriteString("}\n")
			ctx.GoFileFuncBuffer.Write(thisFuncBuffer.Bytes())
			return
		}
		var1Name := ctx.GetNextVarString()
		ctx.GenGoUnmarshal(_buf, t.GetUnderType(), var1Name)
		_buf.WriteString("	" + varName + ":=" + ctx.GoFileContext.MustWriteGoTypes(typ) + "(" + var1Name + ")\n")
	case udwGoParser.Ptr:
		t := typ.(*udwGoParser.PointerType)
		var1Name := ctx.GetNextVarString()
		_buf.WriteString(var1Name + `:=_buf.ReadBool()
	var ` + varName + ` ` + ctx.GoFileContext.MustWriteGoTypes(t) + `
	if ` + var1Name + `{
`)
		var2Name := ctx.GetNextVarString()
		ctx.GenGoUnmarshal(_buf, t.Elem, var2Name)
		_buf.WriteString("	" + varName + "=&" + var2Name + "\n")
		_buf.WriteString(`}
`)

	default:
		panic(fmt.Errorf("type[%s] fn.Name[%s] not support", typ.String(), ctx.CurrentProcessFnName))
	}
}
