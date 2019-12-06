package udwShm

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwBufio"
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

func Unmarshal(data []byte, obj interface{}) (errMsg string) {
	r := NewShmReader(udwBytes.NewBufReader(data), 0)
	return r.ReadValue(obj)
}

type ShmReader struct {
	r                udwBufio.BufioReader
	isLastEndOfArray bool
}

func NewShmReader(r io.Reader, maxStringSize int) *ShmReader {
	return &ShmReader{
		r: *udwBufio.NewBufioReader(r, maxStringSize),
	}
}

func (r *ShmReader) ReadOneString() (b []byte, errMsg string) {
	x, errMsg := r.r.ReadUvarint()
	if errMsg != "" {
		r.isLastEndOfArray = false
		return nil, errMsg
	}
	if x == 0 {
		r.isLastEndOfArray = true
		return nil, "v747qhnfzu"
	}
	if x == 1 {
		return nil, "3a64xt8n9t"
	}
	r.isLastEndOfArray = false
	sl := x - 2
	b, errMsg = r.r.ReadBySize(int(sl))
	if errMsg != "" {
		return nil, errMsg
	}
	return b, ""
}

func (r *ShmReader) ReadValue(a interface{}) (errMsg string) {
	v := reflect.ValueOf(a)

	if v.Kind() != reflect.Ptr {
		return "qdyrnapykr"
	}
	value := v
	for {
		if value.CanSet() {
			break
		}
		if value.IsNil() {
			return "qdyrnapykr"
		}
		value = value.Elem()
		if value.Kind() != reflect.Ptr {
			break
		}
	}
	return unpackValue(r, v)
}

func (r *ShmReader) ReadArrayStart() (errMsg string) {
	b, errMsg := r.r.ReadByteErrMsg()
	r.isLastEndOfArray = false
	if errMsg != "" {
		return errMsg
	}
	if b != 1 {
		return "dcxgumhkfd"
	}
	return ""
}

func (r *ShmReader) ReadArrayEnd() (errMsg string) {
	b, errMsg := r.r.ReadByteErrMsg()
	r.isLastEndOfArray = false
	if errMsg != "" {
		return errMsg
	}
	if b != 0 {
		return "t83pamau27"
	}
	return ""
}

type ReadUvarintResp struct {
	IsArrayEnd bool
	B          []byte
	ErrMsg     string
}

func (r *ShmReader) ReadUvarint() (v uint64, errMsg string) {
	b, errMsg := r.ReadOneString()
	if errMsg != "" {
		return 0, errMsg
	}
	v, size := udwBytes.ReadUvarint(b)
	if size <= 0 {
		return 0, "624avkdeme"
	}
	if size != len(b) {
		return 0, "jrrzthx55p"
	}
	return v, ""
}

