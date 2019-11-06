package udwConsole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwReflect"
	"github.com/tachyon-protocol/udw/udwStrings"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func MustRunCommandLineFromFuncV2(f interface{}) {
	errMsg := RunCommandLineFromFuncV3(RunCommandLineFromFuncV3Request{
		OsArgList: os.Args,
		F:         f,
	})
	if errMsg != "" {
		fmt.Println(errMsg)
		os.Exit(1)
		return
	}
	os.Exit(0)
}

type RunCommandLineFromFuncV3Request struct {
	OsArgList []string
	F         interface{}
}

func RunCommandLineFromFuncV3(req RunCommandLineFromFuncV3Request) (errMsg string) {
	gRunCommandLineFromFuncV3IsRunningLocker.Lock()
	if gRunCommandLineFromFuncV3IsRunning == true {
		gRunCommandLineFromFuncV3IsRunningLocker.Unlock()

		return "[RunCommandLineFromFuncV3] 96ny55peyd another insance is running in this process."
	}
	gRunCommandLineFromFuncV3IsRunning = true
	gRunCommandLineFromFuncV3IsRunningLocker.Unlock()
	defer func() {
		gRunCommandLineFromFuncV3IsRunningLocker.Lock()
		gRunCommandLineFromFuncV3IsRunning = false
		gRunCommandLineFromFuncV3IsRunningLocker.Unlock()
	}()

	simpleFunc, ok := req.F.(func())
	if ok {
		simpleFunc()
		return ""
	}
	mustEnsureValidFuncV2(req.F)
	v := reflect.ValueOf(req.F)
	t := v.Type()
	reqType := t.In(0)
	originReqTypeKind := reqType.Kind()
	var reqValue reflect.Value
	var originReqValue reflect.Value
	switch originReqTypeKind {
	case reflect.Struct:
		reqValue = reflect.New(reqType).Elem()
		originReqValue = reqValue
	case reflect.Ptr:
		reqType = reqType.Elem()
		originReqValue = reflect.New(reqType)
		reqValue = originReqValue.Elem()
	default:
		return "[MustRunCommandLineFromFunc] xsy3hgsrbg not support reqType Kind " + reqType.Kind().String()
	}

	var fieldList []*flagJson

	for _, field := range udwReflect.StructGetAllField(reqType) {
		argName := field.Name
		defaultValue := ""
		desc := ""
		value := field.Tag.Get("CmdFlag")
		if value == "-" {

			continue
		}
		valuePart := strings.Split(value, ",")
		if len(valuePart) >= 1 && valuePart[0] != "" {
			argName = valuePart[0]
		}
		if len(valuePart) >= 2 && valuePart[1] != "" {
			defaultValue = valuePart[1]
		}
		if len(valuePart) >= 3 && valuePart[2] != "" {
			desc = strings.Join(valuePart[2:], ",")
		}
		if argName == "h" || argName == "jsonHex" {
			return "[MustRunCommandLineFromFunc] txtuqsqs4a you can not define -h or -jsonHex flag"
		}
		if !udwStrings.IsAllAlphabetNum(argName) {
			return "[MustRunCommandLineFromFunc] djuyd4erzf [" + argName + "] has not support char"
		}
		thisFlagJson := &flagJson{}
		thisFlagJson.reflectValue = reqValue.FieldByIndex(field.Index)
		thisFlagJson.defaultValue = defaultValue
		thisFlagJson.desc = desc
		thisFlagJson.argName = argName
		switch field.Type.Kind() {
		case reflect.String:
			thisFlagJson.typS = "string"
		case reflect.Bool:
			thisFlagJson.typS = "bool"
			if thisFlagJson.defaultValue == "" {
				thisFlagJson.defaultValue = "false"
			}
		default:
			thisFlagJson.typS = thisFlagJson.reflectValue.Type().String()
			thisFlagJson.isJson = true
		}
		fieldList = append(fieldList, thisFlagJson)
	}
	inFlagKeyValueMap := map[string]string{}
	isAllFlagKeyValue := true
	if len(req.OsArgList) >= 2 {
		pos := 1
		for {
			if pos >= len(req.OsArgList) {
				break
			}
			s := req.OsArgList[pos]
			if strings.HasPrefix(s, "-") {
				s1 := strings.TrimLeft(s, "-")
				if strings.Contains(s1, "=") {
					key := udwStrings.StringBeforeFirstSubString(s1, "=")
					value := udwStrings.StringAfterFirstSubString(s1, "=")
					inFlagKeyValueMap[key] = value
					pos++
				} else {
					if pos+1 >= len(req.OsArgList) {
						inFlagKeyValueMap[s1] = ""
						pos++
					} else {
						nextValue := req.OsArgList[pos+1]
						if strings.HasPrefix(nextValue, "-") {
							inFlagKeyValueMap[s1] = ""
							pos++
						} else {
							inFlagKeyValueMap[s1] = nextValue
							pos += 2
						}
					}
				}
			} else {
				isAllFlagKeyValue = false
				pos++
			}
		}
	}
	cmdName := "unknownCommandName"
	if len(req.OsArgList) >= 1 {
		cmdName = req.OsArgList[0]
	}
	gUsageFn = func() string {
		_buf := bytes.Buffer{}
		_buf.WriteString("Usage of " + cmdName + ":\n")
		for _, field := range fieldList {
			_buf.WriteString("  -" + field.argName + " " + field.typS + "\n")
			if field.desc != "" {
				_buf.WriteString("        " + field.desc + "\n")
			}
		}
		_buf.WriteString(udwDebug.SepLineString + "\n")
		_buf.WriteString("  -jsonHex string\n")
		_buf.WriteString("        input argument jsonHex (this will override other flags.)\n")
		_buf.WriteString("  -h\n")
		_buf.WriteString("        help info\n")
		return _buf.String()
	}
	if len(inFlagKeyValueMap) == 1 {
		jsonHexS, ok := inFlagKeyValueMap["jsonHex"]
		if ok {

			jsonS := udwHex.MustDecodeStringToString(jsonHexS)
			err := json.Unmarshal([]byte(jsonS), reqValue.Addr().Interface())
			if err != nil {
				return "[MustRunCommandLineFromFunc] ffbqpny8fb " + err.Error()
			}
			v.Call([]reflect.Value{originReqValue})
			return ""
		}
		hS, ok := inFlagKeyValueMap["h"]
		if ok {

			if hS != "" {
				return "[MustRunCommandLineFromFunc] d5d5nv29ws -h should not have value"
			}
			return gUsageFn()
		}

	}
	if len(inFlagKeyValueMap) > 0 && isAllFlagKeyValue == false {

		return "[MustRunCommandLineFromFunc] yssspt9tyk mix -flag mode and step argument mode \n" + gUsageFn()
	}
	if len(inFlagKeyValueMap) == 0 && len(req.OsArgList) >= 2 {

		stepArgList := req.OsArgList[1:]
		if len(stepArgList) >= 3 {
			return "[MustRunCommandLineFromFunc] vrvxmtfff6 len(stepArgList) >=3 " + strconv.Itoa(len(stepArgList)) + "\n" + gUsageFn()
		}
		if len(fieldList) < len(stepArgList) {
			return "[MustRunCommandLineFromFunc] tw3sv54znj len(fieldList)<len(stepArgList) " + strconv.Itoa(len(fieldList)) + " " + strconv.Itoa(len(stepArgList)) + "\n" + gUsageFn()
		}
		for i := 0; i < len(stepArgList); i++ {
			thisField := fieldList[i]
			inFlagKeyValueMap[thisField.argName] = stepArgList[i]
		}
	}

	for _, thisFlagJson := range fieldList {
		value, ok := inFlagKeyValueMap[thisFlagJson.argName]
		if ok == false {

			value = thisFlagJson.defaultValue
		} else {
			if thisFlagJson.typS == "bool" {
				value = "true"
			}
		}
		switch thisFlagJson.typS {
		case "string":
			thisFlagJson.reflectValue.SetString(value)
		case "bool":
			b, err := strconv.ParseBool(value)
			if err != nil {
				return "[MustRunCommandLineFromFuncV2] bool parse fail fail1 [" + thisFlagJson.argName + "] [" + err.Error() + "] [" + value + "]\n" + gUsageFn()
			}
			thisFlagJson.reflectValue.SetBool(b)
		default:
			if value != "" {
				err := json.Unmarshal([]byte(value), thisFlagJson.reflectValue.Addr().Interface())
				if err != nil {
					return "[MustRunCommandLineFromFuncV2] json.Unmarshal fail1 [" + thisFlagJson.argName + "] [" + err.Error() + "] [" + value + "]\n" + gUsageFn()
				}
			}
		}
	}
	v.Call([]reflect.Value{originReqValue})
	return
}

