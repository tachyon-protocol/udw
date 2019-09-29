package udwTypeTransform

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwReflect"
	"reflect"
)

func StringTransformSubType(in interface{}, transformTable map[string]map[string]string) (err error) {
	return stringTransformSubType(reflect.ValueOf(in), transformTable)
}
func stringTransformSubType(in reflect.Value, transformTable map[string]map[string]string) (err error) {
	switch in.Kind() {
	case reflect.String:
		typeName := udwReflect.GetTypeFullName(in.Type())
		thisTable, exist := transformTable[typeName]
		if !exist {
			return
		}
		inS := in.String()
		if inS == "" {
			in.SetString("")
			return
		}
		oVal, exist := thisTable[inS]
		if !exist {
			return fmt.Errorf(`string transform fail! from:"%s",type:"%s"`, inS, typeName)
		}
		in.SetString(oVal)
		return
	case reflect.Ptr:
		return stringTransformSubType(in.Elem(), transformTable)
	case reflect.Slice:
		len := in.Len()
		for i := 0; i < len; i++ {
			err = stringTransformSubType(in.Index(i), transformTable)
			if err != nil {
				return
			}
		}
	case reflect.Struct:
		len := in.Type().NumField()
		for i := 0; i < len; i++ {
			err = stringTransformSubType(in.Field(i), transformTable)
			if err != nil {
				return
			}
		}
	}
	return
}
