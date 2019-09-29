package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
)

func TestMustGetAllFilesFollowSymlink1(t *testing.T) {
	MustDelete("zzzig_test")
	defer MustDelete("zzzig_test")
	basePath := MustGetFullPath("zzzig_test")
	MustWriteFileWithMkdir("zzzig_test/d1/d2/f3", []byte("1"))
	MustWriteFileWithMkdir("zzzig_test/d11/d22/f33", []byte("1"))
	MustSymlink(MustGetFullPath("zzzig_test/d1"), "zzzig_test/d12")
	MustSymlink(MustGetFullPath("zzzig_test/d1/d2/f3"), "zzzig_test/d13")
	fileList := MustGetAllFileFollowSymlink("zzzig_test")
	udwTest.Equal(len(fileList), 4)
	udwTest.Equal(fileList[0], filepath.Join(basePath, "d1/d2/f3"))
	udwTest.Equal(fileList[1], filepath.Join(basePath, "d11/d22/f33"))
	udwTest.Equal(fileList[2], filepath.Join(basePath, "d12/d2/f3"))
	udwTest.Equal(fileList[3], filepath.Join(basePath, "d13"))
	fileList = MustGetAllFileFollowSymlink("zzzig_test/d12")
	udwTest.Equal(len(fileList), 1)
	udwTest.Equal(fileList[0], MustGetFullPath("zzzig_test/d12/d2/f3"))
	fileList = MustGetAllFileFollowSymlink("zzzig_test/d1/d2/f3")
	udwTest.Equal(len(fileList), 1)
	udwTest.Equal(fileList[0], MustGetFullPath("zzzig_test/d1/d2/f3"))
}

func TestMustGetAllFilesFollowSymlink2(t *testing.T) {
	MustDelete("zzzig_test")
	defer MustDelete("zzzig_test")
	basePath := MustGetFullPath("zzzig_test")
	MustWriteFileWithMkdir("zzzig_test/d1/d2/f3", []byte("1"))
	MustSymlink(MustGetFullPath("zzzig_test/d1"), "zzzig_test/s1")
	MustSymlink(MustGetFullPath("zzzig_test/s1"), "zzzig_test/s2")
	fileList := MustGetAllFileFollowSymlink("zzzig_test")
	udwTest.Equal(len(fileList), 3)
	udwTest.Equal(fileList[0], filepath.Join(basePath, "d1/d2/f3"))
	udwTest.Equal(fileList[1], filepath.Join(basePath, "s1/d2/f3"))
	udwTest.Equal(fileList[2], filepath.Join(basePath, "s2/d2/f3"))
}

func TestMustGetAllDirFollowSymlink(t *testing.T) {
	MustDelete("zzzig_test")
	defer MustDelete("zzzig_test")
	basePath := MustGetFullPath("zzzig_test")
	MustWriteFileWithMkdir("zzzig_test/d1/d2/f3", []byte("1"))
	MustWriteFileWithMkdir("zzzig_test/d11/d22/f33", []byte("1"))
	MustSymlink(MustGetFullPath("zzzig_test/d1"), "zzzig_test/d12")
	MustSymlink(MustGetFullPath("zzzig_test/d1/d2/f3"), "zzzig_test/d13")
	fileList := MustGetAllDirFollowSymlink("zzzig_test")
	udwTest.Equal(len(fileList), 7)
	udwTest.Equal(fileList, []string{
		filepath.Join(basePath),
		filepath.Join(basePath, "d1"),
		filepath.Join(basePath, "d1/d2"),
		filepath.Join(basePath, "d11"),
		filepath.Join(basePath, "d11/d22"),
		filepath.Join(basePath, "d12"),
		filepath.Join(basePath, "d12/d2"),
	})
}
