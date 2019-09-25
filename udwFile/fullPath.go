package udwFile

import (
	"path/filepath"
	"strings"
)

func FullPathOnPath(workingPath string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(workingPath, path)
}

func GetFullPath(inPath string) (string, error) {
	return filepath.Abs(inPath)

}

func MustGetFullPath(inPath string) string {
	outPath, err := GetFullPath(inPath)
	if err != nil {
		panic(err)
	}
	return outPath
}

func IsDangerFullPath(thisPath string) bool {
	thisPath = filepath.Clean(thisPath)
	if strings.Contains(thisPath, "..") {
		return true
	}
	if thisPath == "/" || thisPath == "" {
		return true
	}
	if strings.HasSuffix(thisPath, `:\`) {
		return true
	}
	return false
}
