package udwTest

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

func Ok(expectTrue bool, objList ...interface{}) {
	if expectTrue == true {
		return
	}
	if len(objList) == 0 {
		panic("ok fail")
	}
	_buf := bytes.Buffer{}
	_buf.WriteString("ok fail\n")
	for i := range objList {
		f, ok := objList[i].(func() string)
		if ok {
			_buf.WriteString(f())
			continue
		}
		_buf.WriteString(sprintln(objList[i]))
	}
	panic(_buf.String())
}

func Equal(get interface{}, expect interface{}, objList ...interface{}) {

	switch objA := get.(type) {
	case []byte:
		objB, ok := expect.([]byte)
		if ok && bytes.Equal(objA, objB) {
			return
		}
	case []string:
		objB, ok := expect.([]string)
		if ok && len(objA) == len(objB) {
			isAllSame := true
			for i := range objA {
				if objA[i] != objB[i] {
					isAllSame = false
					break
				}
			}
			if isAllSame {
				return
			}
		}
	case string:
		objB, ok := expect.(string)
		if ok && objA == objB {
			return
		}
	case int:
		objB, ok := expect.(int)
		if ok && objA == objB {
			return
		}
	case bool:
		objB, ok := expect.(bool)
		if ok && objA == objB {
			return
		}
	case nil:
		if expect == nil {
			return
		}
	}

	interfaceData := *(*[2]uintptr)(unsafe.Pointer(&get))
	get2 := *(*(interface{}))(unsafe.Pointer(&interfaceData))

	interfaceData2 := *(*[2]uintptr)(unsafe.Pointer(&expect))
	expect2 := *(*(interface{}))(unsafe.Pointer(&interfaceData2))

	if isEqual(expect2, get2) {
		return
	}

	var msg string
	msg = "\tget1: " + valueDetail(get2) + "\n\texpect2: " + valueDetail(expect2) + "\n"
	switch objA := get.(type) {
	case string:
		objB, ok := expect.(string)
		if ok {
			diffPos := findStringFirstDiffPos(objA, objB)
			msg += "diff at 0x" + strconv.FormatUint(uint64(diffPos), 16) + " " + strconv.Itoa(int(diffPos)) + "\n"
		}
	case []byte:
		objB, ok := expect.([]byte)
		if ok {
			diffPos := findStringFirstDiffPos(string(objA), string(objB))
			msg += "diff at 0x" + strconv.FormatUint(uint64(diffPos), 16) + " " + strconv.Itoa(int(diffPos)) + "\n"
		}
	}
	if len(objList) > 0 {
		msg += sprintln(objList...)
	}
	fmt.Println(msg)
	panic("[udwTest.Equal] fail")
}

func isEqual(a interface{}, b interface{}) bool {

	if reflect.DeepEqual(a, b) {
		return true
	}
	rva := reflect.ValueOf(a)
	rvb := reflect.ValueOf(b)

	if isNil(rva) && isNil(rvb) {
		return true
	}
	return false
}

type assertPanicType struct{}

func AssertPanic(f func()) (out interface{}) {
	defer func() {
		out = recover()
		_, ok := out.(assertPanicType)
		if ok {
			panic("should panic")
		}
	}()
	f()
	panic(assertPanicType{})
}

func AssertPanicWithErrorMessage(f func(), errorMsgList ...string) {
	outI := AssertPanic(f)
	msg := ""
	switch out := outI.(type) {
	case error:
		msg = out.Error()
		return
	case string:
		msg = out
	default:
		panic(fmt.Errorf("[AssertPanicWithErrorMessage] not expect panic type %T", outI))
	}
	hasFound := false
	for _, errMsg := range errorMsgList {
		if strings.Contains(msg, errMsg) {
			hasFound = true
			break
		}
	}
	if hasFound == false {
		panic("[AssertPanicWithErrorMessage] panic:[" + msg + "] need:[" + strings.Join(errorMsgList, ", ") + "]")
	}
}

func valueDetail(value interface{}) string {
	switch obj := value.(type) {
	case tStringer:
		return fmt.Sprintf("%s (%T) %#v", obj.String(), value, value)
	case error:
		return fmt.Sprintf("%s (%T)", obj.Error(), value)
	case string:
		return fmt.Sprintf("%s (%T)\n [%s]", strconv.Quote(obj), value, hex.Dump([]byte(obj)))
	case []byte:
		return fmt.Sprintf("%s (%T)\n%s\n%s", strconv.Quote(string(obj)), value, hex.Dump([]byte(obj)), toGoByteSlice([]byte(obj)))
	default:
		return fmt.Sprintf("%s (%T)", fmt.Sprint(value), value)
	}
}

type tStringer interface {
	String() string
}

func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func findStringFirstDiffPos(s1 string, s2 string) int {
	for i := 0; i < len(s1); i++ {
		if i >= len(s2) {
			return i
		}
		if s1[i] != s2[i] {
			return i
		}
	}
	return -1
}

func EqualStringListNoOrder(a []string, b []string, objList ...interface{}) {
	sort.Strings(a)
	sort.Strings(b)
	Equal(a, b, objList...)
}

func toGoByteSlice(buf []byte) string {
	_buf := bytes.Buffer{}
	_buf.WriteString("[]byte{")
	for i, b := range buf {
		_buf.WriteString("0x")
		s := strconv.FormatInt(int64(b), 16)
		if len(s) == 1 {
			s = "0" + s
		}
		_buf.WriteString(s)
		_buf.WriteString(",")
		if i%16 == 15 {
			_buf.WriteString("\n")
		}
	}
	_buf.WriteString("}")
	return _buf.String()
}
