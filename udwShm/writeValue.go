package udwShm

import (
	"github.com/tachyon-protocol/udw/udwBufio"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwStrings"
	"io"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func Marshal(obj interface{}) (data []byte, errMsg string) {
	_buf := udwBytes.BufWriter{}
	w := NewShmWriter(&_buf, 0)
	errMsg = w.WriteValue(obj)
	if errMsg != "" {
		return nil, errMsg
	}
	errMsg = w.Flush()
	if errMsg != "" {
		return nil, errMsg
	}
	return _buf.GetBytes(), ""
}

type ShmWriter struct {
	w udwBufio.BufioWriter
}

func NewShmWriter(w io.Writer, softMaxBufferSize int) *ShmWriter {
	return &ShmWriter{
		w: *udwBufio.NewBufioWriter(w, softMaxBufferSize),
	}
}

func (w *ShmWriter) WriteValue(a interface{}) (errMsg string) {
	errMsg = packValue(w, reflect.ValueOf(a))
	if errMsg != "" {
		return errMsg
	}
	return ""
}

func (w *ShmWriter) WriteByteSlice(buf []byte) {
	w.w.WriteUvarint(uint64(len(buf)) + 2)
	w.w.Write_(buf)
	return
}

func (w *ShmWriter) WriteString(s string) {
	w.WriteByteSlice([]byte(s))
}

func (w *ShmWriter) WriteArrayStart() {
	w.w.WriteByte_(1)
	return
}

func (w *ShmWriter) WriteArrayEnd() {
	w.w.WriteByte_(0)
	return
}

func packValue(w *ShmWriter, value reflect.Value) (errMsg string) {
	switch value.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		w.WriteUvarint(uint64(value.Int()))
		return ""
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		w.WriteUvarint(uint64(value.Uint()))
		return ""
	case reflect.Map:
		w.WriteArrayStart()
		for _, key := range value.MapKeys() {
			errMsg = packValue(w, key)
			if errMsg != `` {
				return errMsg
			}
			errMsg = packValue(w, value.MapIndex(key))
			if errMsg != `` {
				return errMsg
			}
		}
		w.WriteArrayEnd()
		return ""
	case reflect.Array, reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			if value.Kind() == reflect.Slice {
				bs := value.Bytes()
				w.WriteByteSlice(bs)
				return ""
			} else {
				bs := make([]byte, value.Len())
				for i := 0; i < value.Len(); i++ {
					bs[i] = uint8(value.Index(i).Uint())
				}
				w.WriteByteSlice(bs)
				return ""
			}
		}
		w.WriteArrayStart()
		for i := 0; i < value.Len(); i++ {
			errMsg = packValue(w, value.Index(i))
			if errMsg != `` {
				return errMsg
			}
		}
		w.WriteArrayEnd()
		return ""
	case reflect.String:
		w.WriteString(value.String())
		return ""
	case reflect.Ptr:
		if value.IsNil() {
			w.WriteByteSlice(nil)
			return ""
		}
		for {
			value = value.Elem()
			if value.Kind() != reflect.Ptr {
				break
			}
		}
		return packValue(w, value)
	case reflect.Struct:
		t, ok := value.Interface().(time.Time)
		if ok {
			intValue := t.UnixNano()
			w.WriteUvarint(uint64(intValue))
			return ""
		}
		w.WriteArrayStart()
		for i := 0; i < value.NumField(); i++ {
			f := value.Field(i)
			ft := value.Type().Field(i)
			if shouldSkipFiled(ft) || isZeroValue(f) {
				continue
			}
			if ft.Anonymous {
				return `Unsupported Anonymous`
			}
			w.WriteString(ft.Name)
			errMsg = packValue(w, f)
			if errMsg != `` {
				return errMsg
			}
		}
		w.WriteArrayEnd()
		return ""
	case reflect.Bool:
		if value.Bool() {
			w.WriteByteSlice([]byte{1})
			return ""
		} else {
			w.WriteByteSlice([]byte{0})
			return ""
		}
	case reflect.Float64:
		w.w.WriteByte_(8 + 2)
		w.w.WriteLittleEndFloat64(value.Float())
		return ""
	case reflect.Float32:
		w.w.WriteByte_(4 + 2)
		w.w.WriteLittleEndFloat32(float32(value.Float()))
		return ""
	default:
		return `Unsupported kind ` + value.Kind().String()
	}
	return ""
}

func (w *ShmWriter) WriteUvarint(v uint64) {
	w.w.AddPos(1)
	startPos := w.w.GetPos()
	w.w.WriteUvarint(uint64(v))
	endPos := w.w.GetPos()
	size := endPos - startPos
	w.w.SetPos(startPos - 1)
	w.w.WriteByte_(uint8(size + 2))
	w.w.SetPos(endPos)
	return
}

func (w *ShmWriter) Flush() (errMsg string) {
	return w.w.Flush()
}

func shouldSkipFiled(ft reflect.StructField) bool {
	if udwStrings.IsInSlice(strings.Split(ft.Tag.Get(`json`), `,`), `-`) {
		return true
	}
	ch, _ := utf8.DecodeRuneInString(ft.Name)
	return !unicode.IsUpper(ch)
}

func isZeroValue(v reflect.Value) bool {
	t, ok := v.Interface().(time.Time)
	if ok && t.IsZero() {
		return true
	}
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Array:
		if v.Len() == 0 {
			return true
		}
		for i := 0; i < v.Len(); i++ {
			if isZeroValue(v.Index(i)) == false {
				return false
			}
		}
		return true
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float64, reflect.Float32:
		return v.Float() == 0
	case reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		t, ok := v.Interface().(time.Time)
		if ok {
			return t.IsZero()
		}
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			ft := v.Type().Field(i)
			isSkip := shouldSkipFiled(ft) || isZeroValue(f)
			if isSkip == false {
				return false
			}
		}
		return true
	}
	return false
}