type flagJson struct {
	value        string
	argName      string
	defaultValue string
	typS         string
	desc         string
	reflectValue reflect.Value
	isJson       bool
}

func MustNewCommandLineFuncFromFuncV2(f interface{}) func() {
	simpleFunc, ok := f.(func())
	if ok {
		return simpleFunc
	}

	mustEnsureValidFuncV2(f)
	return func() {
		MustRunCommandLineFromFuncV2(f)
	}
}

func mustEnsureValidFuncV2(f interface{}) {
	if f == nil {
		panic("[mustEnsureValidFuncV2] f==nil")
	}
	_, ok := f.(func())
	if ok {
		return
	}
	v := reflect.ValueOf(f)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("[mustEnsureValidFuncV2] need to pass in function t.Kind()[%d]!=reflect.Func", t.Kind()))
	}
	if t.NumIn() != 1 {
		panic(fmt.Errorf("[mustEnsureValidFuncV2] function can only have one parameter t.NumIn()[%d]!=1", t.NumIn()))
	}
	if t.NumOut() != 0 {
		panic(fmt.Errorf("[mustEnsureValidFuncV2] function does not allow return parameters t.NumOut()[%d]!=0", t.NumOut()))
	}
	reqType := t.In(0)
	switch reqType.Kind() {
	case reflect.Struct:
		return
	case reflect.Ptr:
		if reqType.Elem().Kind() == reflect.Struct {
			return
		}
		panic(fmt.Errorf("[mustEnsureValidFuncV2] The first argument to the function is either a struct or a pointer to a struct not support reqType Kind [%s] in reflect.Ptr", reqType.Kind()))
	default:
		panic(fmt.Errorf("[mustEnsureValidFuncV2] The first argument to the function is either a struct or a pointer to a struct not support reqType Kind [%s]", reqType.Kind()))
	}

}

func MustGetCommandLineString(header string, obj interface{}) string {
	jsonB, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	jsonS := string(jsonB)
	jsonHexS := udwHex.EncodeStringToString(jsonS)
	return header + " -jsonHex=" + jsonHexS
}

var gUsageFn func() string

func PrintUsageAndExit() {
	fmt.Println(gUsageFn())
	os.Exit(1)
}

func GetUsageString() string {
	return gUsageFn()
}

var gRunCommandLineFromFuncV3IsRunning = false
var gRunCommandLineFromFuncV3IsRunningLocker sync.Mutex
