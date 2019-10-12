package udwCache

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestFileMd5ChangeCacheOneDir(t *testing.T) {
	udwFile.MustDelete("testFile")
	defer udwFile.MustDelete("testFile")
	callLog := make([]string, 32)

	udwFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	udwFile.MustDelete("testFile/d1")

	udwFile.MustMkdirAll("testFile/d1/d2")
	udwFile.MustWriteFile("testFile/d1/d2/f3", []byte("1"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[3] = "f3"
	})
	udwTest.Equal(callLog[3], "f3")

	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[4] = "f3"
	})
	udwTest.Equal(callLog[4], "")

	udwFile.MustWriteFile("testFile/d1/d2/f3", []byte("2"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[5] = "f3"
	})
	udwTest.Equal(callLog[5], "f3")

	udwFile.MustDelete("testFile/d1/d2/f3")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[6] = "f4"
	})
	udwTest.Equal(callLog[6], "f4")

	udwFile.MustWriteFile("testFile/d1/d2/f4", []byte("3"))
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[7] = "f4"
	})
	udwTest.Equal(callLog[7], "f4")

	udwFile.MustReadFile("testFile/d1/d2/f4")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[8] = "f4"
	})
	udwTest.Equal(callLog[8], "")

	udwFile.MustMkdir("testFile/d1/d2/f5")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile/d1",
	}, func() {
		callLog[9] = "f4"
	})
	udwTest.Equal(callLog[9], "")
}

func TestFileMd5ChangeCacheSymlink(t *testing.T) {
	udwFile.MustDelete("testFile")
	defer udwFile.MustDelete("testFile")
	callLog := make([]string, 32)

	udwFile.MustDeleteFile(getFileChangeCachePath("test_file_change_cache"))
	udwFile.MustDelete("testFile")
	udwFile.MustWriteFileWithMkdir("testFile/d1/d2", []byte("1"))
	udwFile.MustSymlink("d1", "testFile/d3")

	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile",
	}, func() {
		callLog[0] = "f3"
	})
	udwTest.Equal(callLog[0], "f3")
	MustMd5FileChangeCache("test_file_change_cache", []string{
		"testFile",
	}, func() {
		callLog[1] = "f3"
	})
	udwTest.Equal(callLog[1], "")
}
