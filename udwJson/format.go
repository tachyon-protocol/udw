package udwJson

import (
	"bytes"
	"encoding/json"
	"github.com/tachyon-protocol/udw/udwFile"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func MustMarshalIndentToString(obj interface{}) string {
	output, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(output)
}

func MustIntentJsonString(s string) string {
	var obj interface{}
	MustUnmarshalFromString(s, &obj)
	return MustMarshalIndentToString(obj)
}

var gHtmlUnescapeReplacerInit sync.Once
var gHtmlUnescapeReplacer *strings.Replacer

func MarshalIndent(obj interface{}) ([]byte, error) {
	gHtmlUnescapeReplacerInit.Do(func() {
		gHtmlUnescapeReplacer = strings.NewReplacer(`\u003c`, "<", `\u003e`, ">", `\u0026`, "&")
	})
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(gHtmlUnescapeReplacer.Replace(string(b))), nil
}

func MustWriteFile(path string, obj interface{}) {
	udwFile.MustMkdirForFile(path)
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, output, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustWriteFileIndent(path string, obj interface{}) {
	udwFile.MustMkdirForFile(path)
	output, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, output, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustFormatStringToString(s string) string {
	var obj interface{}
	MustUnmarshal([]byte(s), &obj)
	return MustMarshalIndentToString(obj)
}

func MustMarshalStringListIndentToString(obj []string) string {
	_buf := bytes.Buffer{}
	_buf.WriteString("[ ")
	for i, o := range obj {
		if i != 0 {
			_buf.WriteString("  ")
		}
		output, err := json.MarshalIndent(o, "", "  ")
		if err != nil {
			panic(err)
		}
		_buf.Write(output)
		if i != len(obj)-1 {
			_buf.WriteString(",\n")
		}
	}
	_buf.WriteString(" ]")
	return _buf.String()
}
