package udwImageResize

import (
	"github.com/tachyon-protocol/udw/udwImage"
	"github.com/tachyon-protocol/udw/udwTest"
	"image"
	"strconv"
	"testing"
)

func TestMustResizeImageToCenterToPng(t *testing.T) {
	type tCas struct {
		inX          int
		inY          int
		reqW         int
		reqH         int
		disableLager bool
		outX         int
		outY         int
	}
	for _, tCas := range []tCas{
		{100, 100, 10, 10, false, 10, 10},
		{100, 100, 50, 10, false, 50, 10},
		{10, 10, 20, 10, false, 20, 10},
		{100, 1, 10, 10, false, 10, 10},
		{2, 2, 10, 1, false, 10, 1},
		{1, 1, 10, 1, false, 10, 1},
		{1, 1, 10, 1, true, 1, 1},
		{10, 10, 100, 100, true, 10, 10},
		{10, 5, 100, 100, true, 5, 5},
	} {
		thisImg := image.NewRGBA(image.Rect(0, 0, tCas.inX, tCas.inY))
		b1 := udwImage.MustPngEncodeFromGoImageToBytes(thisImg)
		b2 := MustResizeImageToCenterToPng(MustResizeImageToCenterToPngReq{
			InContent:     b1,
			Width:         tCas.reqW,
			Height:        tCas.reqH,
			DisableLarger: tCas.disableLager,
		})
		img2 := udwImage.MustPngDecodeFromBuf(b2)
		expectS := strconv.Itoa(tCas.outX) + "x" + strconv.Itoa(tCas.outY)
		fmtS := strconv.Itoa(img2.Bounds().Dx()) + "x" + strconv.Itoa(img2.Bounds().Dy())
		udwTest.Equal(expectS, fmtS)
	}
}
