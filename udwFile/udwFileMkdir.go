package udwFile

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"os"
	"path/filepath"
	"syscall"
)

func Mkdir(path string) (err error) {
	return os.MkdirAll(path, os.FileMode(0777))
}

func MustMkdir(dirname string) {
	err := os.MkdirAll(dirname, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustMkdir777(path string) {

	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return
		}
		panic(&os.PathError{"mkdir", path, syscall.ENOTDIR})
	}

	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) {
		i--
	}
	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) {
		j--
	}
	if j > 1 {

		MustMkdir777(path[0 : j-1])

	}

	err = os.Mkdir(path, os.FileMode(0777))
	if err != nil {

		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return
		}
		panic(err)
	}
	MustChmod(path, os.FileMode(0777))
	return
}

func MkdirForFile(path string) (err error) {
	path = filepath.Dir(path)
	return os.MkdirAll(path, os.FileMode(0777))
}

func MustMkdirForFile(path string) {
	path = filepath.Dir(path)
	err := os.MkdirAll(path, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustMkdirForFile777(path string) {
	path = filepath.Dir(path)
	MustMkdir777(path)
}

func MkdirForFile777(path string) (err error) {
	path = filepath.Dir(path)
	return udwErr.PanicToError(func() {
		MustMkdir777(path)
	})
}
