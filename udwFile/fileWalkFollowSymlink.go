package udwFile

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func MustGetAllFileFollowSymlink(root string) (fileList []string) {
	root, err := GetFullPath(root)
	if err != nil {
		panic(err)
	}
	fi, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	mustGetAllFileFollowSymlinkCb(StatAndFullPath{
		Fi:       fi,
		FullPath: root,
	}, func(fs StatAndFullPath) {
		fileList = append(fileList, fs.FullPath)
	})
	return fileList
}

func mustGetAllFileFollowSymlinkCb(fs StatAndFullPath, cb func(fs StatAndFullPath)) {
	var err error
	isSymlink := FileInfoIsSymlink(fs.Fi)
	if isSymlink {
		fs.Fi, err = os.Stat(fs.FullPath)
		if err != nil {
			panic(err)
		}
	}
	if fs.Fi.IsDir() == false {
		cb(fs)
		return
	}
	fiList, err := ioutil.ReadDir(fs.FullPath)
	if err != nil {
		panic(err)
	}
	for _, fi := range fiList {
		thisRoot := filepath.Join(fs.FullPath, fi.Name())
		thisFs := StatAndFullPath{
			FullPath: thisRoot,
			Fi:       fi,
		}
		mustGetAllFileFollowSymlinkCb(thisFs, cb)
	}
	return
}

func MustGetAllDirFollowSymlink(root string) (dirList []string) {
	root, err := GetFullPath(root)
	if err != nil {
		panic(err)
	}
	fi, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	mustGetAllDirFollowSymlinkCb(StatAndFullPath{
		Fi:       fi,
		FullPath: root,
	}, func(fs StatAndFullPath) {
		dirList = append(dirList, fs.FullPath)
	})
	return dirList
}

func mustGetAllDirFollowSymlinkCb(fs StatAndFullPath, cb func(fs StatAndFullPath)) {
	var err error
	isSymlink := FileInfoIsSymlink(fs.Fi)
	if isSymlink {
		fs.Fi, err = os.Stat(fs.FullPath)
		if err != nil {
			panic(err)
		}
	}
	if fs.Fi.IsDir() == false {
		return
	}
	cb(fs)
	fiList, err := ioutil.ReadDir(fs.FullPath)
	if err != nil {
		panic(err.Error() + " " + fs.FullPath)
	}
	for _, fi := range fiList {
		thisRoot := filepath.Join(fs.FullPath, fi.Name())
		thisFs := StatAndFullPath{
			FullPath: thisRoot,
			Fi:       fi,
		}
		mustGetAllDirFollowSymlinkCb(thisFs, cb)
	}
	return
}
