package udwSort

import (
	"reflect"
)

func InterfaceRevert(objList interface{}) {
	objListReflect := reflect.ValueOf(objList)
	if objListReflect.Kind() != reflect.Slice {
		panic("[InterfaceRevert] only support slice" + objListReflect.Kind().String())
	}
	sliceLen := objListReflect.Len()
	loopLen := sliceLen / 2
	objListElemType := reflect.TypeOf(objList).Elem()
	for i := 0; i < loopLen; i++ {
		a := i
		b := sliceLen - i - 1

		tmp := reflect.New(objListElemType).Elem()
		tmp.Set(objListReflect.Index(b))
		objListReflect.Index(b).Set(objListReflect.Index(a))
		objListReflect.Index(a).Set(tmp)
	}
}
