package udwReflect

import (
	"reflect"
)

func DeepCopy(in interface{}) interface{} {
	inT := reflect.TypeOf(in)
	inV := reflect.ValueOf(in)
	outV := reflect.New(inT).Elem()
	deepCopyL1(inT, inV, outV)
	return outV.Interface()
}

func deepCopyL1(inT reflect.Type, inV reflect.Value, outV reflect.Value) {
	if !outV.CanSet() {
		return
	}
	switch inT.Kind() {
	case reflect.Ptr:
		if inV.IsNil() {
			return
		}
		outV.Set(reflect.New(inT.Elem()))
		deepCopyL1(inT.Elem(), inV.Elem(), outV.Elem())
	case reflect.Map:
		outV.Set(reflect.MakeMap(inT))
		for _, one := range inV.MapKeys() {
			tmpV := reflect.New(inT.Elem()).Elem()
			deepCopyL1(inT.Elem(), inV.MapIndex(one), tmpV)
			outV.SetMapIndex(one, tmpV)
		}
	case reflect.Slice, reflect.Array:
		len := inV.Len()
		outV.Set(reflect.MakeSlice(inT, len, len))
		for i := 0; i < len; i++ {
			tmpV := reflect.New(inT.Elem()).Elem()
			deepCopyL1(inT.Elem(), inV.Index(i), tmpV)
			outV.Index(i).Set(tmpV)
		}
	case reflect.Struct:
		if inT.String() == "time.Time" {
			outV.Set(inV)
		} else {
			for i := 0; i < inT.NumField(); i++ {
				deepCopyL1(inT.Field(i).Type, inV.Field(i), outV.Field(i))
			}
		}
	case reflect.Chan, reflect.Func, reflect.Interface:

	default:
		outV.Set(inV)
	}
}
