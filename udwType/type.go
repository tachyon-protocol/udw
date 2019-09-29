package udwType

import (
	"fmt"
	"reflect"
)

type FromStringInterface interface {
	FromString(s string) (reflect.Value, error)
}
type ToStringInterface interface {
	ToString(v reflect.Value) string
}
type StringConverterInterface interface {
	FromStringInterface
	ToStringInterface
}

type GetElemByStringInterface interface {
	GetElemByString(v reflect.Value, k string) (reflect.Value, UdwType, error)
}

type SaveScaleInterface interface {
	SaveScale(v reflect.Value, value string) error
}

type EditableByPathInterface interface {
	SaveByPath(v *reflect.Value, path Path, value string) (err error)
	DeleteByPath(v *reflect.Value, path Path) (err error)
}

type UdwType interface {
	EditableByPathInterface
	ReflectTypeGetter
}
type UdwTypeAndToStringInterface interface {
	UdwType
	ToStringInterface
}

type ReflectTypeGetter interface {
	GetReflectType() reflect.Type
}

type GetElemByStringAndReflectTypeGetterInterface interface {
	GetElemByStringInterface
	ReflectTypeGetter
}

func TypeOf(rt reflect.Type) (UdwType, error) {
	switch rt {
	case GetDateTimeReflectType():
		t := &DateTimeType{reflectTypeGetterImp: reflectTypeGetterImp{rt}}
		t.saveScaleFromStringer = saveScaleFromStringer{t, t}
		t.saveScaleEditabler = saveScaleEditabler{t, t}
		return t, nil
	}

	switch rt.Kind() {
	case reflect.Ptr:
		t := &PtrType{reflectTypeGetterImp: reflectTypeGetterImp{rt}}
		return t, nil
	case reflect.Bool:
		t := &BoolType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.saveScaleFromStringer = saveScaleFromStringer{t, t}
		t.saveScaleEditabler = saveScaleEditabler{t, t}
		return t, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		t := &IntType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.saveScaleFromStringer = saveScaleFromStringer{t, t}
		t.saveScaleEditabler = saveScaleEditabler{t, t}
		return t, nil
	case reflect.Float32, reflect.Float64:
		t := &FloatType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.saveScaleFromStringer = saveScaleFromStringer{t, t}
		t.saveScaleEditabler = saveScaleEditabler{t, t}
		return t, nil
	case reflect.String:
		t := &StringType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.saveScaleFromStringer = saveScaleFromStringer{t, t}
		t.saveScaleEditabler = saveScaleEditabler{t, t}
		return t, nil
	case reflect.Array:
		t := &ArrayType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.getElemByStringEditorabler = getElemByStringEditorabler{t, t}
		return t, nil
	case reflect.Slice:
		t := &SliceType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.getElemByStringEditorabler = getElemByStringEditorabler{t, t}
		return t, nil
	case reflect.Map:
		t := &MapType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		return t, nil
	case reflect.Struct:
		t := &StructType{
			reflectTypeGetterImp: reflectTypeGetterImp{rt},
		}
		t.getElemByStringEditorabler = getElemByStringEditorabler{t, t}
		return t, nil
	default:
		return nil, fmt.Errorf("not support type kind: %s", rt.Kind().String())
	}
}
