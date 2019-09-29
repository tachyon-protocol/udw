package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestSearchFileInParentDir(t *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/.testFile", []byte("1"))
	MustMkdir("testFile/a/b")
	p, err := SearchFileInParentDir("testFile/a/b", ".testFile")
	udwTest.Equal(err, nil)
	udwTest.Equal(p, MustGetFullPath("testFile"))

	p, err = SearchFileInParentDir("testFile/a/b/c/d", ".testFile")
	udwTest.Equal(err, nil)
	udwTest.Equal(p, MustGetFullPath("testFile"))

	p, err = SearchFileInParentDir("testFile/.testFile", ".testFile")
	udwTest.Equal(err, nil)
	udwTest.Equal(p, MustGetFullPath("testFile"))
}
