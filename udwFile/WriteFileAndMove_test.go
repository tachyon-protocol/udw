package udwFile

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestWriteFileAndMove(t *testing.T) {
	prefix := "zzzig_testFile"
	MustDelete(prefix)
	defer MustDelete(prefix)
	MustWriteFileAndMove(WriteFileAndMoveReq{
		TargetPath: prefix + "/1.txt",
		Content:    []byte("1"),
	})
	udwTest.Equal(MustReadFile(prefix+"/1.txt"), []byte("1"))
	MustWriteFileAndMove(WriteFileAndMoveReq{
		TargetPath: prefix + "/1.txt",
		Content:    []byte("1"),
	})
	udwTest.Equal(MustReadFile(prefix+"/1.txt"), []byte("1"))
	MustWriteFileAndMove(WriteFileAndMoveReq{
		TargetPath: prefix + "/1.txt",
		Content:    []byte("2"),
	})
	udwTest.Equal(MustReadFile(prefix+"/1.txt"), []byte("2"))
}
