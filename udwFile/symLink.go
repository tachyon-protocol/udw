package udwFile

import (
	"fmt"
	"os"
	"path/filepath"
)

func MustSymlink(realPath string, linkPath string) {
	MustMkdirForFile(linkPath)
	err := os.Symlink(realPath, linkPath)
	if err == nil {
		return
	}
	if os.IsExist(err) {
		MustDelete(linkPath)
		err = os.Symlink(realPath, linkPath)
		if err == nil {
			return
		}
		panic(err)
	}
	panic(err)

}

func MustSymlinkRel(realPath string, linkPath string) {
	realPath = MustGetFullPath(realPath)
	linkPath = MustGetFullPath(linkPath)
	s, err := filepath.Rel(filepath.Dir(linkPath), realPath)
	if err != nil {
		panic(err)
	}
	MustSymlink(s, linkPath)
}

func MustSymlinkFullPath(realPath string, linkPath string) {
	realPath = MustGetFullPath(realPath)
	linkPath = MustGetFullPath(linkPath)
	MustSymlink(realPath, linkPath)
}

func FileInfoIsSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink > 0
}
func FileModeIsSymlink(fm os.FileMode) bool {
	return fm&os.ModeSymlink > 0
}

func IsSymlink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return FileInfoIsSymlink(fi)
}

func MustReadSymlink(path string) string {
	link, err := os.Readlink(path)
	if err != nil {
		panic(err)
	}
	return link
}

func MustGetAllSymlinkPathList(path string) (out []string) {
	path = MustGetFullPath(path)
	if !MustFileExist(path) {
		panic(fmt.Errorf("[GetAllSymlinkPathList] file not exist path:[%s]", path))
	}
	for {
		out = append(out, path)
		if !IsSymlink(path) {
			return out
		}
		relPath := MustReadSymlink(path)
		path = FullPathOnPath(filepath.Dir(path), relPath)
	}
}
