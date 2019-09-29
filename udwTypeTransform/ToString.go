package udwTypeTransform

import "reflect"

func ToString(in interface{}) (out string, err error) {
	err = getDefaultTransformer().Transform(in, &out)
	return
}

func ToStringReflect(in reflect.Value) (out string, err error) {
	err = getDefaultTransformer().Tran(in, reflect.ValueOf(&out))
	return
}
