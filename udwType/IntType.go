package udwType

import (
	"reflect"
	"strconv"
)

type IntType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *IntType) ToString(v reflect.Value) string {
	s := strconv.FormatInt(v.Int(), 10)
	return s
}
func (t *IntType) SaveScale(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, t.GetReflectType().Bits())
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}
