package udwErr

import "fmt"

func ErrorSprint(objList ...interface{}) error {
	s := fmt.Sprintln(objList...)
	return stringWrapError(s[:len(s)-1])
}
