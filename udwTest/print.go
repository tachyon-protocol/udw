package udwTest

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func sprintln(objList ...interface{}) string {
	outList := make([]interface{}, len(objList)+1)
	outList[0] = "[udwTest.sprintln]"
	for i := range objList {
		if isNil(reflect.ValueOf(objList[i])) {
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
