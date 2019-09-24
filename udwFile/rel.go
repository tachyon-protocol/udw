package udwFile

import (
	"fmt"
	"path/filepath"
	"strings"
)

func MustGetRelativePath(shortPath string, longPath string) string {
	s, err := GetRelativePath(shortPath, longPath)
	if err != nil {
		panic(err)
	}
	return s
}

func GetRelativePath(shortPath string, longPath string) (string, error) {
	if len(shortPath) > len(longPath) {
		return "", fmt.Errorf("[MustGetRelativePath] len(shortPath)[%s]>len(longPath)[%s]", shortPath, longPath)
	}
	ret, err := filepath.Rel(shortPath, longPath)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func GetSamePrefixPath(shortPath string, longPath string) string {
	if shortPath == longPath {
		return shortPath
	}
	shortList := strings.Split(shortPath, "/")
	longList := strings.Split(longPath, "/")
	var i int
	for i = 0; ; i++ {
		if i >= len(shortList) || i >= len(longList) {
			break
		}
		if shortList[i] != longList[i] {
			break
		}
	}
	return strings.Join(shortList[:i], "/")
}
