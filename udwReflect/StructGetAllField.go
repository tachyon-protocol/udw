package udwReflect

import (
	"reflect"
)

func StructGetAllField(t reflect.Type) (output []*reflect.StructField) {
	fieldMap := map[string]bool{}
	structGetAllFieldImpWithCallback(t, []int{}, fieldMap, func(field *reflect.StructField) {
		output = append(output, field)
	})
	return output
}

func structGetAllFieldImpWithCallback(t reflect.Type, indexs []int, fieldMap map[string]bool, fieldCallback func(field *reflect.StructField)) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	anonymousFieldList := []*reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		f.Index = append(indexs, f.Index...)
		if f.Anonymous {
			anonymousFieldList = append(anonymousFieldList, &f)
		}
		if fieldMap[f.Name] {
			continue
		}
		fieldMap[f.Name] = true
		fieldCallback(&f)

	}
	for _, f := range anonymousFieldList {
		structGetAllFieldImpWithCallback(f.Type, f.Index, fieldMap, fieldCallback)

	}
	return
}

func StructGetAllFieldMap(t reflect.Type) (output map[string]*reflect.StructField) {
	fieldMap := map[string]bool{}
	output = map[string]*reflect.StructField{}
	structGetAllFieldImpWithCallback(t, []int{}, fieldMap, func(field *reflect.StructField) {
		output[field.Name] = field
	})
	return output

}
