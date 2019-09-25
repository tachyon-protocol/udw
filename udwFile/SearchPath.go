package udwFile

import (
	"errors"
	"os"
	"path/filepath"
)

func SearchFileInParentDir(startDirPath string, fileName string) (file string, err error) {
	startDirPath, err = filepath.Abs(startDirPath)
	if err != nil {
		return
	}
	p := startDirPath
	var udwFilePath string
	for {
		udwFilePath = filepath.Join(p, fileName)
		_, err = os.Stat(udwFilePath)
		if err == nil {

			return p, nil
		}
		if ErrorIsFileNotFound(err) == false {
			return "", err
		}
		thisP := filepath.Dir(p)
		if p == thisP {

			return "", errors.New("v2am52e9vb [fileName] " + fileName + " not found at [startDirPath]" + startDirPath)
		}
		p = thisP
	}
}
