package udwGoTypeMarshal

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwMap"
	"github.com/tachyon-protocol/udw/udwSort"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"path"
	"reflect"
	"strconv"
	"time"
)

func MustWriteObjectToMainPackage(obj interface{}) string {
	_buf := &bytes.Buffer{}
	MustWriteObjectToMainPackageWithBuf(obj, _buf)
	return _buf.String()
}

func MustWriteObjectToMainPackageASCIISafe(obj interface{}) string {
	_buf := &bytes.Buffer{}
	MustGoTypeMarshalWithBuf(MustGoTypeMarshalContext{
		AddTypeNameWithImportPackage: addTypeNameWithImportPackageToMainPackage,
		IsStringASCIISafe:            true,
	}, obj, _buf)
	return _buf.String()
}

func MustWriteObjectToMainPackageWithUnexport(obj interface{}) string {
	_buf := &bytes.Buffer{}
	MustGoTypeMarshalWithBuf(MustGoTypeMarshalContext{
		AddTypeNameWithImportPackage: addTypeNameWithImportPackageToMainPackage,
		IsMarshalUnexport:            true,
	}, obj, _buf)
	return _buf.String()
}

func MustWriteObjectToMainPackageWithBuf(obj interface{}, _buf *bytes.Buffer) {
	MustGoTypeMarshalWithBuf(MustGoTypeMarshalContext{
		AddTypeNameWithImportPackage: addTypeNameWithImportPackageToMainPackage,
	}, obj, _buf)
}

func addTypeNameWithImportPackageToMainPackage(pkgPath string, name string, buf *bytes.Buffer) {
	if "main" == pkgPath {
		buf.WriteString(name)
		return
	} else {
		buf.WriteString(path.Base(pkgPath) + "." + name)
		return
	}
}

type MustGoTypeMarshalContext struct {
	AddTypeNameWithImportPackage func(pkgPath string, name string, buf *bytes.Buffer)
	IsStringASCIISafe            bool
	IsMarshalUnexport            bool
}

func MustGoTypeMarshalWithBuf(ctx MustGoTypeMarshalContext, obj interface{}, _buf *bytes.Buffer) {
	switch objI := obj.(type) {

	case map[string]string:
		if len(objI) == 0 {
			_buf.WriteString(`map[string]string(nil)`)
		} else {
			_buf.WriteString("map[string]string{\n")
			keyValueList := udwMap.MapStringStringToKeyValuePairListAes(objI)
			for _, pair := range keyValueList {
				writeStringWithCtx(ctx, pair.Key, _buf)
				_buf.WriteString(":")
				writeStringWithCtx(ctx, pair.Value, _buf)
				_buf.WriteString(",\n")
			}
			_buf.WriteString("}")
			return
		}
	case []byte:
		_buf.WriteString("[]byte(")
		s := string(objI)
		writeStringWithCtx(ctx, s, _buf)
		_buf.WriteString(")")
		return
	case string:
		writeStringWithCtx(ctx, objI, _buf)
		return
	case int:
		_buf.WriteString(strconv.Itoa(objI))
		return
	case int64:
		_buf.WriteString(udwStrconv.FormatInt64(objI))
		return
	case uint8:
		_buf.WriteString(udwStrconv.FormatUint64(uint64(objI)))
		return
	case uint16:
		_buf.WriteString(udwStrconv.FormatUint64(uint64(objI)))
		return
	case uint32:
		_buf.WriteString(udwStrconv.FormatUint64(uint64(objI)))
		return
	case uint64:
		_buf.WriteString(udwStrconv.FormatUint64(objI))
		return
	case uintptr:
		_buf.WriteString(udwStrconv.FormatUint64(uint64(objI)))
		return
	case float64:
		_buf.WriteString(udwStrconv.FormatFloat(objI))
		return
	case bool:
		_buf.WriteString(strconv.FormatBool(objI))
		return
	case nil:
		_buf.WriteString("nil")
		return
	case time.Time:
		ctx.AddTypeNameWithImportPackage("time", "Unix", _buf)
		_buf.WriteByte('(')
		nano := objI.UnixNano()
		_buf.WriteString(strconv.FormatInt(nano/1e9, 10))
		_buf.WriteByte(',')
		_buf.WriteString(strconv.FormatInt(nano%1e9, 10))
		_buf.WriteString(").UTC()")

	default:
		mustMarshalFromReflectL2(ctx, reflect.ValueOf(obj), _buf)
		return
	}
}

