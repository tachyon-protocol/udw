package udwGoParser

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwStrings"
)

func QybStructMethodInputAndOutputToMap(pkgPath string, objectName string) map[string]string {
	out := map[string]string{}
	pkg := MustParsePackage(udwProjectPath.MustGetProjectPath(), pkgPath)
	for _, oneMethod := range pkg.GetAllMethodOnNamedType(objectName) {
		if !oneMethod.IsExport() {
			continue
		}
		for _, input := range oneMethod.InParameter {
			dealTypeL1(oneMethod.Name+".input."+udwStrings.FirstLetterToUpper(input.GetName()), input.Type, out)
		}
		for i, output := range oneMethod.GetOutParameter() {
			name := udwStrings.FirstLetterToUpper(output.GetName())
			if name == "" {
				if output.GetType().Equal(GetErrorType()) {
					name = "Err"
				} else {
					name = fmt.Sprintf("Out_%d", i)
				}
			}
			dealTypeL1(oneMethod.Name+".output."+udwStrings.FirstLetterToUpper(name), output.Type, out)
		}
	}
	return out
}

func QybPkgAllStructToMap(objectPkgPath string) map[string]string {
	out := map[string]string{}
	pkg := MustParsePackage(udwProjectPath.MustGetProjectPath(), objectPkgPath)
	for _, one := range pkg.namedTypeList {
		t, ok := one.GetUnderType().(*StructType)
		if ok {
			dealTypeL1(one.String(), t, out)
		}
	}
	return out
}

func dealTypeL1(name string, ot Type, mark map[string]string) {
	if mark[name] != "" {
		return
	}
	if name != "" {
		val := dealTypeL2(ot)
		if val != "" && val != name {
			mark[name] = val
		}
	}
	getString := func(in Type) string {
		switch in.GetKind() {
		case Array, Slice, Map, Ptr:
			return ""
		default:
			return in.String()
		}
	}
	switch ot.GetKind() {
	case Struct:
		st := ot.(*StructType)
		for _, one := range st.Field {
			dealTypeL1(name+"."+one.Name, one.Elem, mark)
		}
	case Named:
		mt := ot.(*NamedType)
		if IsTypeGoTimeObj(mt) {
			return
		}
		if mt.PkgImportPath == "sync" && (mt.Name == "Mutex" || mt.Name == "RWMutex") {
			return
		}
		if mt.PkgImportPath == "net/http" {
			return
		}
		dealTypeL1(getString(mt), mt.GetUnderType(), mark)
	case Map:
		mt := ot.(*MapType)
		dealTypeL1(getString(mt.Value), mt.Value, mark)
		dealTypeL1(getString(mt.Key), mt.Key, mark)
	case Slice:
		mt := ot.(*SliceType)
		dealTypeL1(getString(mt), mt.Elem, mark)
	case Array:
		mt := ot.(*ArrayType)
		dealTypeL1(getString(mt), mt.Elem, mark)
	case Ptr:
		mt := ot.(*PointerType)
		dealTypeL1(getString(mt), mt.Elem, mark)
	case Func:
		return
	}
}

func dealTypeL2(ot Type) string {
	switch ot.GetKind() {
	case Map:
		mt := ot.(*MapType)
		return mt.GetKind().String() + "[" + dealTypeL2(mt.Key) + "]" + dealTypeL2(mt.Value)
	case Slice:
		mt := ot.(*SliceType)
		return "[]" + dealTypeL2(mt.Elem)
	case Array:
		mt := ot.(*ArrayType)
		return "[]" + dealTypeL2(mt.Elem)
	case Struct:
		return "struct"
	case Ptr:
		mt := ot.(*PointerType)
		return "*" + dealTypeL2(mt.Elem)
	case Named:
		mt := ot.(*NamedType)
		return mt.String()
	case Func, Chan, Interface:
		return ""
	default:
		return ot.GetKind().String()
	}
}
