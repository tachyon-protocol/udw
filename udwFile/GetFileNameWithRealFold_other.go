package udwFile

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MustGetFileOrDirectoryNameWithRealFold(path string) string {
	path = MustGetFullPath(path)

	filePathPartList := strings.Split(path, "/")
	retFileFoldList := []string{}
	for i := range filePathPartList {
		ret := mustGetOneBasePathFoldPart(filepath.Join("/", strings.Join(filePathPartList[:i+1], "/")))
		retFileFoldList = append(retFileFoldList, ret)
	}
	return filepath.Join(retFileFoldList...)
}

func mustGetOneBasePathFoldPart(path string) string {
	if path == "/" {
		return "/"
	}
	fileName := filepath.Base(path)
	dirPath := filepath.Dir(path)
	dirFile, err := os.Open(dirPath)
	if err != nil {
		panic(err)
	}
	defer dirFile.Close()
	names, err := dirFile.Readdirnames(-1)
	for _, n := range names {
		if strings.EqualFold(n, fileName) {
			return n
		}
	}
	panic(fmt.Errorf("[mustGetOneBasePathFoldPart] file not exist path:[%s]", path))
}

func MustMoveNameFoldCorrect(fromPath string, toPath string) {
	MustMove(fromPath, toPath)
	if MustIsFileOrDirectoryNameFoldCorrect(toPath) == true {
		return
	}
	toPath = MustGetFullPath(toPath)
	filePathPartList := strings.Split(toPath, "/")
	for i := 1; i < len(filePathPartList); i++ {
		thisExpectPath := filepath.Join("/", strings.Join(filePathPartList[:i+1], "/"))
		ret := mustGetOneBasePathFoldPart(thisExpectPath)

		if ret != filePathPartList[i] {
			fromPath := filepath.Join(filepath.Dir(thisExpectPath), ret)

			if i < len(filePathPartList)-1 || MustIsFile(fromPath) == false {

				middlePath := thisExpectPath + "." + mustGetRandString()
				MustMove(fromPath, middlePath)
				MustMove(middlePath, thisExpectPath)
			} else {

				MustMove(fromPath, thisExpectPath)
			}
		}
	}
}

func mustGetRandString() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
