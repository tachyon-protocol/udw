package udwFile

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwRand"
	"os"
	"path/filepath"
)

type WriteFileAndMoveReq struct {
	TmpDir     string
	TargetPath string
	Content    []byte
}

func WriteFileAndMove(req WriteFileAndMoveReq) (err error) {

	fi, err := os.Stat(req.TargetPath)
	if err == nil && fi.Size() == int64(len(req.Content)) {
		content, err := ReadFile(req.TargetPath)
		if err == nil && bytes.Equal(req.Content, content) {
			return nil
		}
	}
	if req.TmpDir == "" {
		req.TmpDir = GetTempDir()
	}
	tmpFilePath := filepath.Join(req.TmpDir, "zzzig_w8_"+udwRand.MustCryptoRandToReadableAlphaNum(12))
	err = WriteFileWithMkdir(tmpFilePath, req.Content)
	if err != nil {
		Delete(tmpFilePath)
		return err
	}
	err = MkdirForFile(req.TargetPath)
	if err != nil {
		Delete(tmpFilePath)
		return err
	}
	err = os.Rename(tmpFilePath, req.TargetPath)
	if err != nil {
		Delete(tmpFilePath)
		return err
	}
	return nil
}

func MustWriteFileAndMove(req WriteFileAndMoveReq) {
	err := WriteFileAndMove(req)
	if err != nil {
		panic(err)
	}
}
