package udwDebug

import "fmt"

func FmtSprint(objList ...interface{}) string {
	s := fmt.Sprintln(objList...)
	return s[:len(s)-1]
}
