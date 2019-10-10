package udwJson

import (
	"encoding/json"
	"fmt"
	"github.com/tachyon-protocol/udw/udwJson/udwJsonLib"
	"github.com/tachyon-protocol/udw/udwStrings"
	"reflect"
)

func MustUnmarshalFromStringToStringSlice(s string) (StringList []string) {
	if s == "" {
		return []string{}
	}
	ctx := udwJsonLib.NewContextFromBuffer(udwStrings.GetByteArrayFromStringNoAlloc(s))
	udwJsonLib.ReaderReadSpace(ctx)
	b := udwJsonLib.ReaderReadByte(ctx)
	switch b {
	case 'n':
		udwJsonLib.ReaderReadBack(ctx, 1)
		udwJsonLib.MustReadJsonNull(ctx)
		StringList = nil
	case '[':
		StringList = []string{}
		for {
			b = udwJsonLib.ReaderReadByte(ctx)
			if b == ',' || b == ' ' || b == '\t' || b == '\n' || b == '\r' {
				continue
			} else if b == ']' {
				break
			}
			udwJsonLib.ReaderReadBack(ctx, 1)
			_var11 := udwJsonLib.ReadJsonString(ctx)
			StringList = append(StringList, _var11)
		}
	default:
		panic("need a [ or null but get [" + string(b) + "]")
	}
	return StringList
}

func MustUnmarshalFromStringToStringSliceByIndex(key string, index int) (outS string) {
	kl := MustUnmarshalFromStringToStringSlice(key)
	return kl[index]
}

func MustMarshalKeyPrefix(k ...string) string {
	if len(k) == 0 {
		return ""
	}
	outB := MustMarshalStringArgumentToString(k...)
	return outB[:len(outB)-2]
}

func MustMarshalKeyPrefixWithEnd(k ...string) string {
	if len(k) == 0 {
		return ""
	}
	outB := MustMarshalStringArgumentToString(k...)
	return outB[:len(outB)-1]
}

func MustMarshalStringArgumentToString(sList ...string) string {

	ctx := udwJsonLib.PoolGetContextWithWriteBuf()
	if sList != nil {
		udwJsonLib.WriterWriteString(ctx, `[`)
		for i := 0; i < len(sList); i++ {
			_var13 := sList[i]
			udwJsonLib.WriteJsonString(ctx, _var13)
			if i < len(sList)-1 {
				udwJsonLib.WriterWriteByte(ctx, ',')
			}
		}
		udwJsonLib.WriterWriteByte(ctx, ']')
	}
	out := string(ctx.WriterBytes())
	udwJsonLib.PoolPutContextWithWriteBuf(ctx)
	return out
}

func StringSliceToObjSlice(outList []string, obj interface{}) (err error) {
	return stringSliceToObjSliceL1(outList, reflect.ValueOf(obj))
}

func stringSliceToObjSliceL1(outList []string, obj reflect.Value) (err error) {
	switch obj.Kind() {
	case reflect.Ptr:
		return stringSliceToObjSliceL1(outList, obj.Elem())
	case reflect.Slice:
		newSlice := reflect.MakeSlice(obj.Type(), len(outList), len(outList))
		elemType := obj.Type().Elem()
		for i, s := range outList {
			thisValue := newSlice.Index(i)
			thisElem := reflect.New(elemType)
			err = json.Unmarshal([]byte(s), thisElem.Interface())
			if err != nil {
				return err
			}
			thisValue.Set(thisElem.Elem())
		}
		obj.Set(newSlice)
		return nil
	default:
		return fmt.Errorf("[mgetNotExistCheckGobUnmarshal] Unmarshal unexpect Kind %s", obj.Kind().String())
	}
}
