package udwTypeTransform

import (
	"fmt"
	"reflect"
)

func SliceStructToMapStruct(in interface{}, out interface{}, idFieldName string) (err error) {
	inV := reflect.ValueOf(in)
	outV := reflect.ValueOf(out)
	if inV.Kind() == reflect.Ptr {
		inV = inV.Elem()
	}
	if outV.Kind() == reflect.Ptr {
		outV = outV.Elem()
	}
	return sliceStructToMapStructValue(inV, outV, idFieldName)
}
func sliceStructToMapStructValue(in reflect.Value, out reflect.Value, idFieldName string) (err error) {
	out.Set(reflect.MakeMap(out.Type()))
	len := in.Len()
	isOutPtr := out.Type().Elem().Kind() == reflect.Ptr
	for i := 0; i < len; i++ {
		thisValue := in.Index(i)
		oKey := thisValue.FieldByName(idFieldName)
		if !oKey.IsValid() {
			return fmt.Errorf(`id field name "%s" not exist in "%s"`, idFieldName, thisValue.Type().Name())
		}
		if isOutPtr {
			thisValue = thisValue.Addr()
		}
		oExist := out.MapIndex(oKey)
		if oExist.IsValid() {
			return fmt.Errorf(`%s:%v repeat`, idFieldName, oKey.Interface())
		}
		out.SetMapIndex(oKey, thisValue)
	}
	return
}
