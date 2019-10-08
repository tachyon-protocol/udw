// +build appengine !linux,!darwin,!freebsd,!openbsd,!netbsd

package udwFileFastWalk

import (
	"io/ioutil"
	"os"
)

func readDir(dirName string, fn func(dirName, entName string, typ os.FileMode) error) error {
	fis, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	skipFiles := false
	for _, fi := range fis {
		if fi.Mode().IsRegular() && skipFiles {
			continue
		}
		if err := fn(dirName, fi.Name(), fi.Mode()&os.ModeType); err != nil {
			if err == SkipFiles {
				skipFiles = true
				continue
			}
			return err
		}
	}
	return nil
}
