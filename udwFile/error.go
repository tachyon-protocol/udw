package udwFile

import (
	"os"
	"strings"
)

func ErrorIsFileNotFound(err error) bool {
	return err != nil && (os.IsNotExist(err) || strings.Contains(err.Error(), "not a directory"))
}
func ErrorIsDirectory(err error) bool {
	return err != nil && strings.Contains(err.Error(), "is a directory")
}
