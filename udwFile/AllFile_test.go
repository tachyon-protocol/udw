package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
	"time"
)

func TestGetAllFileAndDirectoryStat(t *testing.T) {
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	MustMkdirAll("zzzig_testFile/d1/d2")
	MustWriteFile("zzzig_testFile/d1/d2/f3", []byte("1"))
	out, err := GetAllFileAndDirectoryStat("zzzig_testFile")
	udwTest.Equal(err, nil)
	udwTest.Equal(len(out), 4)

	out, err = GetAllFileAndDirectoryStat("zzzig_testFile/d1/d2/f3")
	udwTest.Equal(err, nil)
	udwTest.Equal(len(out), 1)
}

func TestGetFileModeTime(t *testing.T) {
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	startTime := time.Now()
	MustWriteFileWithMkdir("zzzig_testFile/d1/d2/f3", []byte("1"))
	endTime := time.Now()
	modeTime, err := GetFileModifyTime("zzzig_testFile/d1/d2/f3")
	udwTest.Equal(err, nil)

	udwTest.Ok(modeTime.Before(endTime), modeTime.String(), endTime.String())
	udwTest.Ok(modeTime.After(startTime), modeTime.String(), startTime.String())
}
func TestMustGetAllDir(t *testing.T) {
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	MustWriteFileWithMkdir("zzzig_testFile/d1/d2/f3", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/d11/d22/f33", []byte("1"))
	dirs := MustGetAllDir("zzzig_testFile")
	fullPath := MustGetFullPath("zzzig_testFile")
	shouldExistPathList := []string{
		fullPath,
		filepath.Join(fullPath, "d11"),
		filepath.Join(fullPath, "d1"),
		filepath.Join(fullPath, "d1", "d2"),
		filepath.Join(fullPath, "d11", "d22"),
	}
	udwTest.EqualStringListNoOrder(shouldExistPathList, dirs)

	udwTest.AssertPanicWithErrorMessage(func() {
		MustGetAllDir("zzzig_testFileNotExist")
	}, "no such file or directory")
}

func TestMustGetAllFileOneLevel(t *testing.T) {
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	MustWriteFileWithMkdir("zzzig_testFile/d1/d2/f3", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/d11/d22/f33", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/f44", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/f55", []byte("1"))
	fileOneLevelList := MustGetAllFileOneLevel("zzzig_testFile")
	fullPath := MustGetFullPath("zzzig_testFile")
	shouldExistPathList := []string{
		filepath.Join(fullPath, "f44"),
		filepath.Join(fullPath, "f55"),
	}
	udwTest.EqualStringListNoOrder(shouldExistPathList, fileOneLevelList)

	udwTest.AssertPanicWithErrorMessage(func() {
		MustGetAllFileOneLevel("zzzig_testFileNotExist")
	}, "no such file or directory")
}

func TestMustGetAllFiles(t *testing.T) {
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	MustWriteFileWithMkdir("zzzig_testFile/d1/d2/f3", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/d11/d22/f33", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/f44", []byte("1"))
	MustWriteFileWithMkdir("zzzig_testFile/f55", []byte("1"))
	fileOneLevelList := MustGetAllFiles("zzzig_testFile")
	fullPath := MustGetFullPath("zzzig_testFile")
	shouldExistPathList := []string{
		filepath.Join(fullPath, "d1/d2/f3"),
		filepath.Join(fullPath, "d11/d22/f33"),
		filepath.Join(fullPath, "f44"),
		filepath.Join(fullPath, "f55"),
	}
	udwTest.EqualStringListNoOrder(shouldExistPathList, fileOneLevelList)
	fileOneLevelList = MustGetAllFiles("zzzig_testFile/f55")
	udwTest.EqualStringListNoOrder([]string{
		filepath.Join(fullPath, "f55"),
	}, fileOneLevelList)

	fileList := MustGetAllFiles("zzzig_testFileNotExist")
	udwTest.Equal(fileList, nil)

}
