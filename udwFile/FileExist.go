package udwFile

import (
	"os"
)

func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func FileExistIgnoreError(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func MustFileExist(path string) bool {
	exist, err := FileExist(path)
	if err != nil {
		panic(err)
	}
	return exist
}

func MustOnlyFileExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return false
		}
		panic(err)
	}
	return !fi.IsDir()
}

func OnlyFileExist(path string) (exist bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return !fi.IsDir(), nil
}

func MustDirectoryExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return false
		}
		panic(err)
	}
	return fi.IsDir()
}
