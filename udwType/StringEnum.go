package udwType

import (
	"fmt"
	"reflect"
)

type StringEnum interface {
	GetEnumList() []string
}

func IsEnumExist(enum StringEnum) bool {
	v := reflect.ValueOf(enum)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("[IsEnumExist] you should pass in a type which underlying type is string ,Get:%s", v.Kind()))
	}
	return IsEnumExistString(enum, v.String())
}

func IsEnumExistString(enum StringEnum, s string) bool {
	for _, enumItem := range enum.GetEnumList() {
		if s == enumItem {
			return true
		}
	}
	return false
}
