package udwGoParser

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"path"
	"strconv"
)

func MustWriteGoTypes(thisPackagePath string, typi Type) (s string, addPkgPathList []string) {
	switch typ := typi.(type) {
	case *FuncType:
		panic("TODO MustWriteGoTypes *FuncType")
	case *NamedType:
		if thisPackagePath == typ.PkgImportPath {
			return typ.Name, nil
		}
		return path.Base(typ.PkgImportPath) + "." + typ.Name, []string{typ.PkgImportPath}
	case *PointerType:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem)
		return "*" + s, addPkgPathList
	case *SliceType:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem)
		return "[]" + s, addPkgPathList
	case *MapType:
		ks, kaddPkgPathList := MustWriteGoTypes(thisPackagePath, typ.Key)
		vs, vaddPkgPathList := MustWriteGoTypes(thisPackagePath, typ.Value)
		return "map[" + ks + "]" + vs, append(kaddPkgPathList, vaddPkgPathList...)
	case BuiltinType:
		return string(typ), nil
	case *ArrayType:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem)
		return "[" + strconv.Itoa(typ.Size) + "]" + s, addPkgPathList
	case InterfaceType:

		return "interface{}", nil
	case *StructType:
		if len(typ.Field) == 0 {
			return "struct{}", nil
		}
		_buf := bytes.Buffer{}
		_buf.WriteString(` struct{
`)
		for _, field := range typ.Field {
			ks, kaddPkgPathList := MustWriteGoTypes(thisPackagePath, field.Elem)
			addPkgPathList = append(addPkgPathList, kaddPkgPathList...)
			if field.IsAnonymousField {
				_buf.WriteString(ks)
			} else {
				_buf.WriteString(field.Name + ` ` + ks)
			}
			if field.Tag != "" {
				_buf.WriteString(` ` + udwGoTypeMarshal.WriteStringToGolang(field.Tag))
			}
			_buf.WriteByte('\n')
		}
		_buf.WriteString(`}`)
		return _buf.String(), addPkgPathList
	default:
		panic(fmt.Errorf("[MustWriteGoTypes] Not implement go/types [%T]",
			typi))
	}

}
