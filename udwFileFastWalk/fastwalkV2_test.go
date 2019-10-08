package udwFileFastWalk_test

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwFileFastWalk"
	"github.com/tachyon-protocol/udw/udwTest"
	"path/filepath"
	"testing"
)

func TestWalkV2(t *testing.T) {
	thisPath := udwFile.NewTmpFilePath()
	defer udwFile.MustDelete(thisPath)
	udwFile.MustWriteFileWithMkdir(filepath.Join(thisPath, "1.txt"), []byte("1"))
	udwFile.MustWriteFileWithMkdir(filepath.Join(thisPath, "dir/1.txt"), []byte("1"))
	list := []udwFileFastWalk.WalkCtx{}
	thisCb := func(ctx udwFileFastWalk.WalkCtx) {
		list = append(list, ctx)
	}
	err := udwFileFastWalk.WalkV2(filepath.Join(thisPath, "1.txt"), thisCb)
	udwErr.PanicIfError(err)
	udwTest.Equal(len(list), 1)
	udwTest.Equal(list[0].Path, filepath.Join(thisPath, "1.txt"))

	list = list[:0]
	err = udwFileFastWalk.WalkV2(filepath.Join(thisPath, "dir"), thisCb)
	udwErr.PanicIfError(err)
	udwTest.Equal(len(list), 2)
	udwTest.Equal(list[0].Path, filepath.Join(thisPath, "dir"))
	udwTest.Equal(list[1].Path, filepath.Join(thisPath, "dir/1.txt"))
}
