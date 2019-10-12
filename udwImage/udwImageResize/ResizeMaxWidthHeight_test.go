package udwImageResize

import (
	"github.com/tachyon-protocol/udw/udwImage"
	"github.com/tachyon-protocol/udw/udwTest"
	"image"
	"strconv"
	"testing"
)

func TestMustMaxWidthHeightResizeToPng(t *testing.T) {
	type tCas struct {
		inX  int
		inY  int
		MaxW int
		maxH int
		outX int
		outY int
	}
	for _, tCas := range []tCas{
		{100, 100, 10, 10, 10, 10},
		{100, 100, 50, 10, 10, 10},
		{100, 200, 10, 10, 5, 10},
		{100, 200, 10, 20, 10, 20},
		{100, 1, 10, 10, 10, 1},
		{100, 200, 100, 0, 100, 200},
		{100, 200, 200, 0, 100, 200},
		{100, 200, 10, 0, 10, 20},
		{100, 200, 0, 20, 10, 20},
	} {
		thisImg := image.NewRGBA(image.Rect(0, 0, tCas.inX, tCas.inY))
		b1 := udwImage.MustPngEncodeFromGoImageToBytes(thisImg)
		b2 := MustMaxWidthHeightResizeToPng(MustMaxWidthHeightRequest{
			InContent: b1,
			MaxWidth:  tCas.MaxW,
			MaxHeight: tCas.maxH,
		})
		img2 := udwImage.MustPngDecodeFromBuf(b2)
		expectS := strconv.Itoa(tCas.outX) + "x" + strconv.Itoa(tCas.outY)
		fmtS := strconv.Itoa(img2.Bounds().Dx()) + "x" + strconv.Itoa(img2.Bounds().Dy())
		udwTest.Equal(expectS, fmtS)
	}
	{
		thisImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
		b1 := udwImage.MustPngEncodeFromGoImageToBytes(thisImg)
		udwTest.AssertPanicWithErrorMessage(func() {
			MustMaxWidthHeightResizeToPng(MustMaxWidthHeightRequest{
				InContent: b1,
			})
		}, "ced365e92v")
	}
}
