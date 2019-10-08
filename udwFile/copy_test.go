package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"os"
	"testing"
	"time"
)

func TestMustCopy(t *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/d1/d2/f3", []byte("1"))
	MustSymlink("f3", "testFile/d1/d2/f4")
	MustCopy("testFile/d1", "testFile/t2")

	udwTest.Equal(MustReadFile("testFile/t2/d2/f3"), []byte("1"))
	udwTest.Equal(MustReadFile("testFile/t2/d2/f4"), []byte("1"))
	udwTest.Equal(IsSymlink("testFile/t2/d2/f3"), false)
	udwTest.Equal(IsSymlink("testFile/t2/d2/f4"), true)
	udwTest.Equal(MustReadSymlink("testFile/t2/d2/f4"), "f3")
}

func TestMustCheckContentCopyWithoutMerge(t *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/d1/d2/f3", []byte("1"))
	MustSymlink("f3", "testFile/d1/d2/f4")
	MustCheckContentCopyWithoutMerge("testFile/d1", "testFile/t2")

	udwTest.Equal(MustReadFile("testFile/t2/d2/f3"), []byte("1"))
	udwTest.Equal(MustReadFile("testFile/t2/d2/f4"), []byte("1"))
	udwTest.Equal(IsSymlink("testFile/t2/d2/f3"), false)
	udwTest.Equal(IsSymlink("testFile/t2/d2/f4"), true)
	udwTest.Equal(MustReadSymlink("testFile/t2/d2/f4"), "f3")
	fi, err := os.Lstat("testFile/t2/d2/f4")
	udwTest.Equal(err, nil)
	oldMtime := fi.ModTime()

	time.Sleep(time.Second)
	MustWriteFile("testFile/t2/d2/f3", []byte("2"))
	MustCheckContentCopyWithoutMerge("testFile/d1", "testFile/t2")
	udwTest.Equal(MustReadFile("testFile/t2/d2/f3"), []byte("1"))

	fi, err = os.Lstat("testFile/t2/d2/f4")
	udwTest.Equal(err, nil)
	udwTest.Equal(oldMtime, fi.ModTime())
}
