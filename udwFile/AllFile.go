package udwFile

import (
	"github.com/tachyon-protocol/udw/udwFileFastWalk"
	"github.com/tachyon-protocol/udw/udwSort"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type StatAndFullPath struct {
	Fi       os.FileInfo
	FullPath string
}

func GetAllFileAndDirectoryStat(root string) (out []StatAndFullPath, err error) {
	return GetAllFileAndDirectoryStatWithFilter(root, func(path string, info os.FileInfo) bool {
		return true
	})
}
func GetAllFileAndDirectoryStatWithFilter(root string, filter func(path string, info os.FileInfo) bool) (out []StatAndFullPath, err error) {
	root, err = GetFullPath(root)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if ErrorIsFileNotFound(err) {
				return nil
			}
			return err
		}
		if filter(path, info) {
			out = append(out, StatAndFullPath{
				FullPath: path,
				Fi:       info,
			})
		}
		return nil
	})
	return
}

func MustGetAllFileAndDirectoryStat(root string) (out []StatAndFullPath) {
	out, err := GetAllFileAndDirectoryStat(root)
	if err != nil {
		panic(err)
	}
	return out
}

const WalkV2RespNext = 1
const WalkV2RespSkipDir = 2
const WalkV2RespBreak = 3

func MustWalkV2(root string, cb func(st StatAndFullPath) int) {
	root, err := GetFullPath(root)
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if ErrorIsFileNotFound(err) {
				return nil
			}
			return err
		}
		ret := cb(StatAndFullPath{
			FullPath: path,
			Fi:       info,
		})
		switch ret {
		case WalkV2RespSkipDir:
			return filepath.SkipDir
		case WalkV2RespBreak:
			return os.ErrNotExist
		case WalkV2RespNext:
			return nil
		default:
			panic("[hpnu4kn7xk] unknown status:[" + strconv.Itoa(ret) + "]")
		}
	})
	if err == os.ErrNotExist {
		return
	}
	if err != nil {
		panic(err)
	}
	return
}

func GetFileModifyTime(path string) (modeTime time.Time, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return modeTime, err
	}

	modeTime = stat.ModTime()
	return modeTime, nil
}

func MustGetFileModifyTime(path string) (modeTime time.Time) {
	modeTime, err := GetFileModifyTime(path)
	if err != nil {
		panic(err)
	}
	return modeTime
}

func GetAllFiles(root string) (out []string, err error) {
	root, err = GetFullPath(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(root)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if FileInfoIsFile(info) {
		return []string{root}, nil
	}
	locker := sync.Mutex{}
	err = udwFileFastWalk.Walk(root, func(path string, typ os.FileMode) error {
		if typ.IsRegular() {
			locker.Lock()
			out = append(out, path)
			locker.Unlock()
		}
		return nil
	})
	if ErrorIsFileNotFound(err) {
		return nil, nil
	}
	return out, err

}

func MustGetAllFiles(root string) (out []string) {
	out, err := GetAllFiles(root)
	if err != nil {
		panic(err)
	}
	return out
}

func MustGetAllFilesSort(root string) (out []string) {
	out = MustGetAllFiles(root)
	udwSort.SortString(out)
	return out
}

func MustGetAllFileAndDirAndSymlink(root string) (out []string) {
	root, err := GetFullPath(root)
	if err != nil {
		panic(err)
	}
	info, err := os.Lstat(root)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return nil
		}
		panic(err)
	}
	if info.IsDir() == false {
		return []string{root}
	}
	locker := sync.Mutex{}
	err = udwFileFastWalk.Walk(root, func(path string, typ os.FileMode) error {
		locker.Lock()
		out = append(out, path)
		locker.Unlock()
		return nil
	})
	if ErrorIsFileNotFound(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return out
}

func MustGetAllFileOneLevel(path string) (fileList []string) {
	path = MustGetFullPath(path)
	fileList, errMsg := udwFileFastWalk.GetAllFileOneLevel(path)
	if errMsg != "" {
		panic(errMsg)
	}
	return fileList

}

func MustGetAllStatAndFullPathOneLevel(path string) (fileAndStatList []StatAndFullPath) {
	fiList, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, fi := range fiList {
		fileAndStatList = append(fileAndStatList, StatAndFullPath{
			FullPath: filepath.Join(path, fi.Name()),
			Fi:       fi,
		})
	}
	return fileAndStatList
}

func MustGetAllFileAndStatOneLevel(path string) (fileAndStatList []StatAndFullPath) {
	fiList, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, fi := range fiList {
		if !fi.IsDir() {
			fileAndStatList = append(fileAndStatList, StatAndFullPath{
				FullPath: filepath.Join(path, fi.Name()),
				Fi:       fi,
			})
		}
	}
	return fileAndStatList
}

func MustGetAllDirOneLevel(path string) (fileList []string) {
	fileList, err := GetAllDirOneLevel(path)
	if err != nil {
		panic(err)
	}
	return fileList
}

func GetAllDirOneLevel(path string) (fileList []string, err error) {
	fiList, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range fiList {
		if fi.IsDir() {
			fileList = append(fileList, filepath.Join(path, fi.Name()))
		}
	}
	return fileList, nil
}

func MustGetAllFileFromPathList(pathlist []string) (outList []string) {
	for _, root := range pathlist {
		out, err := GetAllFiles(root)
		if err != nil {
			panic(err)
		}
		outList = append(outList, out...)
	}
	return outList
}

func GetAllDir(root string) ([]string, error) {
	var out []string
	root, err := GetFullPath(root)
	if err != nil {
		return nil, err
	}
	locker := sync.Mutex{}
	err = udwFileFastWalk.Walk(root, func(path string, typ os.FileMode) error {
		if typ.IsDir() {
			locker.Lock()
			out = append(out, path)
			locker.Unlock()
		}
		return nil
	})

	return out, err
}

func MustGetAllDir(root string) (out []string) {
	out, err := GetAllDir(root)
	if err != nil {
		panic(err)
	}
	return out
}

func MustGetAllDirSort(root string) (out []string) {
	out = MustGetAllDir(root)
	udwSort.SortString(out)
	return out
}

func ReadDirFileOneLevel(path string) (fileList []string, err error) {
	fiList, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range fiList {
		if !fi.IsDir() {
			fileList = append(fileList, fi.Name())
		}
	}
	return fileList, nil
}

func MustReadDirFileOneLevel(path string) (fileList []string) {
	fileList, err := ReadDirFileOneLevel(path)
	if err != nil {
		panic(err)
	}
	return fileList
}
