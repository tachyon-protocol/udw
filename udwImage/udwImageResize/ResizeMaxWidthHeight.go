package udwImageResize

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"image"
	"image/png"
)

type MustMaxWidthHeightRequest struct {
	InContent []byte

	MaxWidth            int
	MaxHeight           int
	PngCompressionLevel png.CompressionLevel
}

func MustMaxWidthHeightResizeToPng(req MustMaxWidthHeightRequest) (b []byte) {
	if req.MaxWidth == 0 && req.MaxHeight == 0 {
		panic("ced365e92v")
	}
	inputImage, _, err := image.Decode(bytes.NewBuffer(req.InContent))
	if err != nil {
		panic(err)
	}
	inW := inputImage.Bounds().Dx()
	inH := inputImage.Bounds().Dy()
	outBuf := &udwBytes.BufWriter{}
	if (req.MaxHeight == 0 || inH <= req.MaxHeight) && (req.MaxWidth == 0 || inW <= req.MaxWidth) {
		encoder := png.Encoder{
			CompressionLevel: req.PngCompressionLevel,
		}
		err = encoder.Encode(outBuf, inputImage)
		if err != nil {
			panic(err)
		}
		return outBuf.GetBytes()
	}
	outHF := float64(inH)
	outWF := float64(inW)
	if req.MaxHeight > 0 && outHF > float64(req.MaxHeight) {
		outWF = float64(outWF) * float64(req.MaxHeight) / float64(outHF)
		outHF = float64(req.MaxHeight)
	}
	if req.MaxWidth > 0 && outWF > float64(req.MaxWidth) {
		outHF = float64(outHF) * float64(req.MaxWidth) / float64(outWF)
		outWF = float64(req.MaxWidth)
	}
	if outHF < 1 {
		outHF = 1
	}
	if outWF < 1 {
		outWF = 1
	}
	image2 := Resize(uint(outWF), uint(outHF), inputImage, Lanczos3)
	encoder := png.Encoder{
		CompressionLevel: req.PngCompressionLevel,
	}
	err = encoder.Encode(outBuf, image2)
	if err != nil {
		panic(err)
	}
	return outBuf.GetBytes()
}

func MustReEncodeToPng(InContent []byte) (b []byte) {
	inputImage, _, err := image.Decode(bytes.NewBuffer(InContent))
	if err != nil {
		panic(err)
	}
	outBuf := &udwBytes.BufWriter{}
	err = png.Encode(outBuf, inputImage)
	if err != nil {
		panic(err)
	}
	return outBuf.GetBytes()
}
