package udwType

import (
	"reflect"
	"strconv"
)

type BoolType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *BoolType) ToString(v reflect.Value) string {
	return strconv.FormatBool(v.Bool())
}
func (t *BoolType) SaveScale(v reflect.Value, value string) error {
	normV, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.SetBool(normV)
	return nil
}
