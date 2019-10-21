package udwCache_test

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCache"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestFileMd5Get(ot *testing.T) {
	dbPath := udwFile.MustGetFullPath("zzzig_testFile/FileMd5Get.db")
	udwFile.MustDelete("testFile")
	udwFile.MustWriteFileWithMkdir("zzzig_testFile/d1/d2", []byte("1"))

	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c4ca4238a0b923820dcc509a6f75849b")
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/dnotExist")), "")
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1")), "")
	})

	udwFile.MustWriteFileWithMkdir("zzzig_testFile/d1/d2", []byte("2"))
	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c81e728d9d4c2f636f067f89cc14862c")
	})
	udwFile.MustDelete("zzzig_testFile/d1/d2")
	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "")
	})

	udwFile.MustWriteFileWithMkdir("zzzig_testFile/d1/d2", []byte("2"))
	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c81e728d9d4c2f636f067f89cc14862c")
	})
	udwFile.MustWriteFileWithMkdir("zzzig_testFile/d1/d2", []byte("1"))
	time.Sleep(time.Second + 10*time.Millisecond)
	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c4ca4238a0b923820dcc509a6f75849b")
	})

	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c4ca4238a0b923820dcc509a6f75849b")
	})
	udwFile.MustWriteFileWithMkdir("zzzig_testFile/d1/d2", []byte("2"))
	udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
		udwTest.Equal(getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath("zzzig_testFile/d1/d2")), "c81e728d9d4c2f636f067f89cc14862c")
	})
}

func TestFileMd5GetSpeed(ot *testing.T) {

	dbPath := udwFile.MustGetFullPath("zzzig_testFile/FileMd5Get1.db")
	udwFile.MustDelete("testFile")

	for i := 0; i < 40; i++ {
		now := time.Now()
		fileList := udwFile.MustGetAllFiles("/usr/local/go/src")
		udwCache.FileMd5Get(dbPath, func(getter udwCache.FileMd5Getter) {
			for _, filePath := range fileList {
				getter.GetMd5ByFullPath__NOLOCK(udwFile.MustGetFullPath(filePath))
			}
		})
		fmt.Println(i, len(fileList), time.Since(now))
	}

}

func TestFileMd5CacheSpeed(ot *testing.T) {

	udwCache.MustMd5FileChangeCacheClean()

	for i := 0; i < 4; i++ {
		now := time.Now()
		udwCache.MustMd5FileChangeCache("test_file_change_cache", []string{"/usr/local/go/src"}, func() {})
		fmt.Println(i, time.Since(now))
	}

}
