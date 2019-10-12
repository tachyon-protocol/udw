package udwGoParser

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwStrings"
	"path"
	"reflect"
	"strings"
)

type Type interface {
	GetKind() Kind
	Equal(typ Type) bool

	String() string

	GetElem() Type
}

type FuncOrMethodDeclaration struct {
	Name            string
	ReceiverVarName string
	ReceiverType    Type
	InParameter     []FuncParameter
	OutParameter    []FuncParameter
	IsVariadic      bool
}

func (t *FuncOrMethodDeclaration) GetName() string {
	return t.Name
}
func (t *FuncOrMethodDeclaration) SetName(name string) {
	t.Name = name
}
func (t *FuncOrMethodDeclaration) GetInParameter() []FuncParameter {
	return t.InParameter
}
func (t *FuncOrMethodDeclaration) GetOutParameter() []FuncParameter {
	return t.OutParameter
}

func (t *FuncOrMethodDeclaration) GetKind() Kind {
	if t.ReceiverType == nil {
		return DefinedFunc
	} else {
		return Method
	}
}

func (t *FuncOrMethodDeclaration) Clone() *FuncOrMethodDeclaration {
	newT := *t
	return &newT
}

func (t *FuncOrMethodDeclaration) IsExport() bool {
	return IsNameGoExport(t.Name)
}

func (t *FuncOrMethodDeclaration) HasInOrOutParameter() bool {
	return len(t.InParameter) > 0 || len(t.OutParameter) > 0
}
func (t *FuncOrMethodDeclaration) Equal(typ Type) bool {
	typS, ok := typ.(*FuncOrMethodDeclaration)
	if !ok {
		return false
	}
	if len(t.InParameter) != len(typS.InParameter) {
		return false
	}
	if len(t.OutParameter) != len(typS.OutParameter) {
		return false
	}
	if t.Name != typS.Name {
		return false
	}
	if t.ReceiverType.Equal(typS.ReceiverType) == false {
		return false
	}
	for i := range t.InParameter {
		if t.InParameter[i].Type.Equal(typS.InParameter[i].Type) == false {
			return false
		}
	}
	for i := range t.OutParameter {
		if t.OutParameter[i].Type.Equal(typS.OutParameter[i].Type) == false {
			return false
		}
	}
	if t.IsVariadic != typS.IsVariadic {
		return false
	}
	return true
}

