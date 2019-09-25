package udwFile

import (
	"bytes"
//	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDotFile(path string) bool {
	if path == "./" {
		return false
	}
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}
	return false
}

func WriteFile(path string, content []byte) (err error) {
	return ioutil.WriteFile(path, content, os.FileMode(0777))
}
func MustWriteFile(path string, content []byte) {
	if path == "" {
		panic(`[MustWriteFile] path==""`)
	}
	err := ioutil.WriteFile(path, content, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func MustWriteFileWithMkdir(path string, content []byte) {
	MustMkdirForFile(path)
	MustWriteFile(path, content)
}

func WriteFileWithMkdir(path string, content []byte) (err error) {
	err = MkdirForFile(path)
	if err != nil {
		return err
	}
	err = WriteFile(path, content)
	if err != nil {
		return err
	}
	return nil
}

func MustCheckContentAndWriteFileWithMkdir(path string, content []byte) {
	if !MustOnlyFileExist(path) {
		MustWriteFileWithMkdir(path, content)
		return
	}
	oldContent := MustReadFile(path)
	if bytes.Equal(oldContent, content) {
		return
	}
	MustWriteFileWithMkdir(path, content)
	return
}

func AppendFile(path string, content []byte) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(content)
	return
}

func MustAppendFile(path string, content []byte) {
	err := AppendFile(path, content)
	if err != nil {
		panic(err)
	}
}

func MustAppendFileAddLineEnd(path string, content []byte) {
	content = append(content, byte('\n'))
	err := AppendFile(path, content)
	if err != nil {
		panic(err)
	}
}

func RemoveExtFromFilePath(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func Delete(path string) (err error) {
	err = os.RemoveAll(path)
	if ErrorIsFileNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func MustDelete(path string) {
	err := os.RemoveAll(path)
	if ErrorIsFileNotFound(err) {
		return
	}
	if err != nil {
		panic(err)
	}
	return
}

func MustChangeDir(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func MustChdir(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func MustRename(oldpath string, newpath string) {
	MustMkdirForFile(newpath)
	err := os.Rename(oldpath, newpath)
	if err == nil {
		return
	}
	for err != nil {
		if strings.HasSuffix(err.Error(), "different disk drive.") ||
			strings.HasSuffix(err.Error(), "cross-device link") {
			MustDeleteFileOrDirectory(newpath)
			MustCopy(oldpath, newpath)
			err = nil
		} else if strings.HasSuffix(err.Error(), "not a directory") ||
			strings.HasSuffix(err.Error(), "Access is denied.") ||
			strings.HasSuffix(err.Error(), "file exists") {
			MustDeleteFileOrDirectory(newpath)
			oldErr := err
			err = os.Rename(oldpath, newpath)
			if err == oldErr {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	MustDeleteFileOrDirectory(oldpath)
}

func MustGetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func MustSetWd(workdingDirectory string) {
	MustChangeDir(workdingDirectory)
}

func MustMove(srcPath string, dstPath string) {
	MustRename(srcPath, dstPath)
}

func MustGetFilePerm(path string) os.FileMode {
	fi, err := os.Lstat(path)
	if err != nil {
		panic(err)
	}
	var mode os.FileMode
	mode = fi.Mode()
	return mode.Perm()
}

func TruncateFileToAimSizeFromEnd(filePath string, length int64) (err error) {
	fileSize, err := GetFileSize(filePath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filePath, os.O_RDWR, os.FileMode(0777))
	if err != nil {
		return err
	}
	if fileSize <= length {
		length = fileSize
	}
	buf := make([]byte, length)
	_, err = f.ReadAt(buf, fileSize-length)
	if err != nil && err != io.EOF {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	n, err := f.Write(buf)
	if err == nil && n < len(buf) {
		return io.ErrShortWrite
	}
	f.Close()
	return nil
}
