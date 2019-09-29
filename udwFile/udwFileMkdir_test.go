package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"os"
	"testing"
)

func TestMustMkdir(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/d1", []byte(""))
	udwTest.AssertPanic(func() {
		MustMkdir("testFile/d1")
	})
	MustMkdir("testFile/d2")
	udwTest.Equal(MustIsDir("testFile/d1"), false)
	udwTest.Equal(MustIsDir("testFile/d2"), true)
}

func TestMustMkdir777(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustMkdir777("testFile/dir1/dir2")
	udwTest.Equal(MustIsDir("testFile/dir1"), true)
	udwTest.Equal(MustIsDir("testFile/dir1/dir2"), true)
	udwTest.Equal(MustGetFilePerm("testFile/dir1"), os.FileMode(0777))
	udwTest.Equal(MustGetFilePerm("testFile/dir1/dir2"), os.FileMode(0777))

	MustWriteFile("testFile/abc.txt", []byte{})
	udwTest.AssertPanicWithErrorMessage(func() {
		MustMkdir777("testFile/abc.txt")
	}, "not a directory")

	MustMkdir("testFile/dir3")
	MustChmod("testFile/dir3", os.FileMode(0700))
	udwTest.Equal(MustGetFilePerm("testFile/dir3"), os.FileMode(0700))
	MustMkdir777("testFile/dir3")
	udwTest.Equal(MustGetFilePerm("testFile/dir3"), os.FileMode(0700))
}
