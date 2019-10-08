package udwTypeTransform

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwReflect"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTime"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Transform(in interface{}, out interface{}) (err error) {
	return getDefaultTransformer().Transform(in, out)
}

func MustTransform(in interface{}, out interface{}) {
	err := getDefaultTransformer().Transform(in, out)
	if err != nil {
		panic(err)
	}
}

func MustTransformToMap(in interface{}) (m map[string]string) {
	m = map[string]string{}
	err := getDefaultTransformer().Transform(in, &m)
	if err != nil {
		panic(err)
	}
	return m
}

func BoolToBool(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetBool(in.Bool())
	return nil
}

func MapToMap(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.Set(reflect.MakeMap(out.Type()))
	for _, key := range in.MapKeys() {
		oKey := reflect.New(out.Type().Key()).Elem()
		oVal := reflect.New(out.Type().Elem()).Elem()
		err = t.Tran(key, oKey)
		if err != nil {
			return
		}
		val := in.MapIndex(key)
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
		out.SetMapIndex(oKey, oVal)
	}
	return
}

func StringToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetString(in.String())
	return nil
}

func NewStringToTimeFunc(location *time.Location) TransformerFunc {
	return func(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
		var t time.Time
		t, err = udwTime.ParseAutoInLocation(in.String(), location)
		if err != nil {
			return
		}
		out.Set(reflect.ValueOf(t))
		return
	}
}

func StringToTime(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
	var t time.Time
	t, err = udwTime.ParseAutoInLocal(in.String())
	if err != nil {
		return
	}
	out.Set(reflect.ValueOf(t))
	return
}

func TimeToString(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
	t := in.Interface().(time.Time)
	out.SetString(t.In(udwTime.GetDefaultTimeZone()).Format(udwTime.FormatMysql))
	return nil
}

func PtrToPtr(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	return t.Tran(in.Elem(), out.Elem())
}

func StructToStruct(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.Set(reflect.New(out.Type()).Elem())
	fieldMap := udwReflect.StructGetAllFieldMap(in.Type())
	for key, field := range fieldMap {
		if field.PkgPath != "" {

			continue
		}
		iVal := in.FieldByName(key)
		oVal := out.FieldByName(key)
		if !oVal.IsValid() {
			continue
		}
		err = t.Tran(iVal, oVal)
		if err != nil {
			return
		}
	}
	return nil
}

func SliceToSlice(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	len := in.Len()
	out.Set(reflect.MakeSlice(out.Type(), len, len))
	for i := 0; i < len; i++ {
		val := in.Index(i)
		err = t.Tran(val, out.Index(i))
		if err != nil {
			return
		}
	}
	return
}

func StringToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	inS = strings.TrimSpace(inS)
	if inS == "" {
		out.SetInt(int64(0))
		return nil
	}
	i, err := strconv.ParseInt(inS, 10, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetInt(i)
	return
}

func StringToUint(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	inS = strings.TrimSpace(inS)
	if inS == "" {
		out.SetUint(uint64(0))
		return nil
	}
	i, err := strconv.ParseUint(inS, 10, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetUint(i)
	return
}

func StringToFloat(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetFloat(0.0)
		return nil
	}
	i, err := strconv.ParseFloat(inS, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetFloat(i)
	return
}

func StringToBool(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetBool(false)
		return nil
	}
	i, err := strconv.ParseBool(inS)
	if err != nil {
		return
	}
	out.SetBool(i)
	return
}

func IntToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetInt(in.Int())
	return nil
}

func FloatToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	outf1 := in.Float()
	if math.Floor(outf1) != outf1 {
		return fmt.Errorf("[typeTransform.tran] it seems to lose some accuracy trying to convert from float to int,float:%f", outf1)
	}
	out.SetInt(int64(outf1))
	return
}

func FloatToFloat(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetFloat(in.Float())
	return
}

func FloatToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	f := in.Float()
	fs := udwStrconv.FormatFloat(f)
	out.SetString(fs)
	return
}

func NonePtrToPtr(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	if out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
	}
	return t.Tran(in, out.Elem())
}
func InterfaceToNoneInterface(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	return t.Tran(in.Elem(), out)
}
func NoneInterfaceToInterface(t Transformer, in reflect.Value, out reflect.Value) (err error) {

	if in.Type().Implements(out.Type()) {
		out.Set(in)
		return
	}
	return t.Tran(in, out.Elem())
}

func IntToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	s := strconv.FormatInt(in.Int(), 10)
	out.SetString(s)
	return nil
}

func FuncToFunc(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.Set(in)
	return nil
}
