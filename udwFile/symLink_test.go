package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"os"
	"path/filepath"
	"testing"
)

func TestMustSymlink(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/d1/d2/f3", []byte("1"))

	MustSymlink("f3", "testFile/d1/d2/f4")
	udwTest.Equal(MustReadFile("testFile/d1/d2/f4"), []byte("1"))

	fi, err := os.Lstat("testFile/d1/d2/f4")
	udwTest.Equal(FileInfoIsSymlink(fi), true)

	link, err := os.Readlink("testFile/d1/d2/f4")
	udwTest.Equal(err, nil)
	udwTest.Equal(link, "f3")

	MustMove("testFile/d1/d2/f4", "testFile/d1/d2/f5")
	udwTest.Equal(MustReadFile("testFile/d1/d2/f5"), []byte("1"))

	MkdirForFile("testFile/d1/d3/f5")
	MustMove("testFile/d1/d2/f5", "testFile/d1/d3/f5")
	udwTest.AssertPanicWithErrorMessage(func() {
		MustReadFile("testFile/d1/d3/f5")
	}, "no such file or directory")

	MustMkdir("testFile/d1/d2/f7")
	MustSymlink("f3", "testFile/d1/d2/f7")
	udwTest.Equal(MustReadFile("testFile/d1/d2/f7"), []byte("1"))
	udwTest.Equal(MustReadSymlink("testFile/d1/d2/f7"), "f3")

	MustWriteFile("testFile/d1/d2/f8", []byte("0"))
	MustSymlink("f3", "testFile/d1/d2/f8")
	udwTest.Equal(MustReadFile("testFile/d1/d2/f8"), []byte("1"))
	udwTest.Equal(MustReadSymlink("testFile/d1/d2/f8"), "f3")

	MustSymlink("testFile/d1/d2/f3", "testFile/d1/d2/f4")
	udwTest.AssertPanicWithErrorMessage(func() {
		MustReadFile("testFile/d1/d2/f4")
	}, "no such file or directory")
}

func TestMustGetAllSymlinkPathList(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/d1/d2/f3", []byte("1"))

	pwd := MustGetWd()
	MustSymlink("f3", "testFile/d1/d2/f4")
	MustSymlink("f4", "testFile/d1/d2/f5")
	pathList := MustGetAllSymlinkPathList("testFile/d1/d2/f3")
	udwTest.Equal(pathList, []string{
		filepath.Join(pwd, "testFile/d1/d2/f3"),
	})

	pathList = MustGetAllSymlinkPathList("testFile/d1/d2/f4")
	udwTest.Equal(pathList, []string{
		filepath.Join(pwd, "testFile/d1/d2/f4"),
		filepath.Join(pwd, "testFile/d1/d2/f3"),
	})

	pathList = MustGetAllSymlinkPathList("testFile/d1/d2/f5")
	udwTest.Equal(pathList, []string{
		filepath.Join(pwd, "testFile/d1/d2/f5"),
		filepath.Join(pwd, "testFile/d1/d2/f4"),
		filepath.Join(pwd, "testFile/d1/d2/f3"),
	})

	MustSymlink("f7", "testFile/d1/d2/f6")
	udwTest.AssertPanic(func() {
		MustGetAllSymlinkPathList("testFile/d1/d2/f6")
	})
}
