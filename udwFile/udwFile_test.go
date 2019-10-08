package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestUdwFile(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	err := WriteFile("testFile", []byte(""))
	udwTest.Equal(err, nil)
	MustDeleteFile("testFile")
	MustDeleteFile("testFile")
	MustDeleteFileOrDirectory("testFile")
}

func TestMustDelete(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFile("testFile", []byte{})
	MustDelete("testFile/testFile")
	udwTest.Equal(MustFileExist("testFile"), true)
	MustDeleteFile("testFile")
	udwTest.Equal(MustFileExist("testFile"), false)
}

func TestHasExt(ot *testing.T) {
	udwTest.Equal(HasExt("/abc.go", ".go"), true)
	udwTest.Equal(HasExt("/abc.go", ".Go"), true)
}
