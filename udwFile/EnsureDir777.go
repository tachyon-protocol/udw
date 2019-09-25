package udwFile

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"os"
)

func MustEnsureDir777(path string) {
	path = MustGetFullPath(path)
	stat, err := os.Lstat(path)
	if err == nil {
		if stat.IsDir() {
			mustEnsure777L1(path, stat)
			return
		}
		panic(fmt.Sprintf("[MustEnsureDir777] path:[%s] is not a dir", path))
	}
	if err != nil {
		if ErrorIsFileNotFound(err) {
			err := os.MkdirAll(path, 0777)
			if err != nil {
				if os.IsPermission(err) {
					udwCmd.CmdSlice([]string{"sudo", "mkdir", "-p", path}).MustRun()
				} else {
					panic(err)
				}
			}
			stat, err := os.Lstat(path)
			if err != nil {
				panic(err)
			}
			if !stat.IsDir() {
				panic(fmt.Sprintf("[MustEnsureDir777] path:[%s] is not a dir", path))
			}
			mustEnsure777L1(path, stat)
			return
		}
		panic(fmt.Sprintf("[MustEnsureDir777] path:[%s] err[%s]", path, err.Error()))
	}
}

func mustEnsure777L1(path string, stat os.FileInfo) {
	if stat.Mode().Perm() == 0777 {
		return
	}
	err := os.Chmod(path, 0777)
	if err == nil {
		return
	}
	if os.IsPermission(err) {
		udwCmd.CmdSlice([]string{"sudo", "chmod", "-R", "777", path}).MustRun()
	} else {
		panic(err)
	}
	return
}
