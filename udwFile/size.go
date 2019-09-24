package udwFile

import (
	"os"
)

func MustGetDirectorySize(path string) int64 {
	statPathList := MustGetAllFileAndDirectoryStat(path)
	dirSize := int64(0)
	for _, stat := range statPathList {
		dirSize += stat.Fi.Size()
	}
	return dirSize
}

func MustGetFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

func GetFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func GetFileSizeOrZero(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fi.Size()
}
