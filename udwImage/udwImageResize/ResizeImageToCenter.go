package udwImageResize

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"image"
	"image/draw"
	"image/png"
)

type ResizeImageToCenterReq struct {
	InImg         image.Image
	Width         int
	Height        int
	DisableLarger bool
}

func MustResizeImageToCenter(req ResizeImageToCenterReq) image.Image {
	inW := req.InImg.Bounds().Dx()
	inH := req.InImg.Bounds().Dy()
	if inW == req.Width && inH == req.Height {
		return req.InImg
	}
	inRate := float64(inW) / float64(inH)
	outRate := float64(req.Width) / float64(req.Height)
	if inRate == outRate {
		if req.DisableLarger && inW <= req.Width {
			return req.InImg
		}
		return Resize(uint(req.Width), uint(req.Height), req.InImg, Lanczos3)
	}
	cutX0F := float64(0)
	cutY0F := float64(0)
	cutX1F := float64(inW)
	cutY1F := float64(inH)
	if inRate > outRate {

		afterCutW := float64(inH) * outRate
		cutX0F = (float64(inW) - afterCutW) / 2
		cutX1F = float64(inW) - (float64(inW)-afterCutW)/2
	} else {

		afterCutH := float64(inW) / outRate
		cutY0F = (float64(inH) - afterCutH) / 2
		cutY1F = float64(inH) - (float64(inH)-afterCutH)/2
	}
	cutX0 := int(cutX0F)
	cutX1 := int(cutX1F)
	cutY0 := int(cutY0F)
	cutY1 := int(cutY1F)
	if cutX1-cutX0 <= 0 {
		cutX1 = cutX0 + 1
	}
	if cutY1-cutY0 <= 0 {
		cutY1 = cutY0 + 1
	}
	cutW := int(cutX1 - cutX0)
	cutH := int(cutY1 - cutY0)
	inImg2 := image.NewRGBA(image.Rect(0, 0, cutW, cutH))
	draw.Draw(inImg2, image.Rect(0, 0, cutW, cutH), req.InImg, image.Pt(cutX0, cutY0), draw.Over)
	if req.DisableLarger && cutW <= req.Width {
		return inImg2
	}
	return Resize(uint(req.Width), uint(req.Height), inImg2, Lanczos3)
}

type MustResizeImageToCenterToPngReq struct {
	InContent           []byte
	Width               int
	Height              int
	DisableLarger       bool
	PngCompressionLevel png.CompressionLevel
}

func MustResizeImageToCenterToPng(req MustResizeImageToCenterToPngReq) (b []byte) {
	inputImage, _, err := image.Decode(bytes.NewBuffer(req.InContent))
	if err != nil {
		panic(err)
	}
	outImg := MustResizeImageToCenter(ResizeImageToCenterReq{
		InImg:         inputImage,
		Width:         req.Width,
		Height:        req.Height,
		DisableLarger: req.DisableLarger,
	})
	outBuf := &udwBytes.BufWriter{}
	encoder := png.Encoder{
		CompressionLevel: req.PngCompressionLevel,
	}
	err = encoder.Encode(outBuf, outImg)
	if err != nil {
		panic(err)
	}
	return outBuf.GetBytes()
}
