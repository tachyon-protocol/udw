package udwFile

import (
	"github.com/tachyon-protocol/udw/udwRand"
	"os"
	"path/filepath"
)

func NewTmpFilePath() string {
	return NewTmpFilePathWithExt("")
}

func NewTmpFilePathWithExt(ext string) string {
	dir := os.TempDir()
	file := "w8_" + udwRand.MustCryptoRandToReadableAlphaNum(12)
	if ext != "" {
		file += "." + ext
	}
	return filepath.Join(dir, file)
}

func MustChangeToTmpPath() string {
	folder := NewTmpFilePath()
	MustMkdir(folder)
	MustChangeDir(folder)
	return folder
}

func GetTempDir() string {
	return os.TempDir()
}
