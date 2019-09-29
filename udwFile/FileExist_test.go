package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestFileExist(ot *testing.T) {
	MustDelete("testFile")
	defer MustDelete("testFile")
	MustWriteFileWithMkdir("testFile/1.txt", []byte{})
	udwTest.Equal(MustFileExist("testFile/1.txt"), true)
	udwTest.Equal(MustFileExist("testFile/2.txt"), false)
	udwTest.Equal(MustFileExist("testFile/sub/2.txt"), false)
	udwTest.Equal(MustFileExist("testFile/1.txt/2.txt"), false)
}
