package udwEncodeBinary

import (
	"reflect"
	"github.com/tachyon-protocol/udwBytes"
)

func Marshal(obj interface{}) (b []byte,errMsg string){
	_buf:=&udwBytes.BufWriter{}
	value:=reflect.ValueOf(obj)
	errMsg = marshalToBufL1(value,_buf)
	if errMsg!=""{
		return nil,errMsg
	}
	return _buf.GetBytes(),""
}

func marshalToBufL1(value reflect.Value,_buf *udwBytes.BufWriter) (errMsg string){
	kind:=value.Kind()
	switch kind {
	case reflect.String:
		_buf.WriteStringLenUvarint(value.String())
		return ""
	case reflect.Struct:
		fieldNum:=value.NumField()
		for i:=0;i<fieldNum;i++{
			errMsg = marshalToBufL1(value.Field(i),_buf)
			if errMsg!=""{
				return errMsg
			}
		}
		return ""
	default:
		return "mkxvg22qn7 "+kind.String()
	}
}