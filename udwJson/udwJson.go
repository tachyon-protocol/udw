package udwJson

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/tachyon-protocol/udw/udwTypeTransform"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func MustReadFile(path string, obj interface{}) {
	err := ReadFile(path, obj)
	if err != nil {
		panic(err)
	}
}

func ReadFileTypeFix(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var iobj interface{}
	err = json.Unmarshal(b, &iobj)
	if err != nil {
		return err
	}
	return udwTypeTransform.Transform(iobj, obj)
}

func WriteFile(path string, obj interface{}) (err error) {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

func UnmarshalNoType(r []byte) (interface{}, error) {
	var obj interface{}
	err := json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func MustUnmarshal(r []byte, obj interface{}) {
	err := json.Unmarshal(r, obj)
	if err != nil {
		panic(err)
	}
	return
}

func MustUnmarshalFromString(r string, obj interface{}) {
	err := json.Unmarshal([]byte(r), obj)
	if err != nil {
		panic(err)
	}
	return
}

func UnmarshalFromString(r string, obj interface{}) error {
	return json.Unmarshal([]byte(r), obj)
}

func MustUnmarshalIgnoreEmptyString(jsonStr string, obj interface{}) {
	if jsonStr == "" {
		return
	}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		panic(err)
	}
	return
}

func MustUnmarshalToMap(r []byte) (obj map[string]interface{}) {
	err := json.Unmarshal(r, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func MustUnmarshalToMapDeleteBOM(r []byte) (obj map[string]interface{}) {
	r = DeleteBOM(r)
	err := json.Unmarshal(r, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func UnmarshalToMapDeleteBOM(r []byte) (obj map[string]interface{}, err error) {
	r = DeleteBOM(r)
	err = json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func MustMarshal(obj interface{}) []byte {
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return output
}

func MustMarshalToString(obj interface{}) string {
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(output)
}

func MarshalToString(obj interface{}) (string, error) {
	output, err := json.Marshal(obj)
	return string(output), err
}

func MustMarshalToStringForEqual(obj interface{}) string {
	output := MustMarshalToString(obj)
	return strings.Replace(output, "[]", "null", -1)
}

func DeleteBOM(fileBytes []byte) []byte {
	trimmedBytes := bytes.Trim(fileBytes, "\xef\xbb\xbf")
	return trimmedBytes
}

func JsLiteralObjectToJson(jsLiteralObj []byte) (json []byte, err error) {
	var (
		result  []byte
		isMatch bool
	)
	for i := 0; i < len(jsLiteralObj); i++ {
		if jsLiteralObj[i] == '\\' {
			if i+1 < len(jsLiteralObj) {
				switch jsLiteralObj[i+1] {
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u':
					if isMatch {
						result = append(result, jsLiteralObj[i])
						result = append(result, jsLiteralObj[i+1])
					}
				default:
					if !isMatch {
						isMatch = true
						result = make([]byte, 0, len(jsLiteralObj))

						for _, b := range jsLiteralObj[:i] {
							result = append(result, b)
						}
					}
					result = append(result, jsLiteralObj[i+1])
				}
				i++
				continue
			} else {
				return nil, errors.New("[uqb6j4f8q9] warning not a valid JavaScript literal object")
			}
		}
		if isMatch {
			result = append(result, jsLiteralObj[i])
		}
	}
	if isMatch {
		return result, nil
	}
	return jsLiteralObj, nil
}
