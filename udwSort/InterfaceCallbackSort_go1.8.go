package udwSort

import (
	"reflect"
	"sort"
	"strconv"
)

func InterfaceCallbackSort(objList interface{}, lessFn interface{}) {

	lessFnReflect := reflect.ValueOf(lessFn)
	if lessFnReflect.Kind() != reflect.Func {
		panic("[InterfaceCallbackSort] lessFn must be func,but get kind [" + lessFnReflect.Kind().String() + "]")
	}
	objListReflect := reflect.ValueOf(objList)
	kind := objListReflect.Kind()
	lessFnType := reflect.TypeOf(lessFn)
	if lessFnType.NumIn() != 2 {
		panic("[InterfaceCallbackSort] lessFn must have two in parameter,but get kind [" + strconv.Itoa(lessFnType.NumIn()) + "]")
	}
	if lessFnType.NumOut() != 1 {
		panic("[InterfaceCallbackSort] lessFn must have one out parameter,but get kind [" + strconv.Itoa(lessFnType.NumOut()) + "]")
	}
	objListElemType := reflect.TypeOf(objList).Elem()
	if !lessFnType.In(0).AssignableTo(objListElemType) {
		panic("[InterfaceCallbackSort] lessFn first in parameter must have the same type of objList Element,but get [" + lessFnType.In(0).String() + "]")
	}
	if !lessFnType.In(1).AssignableTo(objListElemType) {
		panic("[InterfaceCallbackSort] lessFn second in parameter must have the same type of objList Element,but get [" + lessFnType.In(1).String() + "]")
	}
	if !lessFnType.Out(0).AssignableTo(reflect.TypeOf(false)) {
		panic("[InterfaceCallbackSort] lessFn first out parameter must be bool,but get [" + lessFnType.Out(0).String() + "]")
	}
	if kind == reflect.Slice {
		sort.Slice(objList, func(a int, b int) bool {
			outList := lessFnReflect.Call([]reflect.Value{objListReflect.Index(a), objListReflect.Index(b)})
			return outList[0].Bool()
		})
	} else {
		panic("[InterfaceCallbackSort] unsupport objList kind [" + kind.String() + "]")
	}
}

func InterfaceCallbackSortWithIndexLess(objList interface{}, lessFn func(a int, b int) bool) {
	sort.Slice(objList, lessFn)
}
