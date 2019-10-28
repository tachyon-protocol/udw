package udwSqlite3

import (
	"strings"
)

func errorIsTableNotExist(errMsg string) bool {
	if errMsg == "" {
		return false
	}
	if strings.Contains(errMsg, "no such table:") {
		return true
	}
	return false
}

func IsErrorDatabaseCorrupt(errMsg string) bool {
	return strings.Contains(errMsg, "file is not a database") ||
		strings.Contains(errMsg, "file is encrypted or is not a database") ||
		strings.Contains(errMsg, "database disk image is malformed") ||
		strings.Contains(errMsg, "unsupported file format") ||
		strings.Contains(errMsg, "malformed database schema") ||
		(strings.Contains(errMsg, "has no column named"))
}
