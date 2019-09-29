package udwType

import (
	"reflect"
)

type StringType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *StringType) ToString(v reflect.Value) string {
	return v.String()
}
func (t *StringType) SaveScale(v reflect.Value, value string) error {
	v.SetString(value)
	return nil
}
