package udwDebug

import (
	"encoding/json"
	"fmt"
	"github.com/tachyon-protocol/udw/udwReflect"
	"os"
	"reflect"
)

func Println(objList ...interface{}) {
	s := Sprintln(objList...)
	os.Stdout.WriteString(s)
	return
}

func Sprintln(objList ...interface{}) string {
	outList := make([]interface{}, len(objList)+1)
	outList[0] = "[udwDebug.Println]"
	for i := range objList {
		if udwReflect.IsNil(reflect.ValueOf(objList[i])) {
			outList[i+1] = "nil"
			continue
		}
		switch obj := objList[i].(type) {
		case []byte:
			outList[i+1] = fmt.Sprintf("%#v", obj)
		default:
			b, err := json.MarshalIndent(objList[i], "", " ")
			if err != nil {
				outList[i+1] = "[Println]error:" + err.Error()
				continue
			}
			outList[i+1] = string(b)
		}
	}
	return fmt.Sprintln(outList...)
}

func PrintlnSepLine() {
	fmt.Println(SepLineString)
}

const SepLineString = "========================================================="

const SepLineAndTwoLBString = "\n=========================================================\n"

const SepLineAndOneLBString = "=========================================================\n"
