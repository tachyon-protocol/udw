package udwCache

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"time"
)

type NewFileMd5GetterRequest struct {
	DbFilePath        string
	Md5FileSumWithBuf func(path string, innerBuf []byte) (string, error)
}

func NewFileMd5GetterWithLock(dbFilePath string) FileMd5GetterWithLock {
	getterImp := getFileMd5GetterByDbFilePath(dbFilePath)

	getterImp.cacheNow = time.Now()
	getterImp.hasChange = false
	getterImp.dbFilePath = dbFilePath
	outGetter := FileMd5GetterWithLock{imp: getterImp}
	return outGetter
}

type FileMd5GetterWithLock struct {
	imp *fileMd5GetterImp
}

func (getter FileMd5GetterWithLock) GetMd5ByFullPath(path string) string {
	getter.imp.locker.Lock()
	out := getter.imp.getMd5ByFullPath__NOLOCK(path)
	getter.imp.locker.Unlock()
	return out
}
func (getter FileMd5GetterWithLock) GetMd5ByStatAndFullPath(statAndFullPath udwFile.StatAndFullPath) string {
	getter.imp.locker.Lock()
	out := getter.imp.getMd5ByStatAndFullPath__NOLOCK(statAndFullPath)
	getter.imp.locker.Unlock()
	return out
}
func (getter FileMd5GetterWithLock) Close() {
	getter.imp.locker.Lock()
	getter.imp.close__NOLOCK()
	getter.imp.locker.Unlock()
}
