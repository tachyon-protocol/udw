package udwImage

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func MustDecodeConfigFromFile(path string) (conf image.Config) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	conf, _, err = image.DecodeConfig(file)
	if err != nil {
		panic(err)
	}
	return conf
}

func MustDecodeImageFromFile(path string) (conf image.Image) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	conf, _, err = image.Decode(file)
	if err != nil {
		panic(err)
	}
	return conf
}

func MustDecodeImageFromByte(content []byte) image.Image {
	buf := bytes.NewBuffer(content)
	img, _, err := image.Decode(buf)
	if err != nil {
		panic(err)
	}
	return img
}

type ImageConfig struct {
	W      int
	H      int
	ImgExt string
}

func MustDecodeConfigFromByte(content []byte) (config ImageConfig) {
	buf := bytes.NewBuffer(content)
	imgConfig, imageType, err := image.DecodeConfig(buf)
	if err != nil {
		panic(err)
	}
	return ImageConfig{
		W:      imgConfig.Width,
		H:      imgConfig.Height,
		ImgExt: imageType,
	}
}

func DecodeConfigFromByte(content []byte) (config ImageConfig, err error) {
	buf := bytes.NewBuffer(content)
	imgConfig, imageType, err := image.DecodeConfig(buf)
	if err != nil {
		return config, err
	}
	return ImageConfig{
		W:      imgConfig.Width,
		H:      imgConfig.Height,
		ImgExt: imageType,
	}, nil
}

func GenerateOneColorImage(w int, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, c)
		}
	}
	return img
}
