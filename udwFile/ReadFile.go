package udwFile

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
	"io/ioutil"
	"os"
)

func ReadFile(path string) (content []byte, err error) {
	return ioutil.ReadFile(path)
}

func MustReadFile(path string) (content []byte) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func MustReadFileOrIgnore(path string) (content []byte) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if ErrorIsFileNotFound(err) {
			return nil
		}
		panic(err)
	}
	return content
}

func MustReadFileAtStartPosWithLen(path string, startPos int64, thisLen int64) (content []byte) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if startPos < 0 {
		panic("[MustReadFileAtStartPos] startPos<0")
	}
	content = make([]byte, thisLen)
	n, err := f.ReadAt(content, startPos)
	if err != io.EOF && err != nil {
		panic(err)
	}
	return content[:n]
}

func TailByte(filePath string, size int64) (content []byte, err error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	start := int64(0)
	if info.Size() >= size {
		start = info.Size() - size
	} else {
		size = info.Size()
	}
	content = make([]byte, size)
	_, err = f.ReadAt(content, start)
	if err != io.EOF && err != nil {
		return nil, err
	}
	return content, nil
}

func ReadFileToBufW(path string, bufW *udwBytes.BufWriter) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var n int64
	fi, err := f.Stat()
	if err == nil {

		n = fi.Size()
	}
	n = n + 512
	if n > 32*1024 {
		n = 32 * 1024
	}
	bufW.TryGrow(int(n))
	_, err = bufW.ReadFrom(f)
	if err != nil {
		return err
	}
	return nil
}

func MustReadFileToBufW(path string, bufW *udwBytes.BufWriter) {
	err := ReadFileToBufW(path, bufW)
	if err != nil {
		panic(err)
	}
}
