package udwRpc2Builder

import (
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
	"path"
	"strconv"
)

type RpcService struct {
	List []RpcApi
}
type RpcApi struct {
	Pos                 int
	Name                string
	InputParameterList  []RpcParameter
	OutputParameterList []RpcParameter
}
type RpcParameter struct {
	Pos  int
	Name string
	Type RpcType

	IsClientIp  bool
	IsHijackCtx bool
}

type RpcType struct {
	Kind       string
	Elem       *RpcType
	FieldList  []RpcField
	StructName string
	GoPkg      string
}

type RpcField struct {
	Pos  int
	Name string
	Type RpcType
}

const (
	RpcTypeKindNamedStruct = "namedStruct"
	RpcTypeKindString      = "string"
	RpcTypeKindInt         = "int"
	RpcTypeKindBool        = "bool"
	RpcTypeKindSlice       = "slice"
)

func initRpcDefine(RpcDefine *RpcService) {
	for i := range RpcDefine.List {
		api := &RpcDefine.List[i]
		api.Pos = i + 1
		for i := range api.InputParameterList {
			parameter := &api.InputParameterList[i]
			parameter.Pos = i + 1
			if parameter.Name == "" {
				parameter.Name = "fi" + strconv.Itoa(i+1)
			}
			thisTyp := &parameter.Type
			initRpcType(thisTyp)
			if thisTyp.Kind == RpcTypeKindNamedStruct && thisTyp.StructName == "PeerIp" && thisTyp.GoPkg == "github.com/tachyon-protocol/udw/udwRpc2" {
				parameter.IsClientIp = true
			}
		}
		for i := range api.OutputParameterList {
			parameter := &api.OutputParameterList[i]
			parameter.Pos = i + 1
			if parameter.Name == "" {
				parameter.Name = "fo" + strconv.Itoa(i+1)
			}
			initRpcType(&parameter.Type)
		}
	}
}

func initRpcType(Type *RpcType) {
	switch Type.Kind {
	case RpcTypeKindSlice:
		initRpcType(Type.Elem)
	case RpcTypeKindNamedStruct:
		for i := range Type.FieldList {
			field := &Type.FieldList[i]
			field.Pos = i + 1
			initRpcType(&field.Type)
		}
	}
}

func mustRpcTypeToGoTypeName(writer *udwGoParser.GoFileWriter, Type *RpcType) string {
	switch Type.Kind {
	case RpcTypeKindString:
		return "string"
	case RpcTypeKindInt:
		return "int"
	case RpcTypeKindNamedStruct:
		writer.AddImportPath(Type.GoPkg)
		if Type.GoPkg == writer.GetPkgImportPath() {
			return Type.StructName
		}
		return path.Base(Type.GoPkg) + "." + Type.StructName
	default:
		panic("not support Type [" + Type.Kind + "]")
	}
}
