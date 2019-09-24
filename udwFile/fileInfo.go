package udwFile

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

func MustIsFile(filepath string) bool {
	fi, err := os.Lstat(filepath)
	if ErrorIsFileNotFound(err) {
		return false
	}
	if err != nil {
		panic(fmt.Errorf("[MustIsFile] os.Lstat %s", err))
	}
	return FileInfoIsFile(fi)
}

func FileInfoIsFile(fi os.FileInfo) bool {
	if FileInfoIsSymlink(fi) {
		return false
	}
	return fi.IsDir() == false
}

func FileInfoIsDir(fi os.FileInfo) bool {
	if FileInfoIsSymlink(fi) {
		return false
	}
	return fi.IsDir() == true
}

func MustIsDir(filepath string) bool {
	fi, err := os.Lstat(filepath)
	if ErrorIsFileNotFound(err) {
		return false
	}
	if err != nil {
		panic(fmt.Errorf("[MustIsDir] os.Lstat %s", err))
	}
	if FileInfoIsSymlink(fi) {
		return false
	}
	return fi.IsDir()
}
func SymlinkIsDirectory(path string) bool {
	if IsSymlink(path) {
		path, err := filepath.EvalSymlinks(path)
		if err == nil {
			fi, err := os.Lstat(path)
			if ErrorIsFileNotFound(err) {
				return false
			}
			return fi.IsDir()
		}
	}
	return false
}

func GetFileNameParseUrlv(in string) (out string) {
	appName, err := url.QueryUnescape(in)
	if err != nil {
		return ""
	}
	return appName
}
