package udwImageResize

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"image"
	"image/png"
)

type MustResizeGoImageToPngFileRequest struct {
	Height      int
	Width       int
	Image       image.Image
	OutFilePath string
}

func MustResizeGoImageToPngFile(req MustResizeGoImageToPngFileRequest) {
	dst := MustResizeGoImageToGoImage(MustResizeGoImageToGoImageRequest{
		Height:  req.Height,
		Width:   req.Width,
		InImage: req.Image,
	})
	buf := &bytes.Buffer{}
	err := png.Encode(buf, dst)
	if err != nil {
		panic(err)
	}
	udwFile.MustWriteFileWithMkdir(req.OutFilePath, buf.Bytes())
}

type MustResizeGoImageToGoImageRequest struct {
	Height  int
	Width   int
	InImage image.Image
}

func MustResizeGoImageToGoImage(req MustResizeGoImageToGoImageRequest) image.Image {
	return Resize(uint(req.Width), uint(req.Height), req.InImage, NearestNeighbor)
}
