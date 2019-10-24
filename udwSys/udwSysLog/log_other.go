// +build !ios,!android,!macAppStore

package udwSysLog

import "fmt"

func log(s string) {
	fmt.Println(s)
}
