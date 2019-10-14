package udwGoParser

import (
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestTypeEqual(ot *testing.T) {
	var _ Type = &PointerType{}
	var _ Type = &NamedType{}
	var _ Type = BuiltinType("string")
	var _ Type = &StructType{}
	var _ Type = &MapType{}
	var _ Type = &SliceType{}
	var _ Type = &ArrayType{}
	var _ Type = &ChanType{}
	var _ Type = &FuncOrMethodDeclaration{}
	var _ Type = &InterfaceType{}
	var _ Type = &FuncType{}
	udwTest.Equal(GetErrorType().Equal(GetErrorType()), true)
	pkg := MustParsePackage(udwProjectPath.MustGetProjectPath(), "github.com/tachyon-protocol/udw/udwGoSource/udwGoParser")
	typ := pkg.LookupNamedType("NamedType")
	udwTest.Equal(typ.Equal(NewNamedType("github.com/tachyon-protocol/udw/udwGoSource/udwGoParser", "NamedType")), true)
	udwTest.Equal(typ.Equal(NewPointer(NewNamedType("github.com/tachyon-protocol/udw/udwGoSource/udwGoParser", "NamedType"))), false)
	pkg.GetNamedTypeMethodSet(typ)
	typp := NewPointer(typ)
	udwTest.Equal(typp.Equal(NewPointer(NewNamedType("github.com/tachyon-protocol/udw/udwGoSource/udwGoParser", "NamedType"))), true)
}
