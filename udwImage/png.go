package udwImage

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"image"
	"image/png"
	"os"
)

func PngDecodeConfigFromFile(path string) (conf image.Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	return png.DecodeConfig(file)
}

func MustPngDecodeConfigFromFile(path string) (conf image.Config) {
	conf, err := PngDecodeConfigFromFile(path)
	if err != nil {
		panic(err)
	}
	return conf
}

func MustPngDecodeConfigFromBytes(content []byte) (conf image.Config) {
	buf := bytes.NewBuffer(content)
	conf, err := png.DecodeConfig(buf)
	if err != nil {
		panic(err)
	}
	return conf
}

func MustPngEncodeFromGoImageToBytes(img image.Image) []byte {
	buf := bytes.Buffer{}
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func MustPngEncodeFromGoImageToFile(img image.Image, inPath string) {
	content := MustPngEncodeFromGoImageToBytes(img)
	udwFile.MustWriteFileWithMkdir(inPath, content)
}

func MustPngDecodeFromBuf(content []byte) image.Image {
	buf := bytes.NewBuffer(content)
	img, err := png.Decode(buf)
	if err != nil {
		panic(err)
	}
	return img
}
