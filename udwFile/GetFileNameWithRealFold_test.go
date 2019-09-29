package udwFile

import (
	"github.com/tachyon-protocol/udw/udwPlatform"
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
)

func TestMustGetFileOrDirectoryNameWithRealFold(ot *testing.T) {
	if udwPlatform.IsLinux() {
		MustWriteFileWithMkdir("zzzig_testFile/a.txt", []byte("a"))
		fullPath := MustGetFullPath(".")
		realFoldPath := MustGetFileOrDirectoryNameWithRealFold("testFile/a.txt")
		udwTest.Equal(realFoldPath, filepath.Join(fullPath, "testFile/a.txt"))
		udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("testFile/a.txt"), true)
		realFoldPath = MustGetFileOrDirectoryNameWithRealFold("testFile/A.txt")
		udwTest.Equal(realFoldPath, filepath.Join(fullPath, "testFile/A.txt"))
		udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("testFile/A.txt"), true)
		return
	}
	MustDelete("zzzig_testFile")
	defer MustDelete("zzzig_testFile")
	MustWriteFileWithMkdir("zzzig_testFile/a.txt", []byte("a"))
	fullPath := MustGetFullPath(".")
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/a.txt"), true)
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/A.txt"), false)
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFilE/a.txt"), false)

	realFoldPath := MustGetFileOrDirectoryNameWithRealFold("zzzig_testFile/a.txt")
	udwTest.Equal(realFoldPath, filepath.Join(fullPath, "zzzig_testFile/a.txt"))
	realFoldPath = MustGetFileOrDirectoryNameWithRealFold("zzzig_testFile/A.txt")
	udwTest.Equal(realFoldPath, filepath.Join(fullPath, "zzzig_testFile/a.txt"))
	realFoldPath = MustGetFileOrDirectoryNameWithRealFold("zzzig_testfile/A.txt")
	udwTest.Equal(realFoldPath, filepath.Join(fullPath, "zzzig_testFile/a.txt"))
	udwTest.AssertPanic(func() {
		MustGetFileOrDirectoryNameWithRealFold("zzzig_testFile/b.txt")
	})
	udwTest.AssertPanic(func() {
		MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/b.txt")
	})

	MustWriteFileWithMkdir("zzzig_testFile/a/a.txt", []byte("1"))
	MustMoveNameFoldCorrect("zzzig_testFile/a/a.txt", "zzzig_testFile/a/A.txt")
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/a/A.txt"), true)

	MustWriteFileWithMkdir("zzzig_testFile/b/b.txt", []byte("1"))
	udwTest.Equal(MustReadFile("zzzig_testFile/b/b.txt"), []byte("1"))
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/b/b.txt"), true)

	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/b"), true)
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/B"), false)
	MustMoveNameFoldCorrect("zzzig_testFile/b/b.txt", "zzzig_testFile/B/b.txt")
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/b/b.txt"), false)
	udwTest.Equal(MustIsFileOrDirectoryNameFoldCorrect("zzzig_testFile/B/b.txt"), true)
}
