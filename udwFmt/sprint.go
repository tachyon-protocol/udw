package udwFmt

import "fmt"

func Sprint(a ...interface{}) string {
	s := fmt.Sprintln(a...)
	return s[:len(s)-1]
}

func Sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}
