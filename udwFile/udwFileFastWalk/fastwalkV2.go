package udwFileFastWalk

import (
	"os"
	"sync"
)

type WalkCtx struct {
	Path string
	Typ  os.FileMode
}

func WalkV2(root string, walkFn func(walkCtx WalkCtx)) error {
	st, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if st.IsDir() == false {
		walkFn(WalkCtx{
			Path: root,
			Typ:  st.Mode(),
		})
		return nil
	}
	locker := sync.Mutex{}
	return Walk(root, func(path string, typ os.FileMode) error {
		locker.Lock()
		walkFn(WalkCtx{
			Path: path,
			Typ:  typ,
		})
		locker.Unlock()
		return nil
	})
}

func GetAllFileOneLevel(path string) (fileList []string, errMsg string) {
	err := readDir(path, func(dirName, baseName string, typ os.FileMode) error {
		thisPath := dirName + string(os.PathSeparator) + baseName
		if typ.IsRegular() {
			fileList = append(fileList, thisPath)
		}
		return nil
	})
	if err != nil {
		return nil, err.Error()
	}
	return fileList, ""
}
