package udwCryptoMd5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

func Md5Hex(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5HexFromString(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5Byte(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func MustMd5File(path string) string {
	m, err := Md5File(path)
	if err != nil {
		panic(err)
	}
	return m
}

func MustMd5FileIgnoreNotExist(path string) string {
	m, err := Md5File(path)
	if err != nil {
		if err != nil && (os.IsNotExist(err) || strings.Contains(err.Error(), "not a directory")) {
			return ""
		}
		panic(err)
	}
	return m
}

func Md5File(path string) (string, error) {
	hash := md5.New()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func MustMd5FromIoReadCloser(r io.ReadCloser) string {
	defer r.Close()
	hash := md5.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5FromIoReadCloserWithBuf(r io.ReadCloser, innerBuf []byte) (string, error) {
	defer r.Close()
	hash := md5.New()
	_, err := io.CopyBuffer(hash, r, innerBuf)
	if err != nil {
		return "", err
	}
	hashB := hash.Sum(innerBuf[0:0])
	hex.Encode(innerBuf[16:], hashB)
	return string(innerBuf[16 : 16+32]), nil
}

func MustMd5FileWithBuf(path string, innerBuf []byte) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	out, err := Md5FromIoReadCloserWithBuf(f, innerBuf)
	if err != nil {
		panic(err)
	}
	return out
}
