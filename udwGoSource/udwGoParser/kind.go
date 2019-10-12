package udwGoParser

import "strconv"

type Kind uint

const (
	Invalid       Kind = 0
	Bool          Kind = 1
	Int           Kind = 2
	Int8          Kind = 3
	Int16         Kind = 4
	Int32         Kind = 5
	Int64         Kind = 6
	Uint          Kind = 7
	Uint8         Kind = 8
	Uint16        Kind = 9
	Uint32        Kind = 10
	Uint64        Kind = 11
	Uintptr       Kind = 12
	Float32       Kind = 13
	Float64       Kind = 14
	Complex64     Kind = 15
	Complex128    Kind = 16
	Array         Kind = 17
	Chan          Kind = 18
	Func          Kind = 19
	Interface     Kind = 20
	Map           Kind = 21
	Ptr           Kind = 22
	Slice         Kind = 23
	String        Kind = 24
	Struct        Kind = 25
	UnsafePointer Kind = 26
	Method        Kind = 27
	Named         Kind = 28
	DefinedFunc   Kind = 29
)

func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "kind" + strconv.Itoa(int(k))
}

var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Ptr:           "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
	Method:        "method",
	Named:         "Named",
	DefinedFunc:   "DefinedFunc",
}

func getKindFromBuiltinType(typ string) Kind {
	switch typ {
	case "bool":
		return Bool
	case "byte":
		return Uint8
	case "complex128":
		return Complex128
	case "complex64":
		return Complex64
	case "error":
		return Interface
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "int":
		return Int
	case "int16":
		return Int16
	case "int32":
		return Int32
	case "int64":
		return Int64
	case "int8":
		return Int8
	case "rune":
		return Int32
	case "string":
		return String
	case "uint":
		return Uint
	case "uint16":
		return Uint16
	case "uint32":
		return Uint32
	case "uint64":
		return Uint64
	case "uint8":
		return Uint8
	case "uintptr":
		return Uintptr
	default:
		return Invalid
	}
}
