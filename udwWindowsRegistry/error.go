package udwWindowsRegistry

import "strings"

func IsErrorNotExist(err error) bool {
	return err != nil && strings.Contains(err.Error(), "The system cannot find the file specified.")
}
func IsErrorAccessDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Access is denied.")
}

func ErrorMsgIsNotFound(errMsg string) bool {
	return errMsg != "" && (strings.Contains(errMsg, "The system cannot find the file specified.") || strings.Contains(errMsg, "Sistem belirtilen yolu bulamıyor."))
}
