package udwFile

import (
	"io/ioutil"
)

func MustMkdirAll(dirname string) {
	MustMkdir(dirname)
}

func MustReadFileAll(path string) (content []byte) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func ReadFileAll(path string) (content []byte, err error) {
	return ioutil.ReadFile(path)
}

func MustDeleteFile(path string) {
	MustDelete(path)

}

func MustDeleteFileOrDirectory(path string) {
	MustDelete(path)

	return
}

func Realpath(inPath string) (string, error) {
	return GetFullPath(inPath)
}

func MustRealPath(inPath string) string {
	return MustGetFullPath(inPath)
}