func (t *FuncOrMethodDeclaration) EqualMethod(typ Type) bool {
	typS, ok := typ.(*FuncOrMethodDeclaration)
	if !ok {
		return false
	}
	if len(t.InParameter) != len(typS.InParameter) {
		return false
	}
	if len(t.OutParameter) != len(typS.OutParameter) {
		return false
	}
	if t.Name != typS.Name {
		return false
	}

	for i := range t.InParameter {
		if t.InParameter[i].Type.Equal(typS.InParameter[i].Type) == false {
			return false
		}
	}
	for i := range t.OutParameter {
		if t.OutParameter[i].Type.Equal(typS.OutParameter[i].Type) == false {
			return false
		}
	}
	if t.IsVariadic != typS.IsVariadic {
		return false
	}
	return true
}
func (t *FuncOrMethodDeclaration) String() string {
	buf := udwBytes.BufWriter{}
	buf.WriteString("func ")
	if t.ReceiverType != nil {
		buf.WriteString("(")
		buf.WriteString(t.ReceiverType.String())
		buf.WriteString(")")
	}
	if t.Name != "" {
		buf.WriteString(t.Name)
	}
	buf.WriteString("(")
	for i, para := range t.InParameter {
		if para.Name != "" {
			buf.WriteString(para.Name)
			buf.WriteString(" ")
		}
		buf.WriteString(para.Type.String())
		if i != len(t.InParameter)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	if len(t.OutParameter) > 0 {
		buf.WriteString("(")
		for i, para := range t.OutParameter {
			if para.Name != "" {
				buf.WriteString(para.Name)
				buf.WriteString(" ")
			}
			buf.WriteString(para.Type.String())
			if i != len(t.OutParameter)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")
	}
	return buf.GetString()
}
func (t *FuncOrMethodDeclaration) GetElem() Type {
	panic("Not Allow Call GetElem on FuncOrMethodDeclaration")
}

type FuncType struct {
	InParameter  []FuncParameter
	OutParameter []FuncParameter
}

func (t *FuncType) GetKind() Kind {
	return Func
}
func (t *FuncType) Equal(typ Type) bool {
	return reflect.DeepEqual(t, typ)
}
func (t *FuncType) String() string {
	return "[FuncType.String] TODO"

}
func (t *FuncType) GetElem() Type {
	panic("Not Allow Call GetElem on *FuncType")
}

type NamedType struct {
	PkgImportPath string
	Name          string
	UnderType     Type
	program       *Program
}

func (t *NamedType) GetKind() Kind {
	return Named
}
func (t *NamedType) Equal(typ Type) bool {
	typN, ok := typ.(*NamedType)
	if !ok {
		return false
	}
	return typN.PkgImportPath == t.PkgImportPath && typN.Name == t.Name
}

func (t *NamedType) GetUnderType() Type {
	if t.UnderType == nil {
		definer := t.program.GetNamedType(t.PkgImportPath, t.Name)
		if definer == nil {
			panic(fmt.Errorf("[NamedType.GetUnderType] [%s]", t.String()))
		}
		t.UnderType = definer.UnderType
	}
	return t.UnderType
}
func (t *NamedType) String() string {
	return t.PkgImportPath + "." + t.Name
}

func (t *NamedType) GetPackageName() string {
	return path.Base(t.PkgImportPath)
}

func (t *NamedType) GetElem() Type {
	return t.GetUnderType()
}

func IsTypeGoTimeObj(typ Type) bool {
	return typ.Equal(GetGoTimeType())
}

var gGoTimeType = NamedType{
	PkgImportPath: "time",
	Name:          "Time",
}

func GetGoTimeType() *NamedType {
	return &gGoTimeType
}

func NewNamedType(PackagePath string, Name string) *NamedType {
	return &NamedType{
		PkgImportPath: PackagePath,
		Name:          Name,
	}
}

type StructType struct {
	Field []StructField
}

func (t *StructType) GetKind() Kind {
	return Struct
}

func (t *StructType) Equal(typ Type) bool {
	typS, ok := typ.(*StructType)
	if !ok {
		return false
	}
	if len(t.Field) != len(typS.Field) {
		return false
	}
	for i := range t.Field {
		if t.Field[i].Equal(typS.Field[i]) == false {
			return false
		}
	}
	return true
}
func (t *StructType) String() string {
	s, _ := MustWriteGoTypes("///", t)
	return s
}

func (t *StructType) IsEmpty() bool {
	return len(t.Field) == 0
}
func (t *StructType) GetElem() Type {
	panic("Not Allow Call GetElem On StructType")
}

func GetStructOrNamedStructFromType(elemT Type) *StructType {
	switch elemT.GetKind() {
	case Named:
		t := elemT.(*NamedType).GetUnderType()
		if t.GetKind() == Struct {
			t2 := t.(*StructType)
			return t2
		}
	case Struct:
		t2 := elemT.(*StructType)
		return t2
	}
	return nil
}

type StructField struct {
	Name             string
	Elem             Type
	IsAnonymousField bool
	Tag              string
}

func (sf StructField) Equal(sf2 StructField) bool {
	return sf.Elem.Equal(sf2.Elem) && sf.Name == sf2.Name && sf.IsAnonymousField == sf.IsAnonymousField
}
func (sf StructField) IsExport() bool {
	return IsNameGoExport(sf.Name)
}

func (sf StructField) IsInTagList(key string, search string) bool {
	value := sf.GetTagL1ValueByKeyIgnoreError(key)
	list := strings.Split(value, `,`)
	return udwStrings.IsInSlice(list, search)
}

type MapType struct {
	Key   Type
	Value Type
}

func (t *MapType) GetKind() Kind {
	return Map
}

func (t *MapType) Equal(typ Type) bool {
	typS, ok := typ.(*MapType)
	if !ok {
		return false
	}
	return t.Key.Equal(typS.Key) && t.Value.Equal(typS.Value)
}
func (t *MapType) String() string {
	return "map[" + t.Key.String() + "]" + t.Value.String()
}
func (t *MapType) GetElem() Type {
	return t.Value
}

type InterfaceType struct {
}

func (t InterfaceType) GetKind() Kind {
	return Interface
}

func (t InterfaceType) Equal(typ Type) bool {
	return reflect.DeepEqual(t, typ)
}
func (t InterfaceType) String() string {
	s, _ := MustWriteGoTypes("///", t)
	return s
}
func (t InterfaceType) GetElem() Type {
	panic("Not Allow Call GetElem on InterfaceType")
}

type PointerType struct {
	Elem Type
}

func (t *PointerType) GetKind() Kind {
	return Ptr
}
func (t *PointerType) Equal(typ Type) bool {
	typN, ok := typ.(*PointerType)
	if !ok {
		return false
	}
	return t.Elem.Equal(typN.Elem)
}
func (t *PointerType) String() string {
	return "*" + t.Elem.String()
}
func (t *PointerType) GetElem() Type {
	return t.Elem
}

func NewPointer(elem Type) *PointerType {
	return &PointerType{Elem: elem}
}

type FuncParameter struct {
	Name string
	Type Type
}

func (p FuncParameter) GetName() string {
	return p.Name
}
func (p FuncParameter) GetType() Type {
	return p.Type
}

type SliceType struct {
	Elem Type
}

func (t *SliceType) GetKind() Kind {
	return Slice
}
func (t *SliceType) Equal(typ Type) bool {
	typS, ok := typ.(*SliceType)
	if !ok {
		return false
	}
	return t.Elem.Equal(typS.Elem)
}
func (t *SliceType) String() string {
	return "[]" + t.Elem.String()
}
func (t *SliceType) GetElem() Type {
	return t.Elem
}

type ArrayType struct {
	Size int
	Elem Type
}

func (t *ArrayType) GetKind() Kind {
	return Array
}
func (t *ArrayType) Equal(typ Type) bool {
	typS, ok := typ.(*ArrayType)
	if !ok {
		return false
	}
	return t.Elem.Equal(typS.Elem) && t.Size == typS.Size
}
func (t *ArrayType) String() string {
	s, _ := MustWriteGoTypes("///", t)
	return s
}
func (t *ArrayType) GetElem() Type {
	return t.Elem
}

type ChanType struct {
	Dir  ChanDir
	Elem Type
}

func (t *ChanType) GetKind() Kind {
	return Chan
}
func (t *ChanType) Equal(typ Type) bool {
	typS, ok := typ.(*ChanType)
	if !ok {
		return false
	}
	return t.Elem.Equal(typS.Elem) && t.Dir == typS.Dir
}
func (t *ChanType) String() string {
	s, _ := MustWriteGoTypes("///", t)
	return s
}
func (t *ChanType) GetElem() Type {
	return t.Elem
}

type BuiltinType string

func (t BuiltinType) GetKind() Kind {
	return getKindFromBuiltinType(string(t))
}

func (t BuiltinType) Equal(typ Type) bool {
	typS, ok := typ.(BuiltinType)
	if !ok {
		return false
	}
	return t == typS
}

func (t BuiltinType) String() string {
	return string(t)
}
func (t BuiltinType) GetElem() Type {
	panic("Not Allow Call GetElem On BuiltinType")
}
func GetFloat64Type() Type {
	return BuiltinType("float64")
}

func GetErrorType() BuiltinType {
	return BuiltinType("error")
}
func GetStringType() BuiltinType {
	return BuiltinType("string")
}
func GetEmptyStructType() *StructType {
	return &StructType{}
}

type ChanDir int

const (
	RecvDir ChanDir = 1 << iota
	SendDir
	BothDir = RecvDir | SendDir
)
