package udwImage

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"image/color"
	"testing"
)

func TestColor(t *testing.T) {
	udwTest.Equal(ColorToRGBAHex(color.RGBA{0x12, 0x34, 0x56, 0x78}), "#12345678")
	udwTest.Equal(ColorToRGBHex(color.RGBA{0x12, 0x34, 0x56, 0x78}), "#123456")
	s := "#12345678"
	udwTest.Equal(ColorToRGBAHex(MustColorRGBAHexToObj(s)), s)

	udwTest.Equal(ColorToRGBAHex(GetColorBlack()), "#000000ff")
	udwTest.Equal(ColorToRGBAHex(GetColorWhite()), "#ffffffff")
	udwTest.Equal(ColorToRGBAHex(GetColorRed()), "#ff0000ff")
	udwTest.Equal(ColorToRGBAHex(GetColorGreen()), "#00ff00ff")
	udwTest.Equal(ColorToRGBAHex(GetColorBlue()), "#0000ffff")
}
