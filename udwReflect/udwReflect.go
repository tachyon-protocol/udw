package udwReflect

import (
	"reflect"
)

func GetTypeFullName(t reflect.Type) (name string) {
	if t.Kind() == reflect.Ptr {
		return GetTypeFullName(t.Elem())
	}
	if t.Name() == "" {
		return ""
	}

	if t.PkgPath() == "" {
		return t.Name()
	}
	return t.PkgPath() + "." + t.Name()
}

func IndirectType(v reflect.Type) reflect.Type {
	switch v.Kind() {
	case reflect.Ptr:
		return IndirectType(v.Elem())
	default:
		return v
	}

}