func unpackValue(r *ShmReader, value reflect.Value) (errMsg string) {

	switch value.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		u64, errMsg := r.ReadUvarint()
		if errMsg != "" {
			return errMsg
		}
		value.SetInt(int64(u64))
		return ""
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		u64, errMsg := r.ReadUvarint()
		if errMsg != "" {
			return errMsg
		}
		value.SetUint(uint64(u64))
		return ""
	case reflect.Map:
		b, errMsg := r.r.ReadByteErrMsg()
		r.isLastEndOfArray = false
		if errMsg != "" {
			return errMsg
		}
		if b != 1 {
			return "dcxgumhkfd"
		}
		value.Set(reflect.MakeMap(value.Type()))
		for {
			key := reflect.New(value.Type().Key())
			errMsg = unpackValue(r, key)
			if errMsg != `` {
				if r.isLastEndOfArray {
					return ""
				}
				return errMsg
			}
			value1 := reflect.New(value.Type().Elem())
			errMsg = unpackValue(r, value1)
			if errMsg != `` {
				return errMsg
			}
			value.SetMapIndex(key.Elem(), value1.Elem())
		}
	case reflect.Array, reflect.Slice:
		isSlice := (value.Kind() == reflect.Slice)
		if value.Type().Elem().Kind() == reflect.Uint8 {
			byteList, errMsg := r.ReadOneString()
			if errMsg != "" {
				return errMsg
			}
			if isSlice {
				value.SetBytes(udwBytes.Clone(byteList))
				return ""
			} else {
				if len(byteList) > value.Len() {
					return "8c3gd2w9wm " + strconv.Itoa(len(byteList)) + " " + strconv.Itoa(value.Len())
				}
				for i, b := range byteList {
					value.Index(i).SetUint(uint64(b))
				}
			}
			return ""
		}
		b, errMsg := r.r.ReadByteErrMsg()
		r.isLastEndOfArray = false
		if errMsg != "" {
			return errMsg
		}
		if b != 1 {
			return "jnbtfp9373"
		}
		i := 0
		for {
			key := reflect.New(value.Type().Elem())
			errMsg = unpackValue(r, key)
			if errMsg != `` {
				if r.isLastEndOfArray {
					return ""
				}
				return errMsg
			}
			if isSlice {
				value.Set(reflect.Append(value, key.Elem()))
			} else {
				if i >= value.Len() {
					return "p98uksqr58"
				}
				value.Index(i).Set(key.Elem())
				i++
			}
		}
	case reflect.String:
		b, errMsg := r.ReadOneString()
		if errMsg != "" {
			return errMsg
		}
		value.SetString(string(b))
		return ""
	case reflect.Ptr:
		b, errMsg := r.r.PeekByte()
		r.isLastEndOfArray = false
		if errMsg != "" {
			return errMsg
		}
		if b == 0 {
			r.isLastEndOfArray = true
			_, errMsg = r.r.ReadByteErrMsg()
			return "skmm76yukh"
		}
		if b == 2 {
			_, errMsg = r.r.ReadByteErrMsg()
			return errMsg
		}
		for {
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
			if value.Kind() != reflect.Ptr {
				break
			}
		}
		return unpackValue(r, value)
	case reflect.Struct:
		_, ok := value.Interface().(time.Time)
		if ok {
			v, errMsg := r.ReadUvarint()
			if errMsg != "" {
				return errMsg
			}
			t := time.Unix(0, int64(v))
			value.Set(reflect.ValueOf(t))
			return ""
		}
		b, errMsg := r.r.ReadByteErrMsg()
		r.isLastEndOfArray = false
		if errMsg != "" {
			return errMsg
		}
		if b != 1 {
			return "3kuymxwdxu"
		}
		for {
			nameB, errMsg := r.ReadOneString()
			if errMsg != "" {
				if r.isLastEndOfArray {
					return ""
				}
				return errMsg
			}
			name := string(nameB)

			ft, ok := value.Type().FieldByName(name)
			if len(ft.Index) > 1 {
				return "b8u9jrngas unsupported type"
			}
			if ok == false || shouldSkipFiled(ft) {
				errMsg := unpackSkip(r)
				if errMsg != "" {
					return errMsg
				}
				continue
			}
			f := value.FieldByName(name)
			errMsg = unpackValue(r, f)
			if errMsg != "" {
				return errMsg
			}
		}
	case reflect.Bool:
		b, errMsg := r.ReadOneString()
		if errMsg != "" {
			return errMsg
		}
		if len(b) != 1 {
			return "uw5c3trah5"
		}
		if b[0] == 0 {
			value.SetBool(false)
		} else if b[0] == 1 {
			value.SetBool(true)
		} else {
			return "8eh2khhc25"
		}
	case reflect.Float64:
		b, errMsg := r.ReadOneString()
		if errMsg != "" {
			return errMsg
		}
		if len(b) != 8 {
			return "h24hx3krs4"
		}
		v := binary.LittleEndian.Uint64(b)
		f := math.Float64frombits(v)
		value.SetFloat(f)
		return ""
	case reflect.Float32:
		b, errMsg := r.ReadOneString()
		if errMsg != "" {
			return errMsg
		}
		if len(b) != 4 {
			return "h24hx3krs4"
		}
		v := binary.LittleEndian.Uint32(b)
		f := math.Float32frombits(v)
		value.SetFloat(float64(f))
		return ""
	default:
		return "4fweebt5a3 unsupported type"
	}
	return
}

func unpackSkip(r *ShmReader) (errMsg string) {
	x, errMsg := r.r.ReadUvarint()
	if errMsg != "" {
		return errMsg
	}
	if x == 0 {
		r.isLastEndOfArray = true
		return "nw5ec2wrxc"
	}
	if x == 1 {
		for {
			errMsg = unpackSkip(r)
			if errMsg != "" {
				if r.isLastEndOfArray {
					return ""
				}
				return errMsg
			}
		}
	}
	r.isLastEndOfArray = false
	sl := x - 2
	_, errMsg = r.r.ReadBySize(int(sl))
	if errMsg != "" {
		return errMsg
	}
	return ""
}
