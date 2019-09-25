package udwFile

import (
	"path/filepath"
	"strings"
)

func HasExt(path string, ext string) bool {
	return GetExt(path) == strings.ToLower(ext)
}

func GetFileBaseWithoutExt(p string) string {
	return filepath.Base(p[:len(p)-len(filepath.Ext(p))])
}

func GetExt(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

func GetExtWithoutDot(path string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
}
