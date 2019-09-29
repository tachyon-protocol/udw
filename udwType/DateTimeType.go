package udwType

import (
	"github.com/tachyon-protocol/udw/udwTime"
	"reflect"
	"time"
)

func GetDateTimeReflectType() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil)).Elem()
}

type DateTimeType struct {
	reflectTypeGetterImp
	saveScaleFromStringer
	saveScaleEditabler
}

func (t *DateTimeType) ToString(v reflect.Value) string {
	return v.Interface().(time.Time).Format(udwTime.FormatMysql)
}
func (t *DateTimeType) SaveScale(v reflect.Value, value string) error {
	valueT, err := time.Parse(udwTime.FormatMysql, value)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(valueT))
	return nil
}
