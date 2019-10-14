package udwImageResize

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

type MustResizePngFileToHeightAndWidthRequest struct {
	InFilePath  string
	OutFilePath string
	Height      int
	Width       int
}

func MustResizePngFileToHeightAndWidth(req MustResizePngFileToHeightAndWidthRequest) {
	if req.OutFilePath == "" && req.InFilePath != "" {
		req.OutFilePath = req.InFilePath
	}
	fileContent := udwFile.MustReadFile(req.InFilePath)
	inputImage, _, err := image.Decode(bytes.NewBuffer(fileContent))
	if err != nil {
		panic(err)
	}

	if inputImage.Bounds().Dx() == req.Width && inputImage.Bounds().Dy() == req.Height {
		udwFile.MustWriteFile(req.OutFilePath, fileContent)
		return
	}
	MustResizeGoImageToPngFile(MustResizeGoImageToPngFileRequest{
		Height:      req.Height,
		Width:       req.Width,
		Image:       inputImage,
		OutFilePath: req.OutFilePath,
	})
}

func MustResizePngFileSelfToRate(filepath string, rate float64) {
	fileContent := udwFile.MustReadFile(filepath)
	inputImage, _, err := image.Decode(bytes.NewBuffer(fileContent))
	if err != nil {
		panic(err)
	}
	MustResizeGoImageToPngFile(MustResizeGoImageToPngFileRequest{
		Height:      int(float64(inputImage.Bounds().Dy()) * rate),
		Width:       int(float64(inputImage.Bounds().Dx()) * rate),
		Image:       inputImage,
		OutFilePath: filepath,
	})
}

func MustResizePngContentToSquare(pngContent []byte, width int) (outPngContent []byte) {
	inputImage, _, err := image.Decode(bytes.NewBuffer(pngContent))
	if err != nil {
		panic(err)
	}
	outImage := MustResizeGoImageToGoImage(MustResizeGoImageToGoImageRequest{
		Width:   width,
		Height:  width,
		InImage: inputImage,
	})
	buf := &bytes.Buffer{}
	err = png.Encode(buf, outImage)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
