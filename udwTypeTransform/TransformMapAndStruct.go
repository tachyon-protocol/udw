package udwTypeTransform

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwReflect"
	"reflect"
)

const debugMapToStruct = false

func MapToStruct(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	outT := out.Type()
	out.Set(reflect.New(outT).Elem())
	fieldNameMap := map[string]string{}
	fieldList := udwReflect.StructGetAllField(outT)
	for _, field := range fieldList {
		if field.Tag == "" {
			fieldNameMap[field.Name] = field.Name
			continue
		}
		mapKey := field.Tag.Get("typeTransform")
		if mapKey == "" {
			fieldNameMap[field.Name] = field.Name
			continue
		}
		fieldNameMap[mapKey] = field.Name
	}

	for _, key := range in.MapKeys() {
		err = t.Tran(key, oKey)
		if err != nil {
			return
		}
		mapKeyS := oKey.String()

		structKeyS := fieldNameMap[mapKeyS]
		if debugMapToStruct {
			fmt.Println("debug", "[MapToStruct]", mapKeyS, structKeyS)
		}
		oVal := out.FieldByName(structKeyS)
		if !oVal.IsValid() {
			continue
		}
		val := in.MapIndex(key)
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
	}
	return
}

func StructToMap(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	oValType := out.Type().Elem()
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	if out.IsNil() {
		out.Set(reflect.MakeMap(out.Type()))
	}

	fieldMap := udwReflect.StructGetAllFieldMap(in.Type())
	for key, field := range fieldMap {
		if field.PkgPath != "" {

			continue
		}
		iVal := in.FieldByName(key)
		oVal := reflect.New(oValType).Elem()
		err = t.Tran(iVal, oVal)
		if err != nil {
			return
		}
		oKey.SetString(key)
		out.SetMapIndex(oKey, oVal)
	}
	return nil
}
