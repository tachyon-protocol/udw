package udwEncodeBinary

import (
	"reflect"
	"github.com/tachyon-protocol/udwBytes"
)

func Unmarshal(b []byte,obj interface{}) (errMsg string){
	r:=udwBytes.NewBufReader(b)
	value:=reflect.ValueOf(obj)
	if value.Kind()!=reflect.Ptr{
		return "bz5zz45utu"
	}
	return unmarshalFromReader(value.Elem(),r)
}

func unmarshalFromReader(value reflect.Value,r *udwBytes.BufReader) (errMsg string){
	if value.CanSet()==false{
		return "jtf5u7q66h"
	}
	kind:=value.Kind()
	switch kind {
	case reflect.String:
		s,isOk:=r.ReadStringLenUvarint()
		if isOk==false{
			return "vyypaf2r4w"
		}
		value.SetString(s)
		return ""
	case reflect.Struct:
		fieldNum:=value.NumField()
		for i:=0;i<fieldNum;i++{
			errMsg = unmarshalFromReader(value.Field(i),r)
			if errMsg!=""{
				return errMsg
			}
		}
		return ""
	default:
		return "ztnnxat46u "+kind.String()
	}
}