func mustMarshalFromReflectL1(ctx MustGoTypeMarshalContext, reflectValue reflect.Value, _buf *bytes.Buffer) {
	if reflectValue.CanInterface() {
		MustGoTypeMarshalWithBuf(ctx, reflectValue.Interface(), _buf)
	} else {
		mustMarshalFromReflectL2(ctx, reflectValue, _buf)
	}
}

func mustMarshalFromReflectL2(ctx MustGoTypeMarshalContext, reflectValue reflect.Value, _buf *bytes.Buffer) {
	switch reflectValue.Kind() {
	case reflect.Ptr:
		if reflectValue.IsNil() {
			_buf.WriteString("nil")
			return
		}
		elemValue := reflectValue.Elem()
		switch elemValue.Kind() {
		case reflect.Bool, reflect.String, reflect.Float64, reflect.Int:
			canNotAddressableFnName := ""
			switch elemValue.Type().Name() {
			case "bool":
				canNotAddressableFnName = "PtrBool"
			case "string":
				canNotAddressableFnName = "PtrString"
			case "float64":
				canNotAddressableFnName = "PtrFloat64"
			case "int":
				canNotAddressableFnName = "PtrInt"
			}
			if canNotAddressableFnName != "" {
				ctx.AddTypeNameWithImportPackage("github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal/udwGoTypeMarshalLib", canNotAddressableFnName, _buf)
				_buf.WriteByte('(')
				mustMarshalFromReflectL1(ctx, elemValue, _buf)
				_buf.WriteByte(')')
				return
			}

		case reflect.Struct, reflect.Slice:
			_buf.WriteByte('&')
			mustMarshalFromReflectL1(ctx, elemValue, _buf)
			return
		}

		_buf.WriteString("func()")
		WriteReflectTypeNameToGoFile(ctx, reflectValue.Type(), _buf)
		_buf.WriteString("{_a:=")
		mustMarshalFromReflectL1(ctx, elemValue, _buf)
		_buf.WriteString(";return &_a}()")
		return
	case reflect.Struct:
		t := reflectValue.Type()
		WriteReflectTypeNameToGoFile(ctx, t, _buf)
		_buf.WriteString("{\n")
		nf := t.NumField()
		for i := 0; i < nf; i++ {
			field := t.Field(i)
			if ctx.IsMarshalUnexport == false && field.PkgPath != "" {
				continue
			}
			thisV := reflectValue.Field(i)
			if isEmptyValue(thisV) {
				continue
			}
			_buf.WriteString(field.Name)
			_buf.WriteString(":")
			mustMarshalFromReflectL1(ctx, thisV, _buf)
			_buf.WriteString(",\n")
		}
		_buf.WriteString(`}`)
		return

	case reflect.Slice, reflect.Array:
		t := reflectValue.Type()
		WriteReflectTypeNameToGoFile(ctx, t, _buf)
		l := reflectValue.Len()
		if l == 0 {
			_buf.WriteString("(nil)")
			return
		}
		_buf.WriteString("{\n")
		length := reflectValue.Len()
		for i := 0; i < length; i++ {
			mustMarshalFromReflectL1(ctx, reflectValue.Index(i), _buf)
			_buf.WriteString(",\n")
		}
		_buf.WriteString("}")
		return
	case reflect.Map:
		t := reflectValue.Type()
		_buf.WriteString("map[")
		WriteReflectTypeNameToGoFile(ctx, t.Key(), _buf)
		_buf.WriteString("]")
		WriteReflectTypeNameToGoFile(ctx, t.Elem(), _buf)
		l := reflectValue.Len()
		if l == 0 {
			_buf.WriteString("(nil)")
			return
		}
		_buf.WriteString("{\n")
		keyList := reflectValue.MapKeys()
		keyType := t.Key()
		switch keyType.Kind() {
		case reflect.String:
			udwSort.InterfaceCallbackSortWithIndexLess(keyList, func(a int, b int) bool {
				return keyList[a].String() < keyList[b].String()
			})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			udwSort.InterfaceCallbackSortWithIndexLess(keyList, func(a int, b int) bool {
				return keyList[a].Uint() < keyList[b].Uint()
			})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			udwSort.InterfaceCallbackSortWithIndexLess(keyList, func(a int, b int) bool {
				return keyList[a].Int() < keyList[b].Int()
			})
		default:

			panic(fmt.Errorf("[mustMarshalFromReflect] map sort TODO %s", reflectValue.Kind().String()))
		}
		for i := range keyList {
			mustMarshalFromReflectL1(ctx, keyList[i], _buf)
			_buf.WriteString(":")
			mustMarshalFromReflectL1(ctx, reflectValue.MapIndex(keyList[i]), _buf)
			_buf.WriteString(",\n")
		}
		_buf.WriteString("}")
		return
	case reflect.String,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float64, reflect.Float32:
		isBuiltInConstant := false
		switch reflectValue.Type().Name() {
		case "string", "bool", "int", "uint":
			isBuiltInConstant = true
		}
		if isBuiltInConstant == false {
			t := reflectValue.Type()
			WriteReflectTypeNameToGoFile(ctx, t, _buf)
			_buf.WriteString("(")
		}
		switch reflectValue.Kind() {
		default:
			panic("fdhp4p538j")
		case reflect.String:
			writeStringWithCtx(ctx, reflectValue.String(), _buf)
		case reflect.Bool:
			_buf.WriteString(strconv.FormatBool(reflectValue.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_buf.WriteString(strconv.FormatInt(reflectValue.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			_buf.WriteString(strconv.FormatUint(reflectValue.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			_buf.WriteString(strconv.FormatFloat(reflectValue.Float(), 'f', -1, 64))
		}
		if isBuiltInConstant == false {
			_buf.WriteString(")")
		}
	case reflect.Func:
		t := reflectValue.Type()
		WriteReflectTypeNameToGoFile(ctx, t, _buf)
		_buf.WriteString(`{
	panic("not implement in udwGoTypeMarshal ` + udwStrconv.FormatUint64Hex(uint64(reflectValue.Pointer())) + `")
}`)
	default:
		panic(fmt.Errorf("[mustMarshalFromReflect] TODO %s", reflectValue.Kind().String()))
	}
}

func WriteReflectTypeNameToGoFile(Ctx MustGoTypeMarshalContext, Typ reflect.Type, Buf *bytes.Buffer) {
	writeReflectTypeNameToGoFileL1(WriteReflectTypeNameToGoFileRequest{
		Ctx: Ctx,
		Typ: Typ,
		Buf: Buf,
	})
}

type WriteReflectTypeNameToGoFileRequest struct {
	Ctx      MustGoTypeMarshalContext
	Typ      reflect.Type
	Buf      *bytes.Buffer
	IsPtrAnd bool
}

func writeReflectTypeNameToGoFileL1(req WriteReflectTypeNameToGoFileRequest) {
	switch req.Typ.Kind() {
	case reflect.Ptr:
		if req.IsPtrAnd {
			req.Buf.WriteByte('&')
		} else {
			req.Buf.WriteByte('*')
		}
		typ := req.Typ
		req.Typ = typ.Elem()
		writeReflectTypeNameToGoFileL1(req)
	case reflect.String,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float64, reflect.Float32:

		name := req.Typ.Name()
		switch name {
		case "string", "uint8", "uint16", "uint32", "uint64", "uintptr", "bool", "int", "int8", "int16", "int32", "int64", "float32", "float64":
			req.Buf.WriteString(name)
			return
		}
		req.Ctx.AddTypeNameWithImportPackage(req.Typ.PkgPath(), req.Typ.Name(), req.Buf)
		return
	case reflect.Struct:
		if req.Typ.Name() != "" {
			req.Ctx.AddTypeNameWithImportPackage(req.Typ.PkgPath(), req.Typ.Name(), req.Buf)
		} else {

			req.Buf.WriteString(fmt.Sprintf("%T", reflect.New(req.Typ).Elem().Interface()))
		}
		return
	case reflect.Slice:
		req.Buf.WriteString("[]")
		req.Typ = req.Typ.Elem()
		writeReflectTypeNameToGoFileL1(req)
	case reflect.Array:
		req.Buf.WriteString("[")
		req.Buf.WriteString(strconv.Itoa(req.Typ.Len()))
		req.Buf.WriteString("]")
		req.Typ = req.Typ.Elem()
		writeReflectTypeNameToGoFileL1(req)
	case reflect.Map:
		req.Buf.WriteString("map[")
		typ := req.Typ
		req.Typ = typ.Key()
		writeReflectTypeNameToGoFileL1(req)
		req.Buf.WriteString("]")
		req.Typ = typ.Elem()
		writeReflectTypeNameToGoFileL1(req)
	case reflect.Func:
		req.Buf.WriteString("func (")
		typ := req.Typ
		inNum := typ.NumIn()
		outNum := typ.NumOut()
		for i := 0; i < inNum; i++ {
			req.Typ = typ.In(i)
			writeReflectTypeNameToGoFileL1(req)
			if i != inNum-1 {
				req.Buf.WriteString(",")
			}
		}
		req.Buf.WriteString(")(")
		for i := 0; i < outNum; i++ {
			req.Typ = typ.Out(i)
			writeReflectTypeNameToGoFileL1(req)
			if i != outNum-1 {
				req.Buf.WriteString(",")
			}
		}
		req.Buf.WriteString(")")
	case reflect.Interface:
		name := req.Typ.Name()
		switch name {
		case "error":
			req.Buf.WriteString(name)
			return
		}
		req.Ctx.AddTypeNameWithImportPackage(req.Typ.PkgPath(), req.Typ.Name(), req.Buf)

	default:
		panic(fmt.Errorf("[writeReflectTypeNameToGoFile] TODO %s", req.Typ.Kind().String()))
	}
}

func isBuiltInSingleType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float64, reflect.Float32:
		name := typ.Name()
		switch name {
		case "string", "uint8", "uint16", "uint32", "uint64", "uintptr", "bool", "int", "int8", "int16", "int32", "int64", "float32", "float64":
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		return v.Bool() == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr:
		return v.IsNil() || isEmptyValue(v.Elem())
	case reflect.Float64, reflect.Float32:
		return v.Float() == 0
	case reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		l := v.Len()
		if l == 0 {
			return true
		}
		for i := 0; i < l; i++ {
			if isEmptyValue(v.Index(i)) == false {
				return false
			}
		}
		return true
	case reflect.Struct:
		nf := v.NumField()
		for i := 0; i < nf; i++ {
			thisV := v.Field(i)
			if isEmptyValue(thisV) == false {
				return false
			}
		}
		return true
	case reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

func isEmptyStructField(v reflect.Value, field *reflect.StructField) bool {
	if len(field.Index) == 1 {
		return isEmptyValue(v.Field(field.Index[0]))
	}
	for i, x := range field.Index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					return true
				}
				if v.Elem().Kind() == reflect.Struct {
					v = v.Elem()
				}
			}
		}
		v = v.Field(x)
	}
	if isEmptyValue(v) {
		return true
	}
	return false
}

func writeStringWithCtx(ctx MustGoTypeMarshalContext, s string, _buf *bytes.Buffer) {
	if ctx.IsStringASCIISafe {
		out := WriteStringToGolangASCII(s)
		_buf.WriteString(out)
	} else {
		WriteStringToGolangToBuf(s, _buf)
	}
}

func IsObjHasInitAlloc(val reflect.Value, typ reflect.Type) bool {
	isAlloc := IsReflectTypeHasInitAlloc(typ)
	if isAlloc == false {
		return false
	}
	switch typ.Kind() {
	case reflect.Ptr:
		return IsObjHasInitAlloc(val.Elem(), typ.Elem())
	case reflect.String,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float64, reflect.Float32:
		return false
	case reflect.Struct:
		if typ.PkgPath() == "time" && typ.Name() == "Time" {
			return true
		}
		fieldNum := typ.NumField()
		for i := 0; i < fieldNum; i++ {
			thisField := typ.Field(i)
			thisVal := val.Field(i)
			isAlloc := IsObjHasInitAlloc(thisVal, thisField.Type)
			if isAlloc == true {
				return true
			}
		}
		return false
	case reflect.Slice:
		l := val.Len()
		for i := 0; i < l; i++ {
			isAlloc := IsObjHasInitAlloc(val.Index(i), typ.Elem())
			if isAlloc == true {
				return true
			}
		}
		return false
	case reflect.Array:
		l := val.Len()
		for i := 0; i < l; i++ {
			isAlloc := IsObjHasInitAlloc(val.Index(i), typ.Elem())
			if isAlloc == true {
				return true
			}
		}
		return false

	case reflect.Map:
		if val.Len() == 0 {
			return false
		}
		return true
	default:
		panic("[IsReflectTypeNameInitAlloc] TODO " + typ.Kind().String())
	}
}

func IsReflectTypeHasInitAlloc(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Ptr:
		return IsReflectTypeHasInitAlloc(typ.Elem())
	case reflect.String,
		reflect.Uint8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool, reflect.Float64, reflect.Float32:
		return false
	case reflect.Struct:
		if typ.PkgPath() == "time" && typ.Name() == "Time" {
			return true
		}
		fieldNum := typ.NumField()
		for i := 0; i < fieldNum; i++ {
			thisField := typ.Field(i)
			isAlloc := IsReflectTypeHasInitAlloc(thisField.Type)
			if isAlloc == true {
				return true
			}
		}
		return false
	case reflect.Slice:
		return IsReflectTypeHasInitAlloc(typ.Elem())
	case reflect.Array:
		return IsReflectTypeHasInitAlloc(typ.Elem())
	case reflect.Map:
		return true
	default:
		panic("[IsReflectTypeNameInitAlloc] TODO " + typ.Kind().String())
	}
}
