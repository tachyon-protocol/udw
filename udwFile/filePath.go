package udwFile

import (
	"path/filepath"
	"strings"
)

func PathBaseWithoutExt(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(filepath.Base(path), ext)
}

func PathTrimExt(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

type PathAndContentPair struct {
	Path    string
	Content []byte
}

func IsExtInList(fileName string, extList []string) (isExtInList bool) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	for _, ext := range extList {
		if strings.ToLower(ext) == fileExt {
			return true
		}
	}
	return false
}
