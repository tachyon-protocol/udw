package udwFile

import (
	"os"
	"strings"
)

func MustChmod(filepath string, mode os.FileMode) {
	err := Chmod(filepath, mode)
	if err != nil {
		panic(err)
	}
}

func Chmod(filepath string, mode os.FileMode) (err error) {
	err = os.Chmod(filepath, mode)
	if err != nil {
		errS := err.Error()
		if strings.Contains(errS, "operation not permitted") {
			fi, err1 := os.Lstat(filepath)
			if err1 == nil && fi.Mode().Perm() == mode {

				return nil
			}
		}
		return err
	}
	return nil
}
