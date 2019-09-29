package udwType

import (
	"reflect"
	"strconv"
)

type FloatType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *FloatType) ToString(v reflect.Value) string {
	return strconv.FormatFloat(v.Float(), 'g', -1, t.GetReflectType().Bits())
}
func (t *FloatType) SaveScale(v reflect.Value, value string) error {
	f, err := strconv.ParseFloat(value, t.GetReflectType().Bits())
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}
