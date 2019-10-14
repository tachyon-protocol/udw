package udwImage

import (
	"github.com/tachyon-protocol/udw/udwHex"
	"image/color"
)

func ColorToRGBHex(c color.Color) string {
	rgba8 := color.RGBAModel.Convert(c).(color.RGBA)
	outBuf := make([]byte, 7)
	outBuf[0] = '#'
	outBuf[1] = udwHex.HexTableLower[rgba8.R>>4&0xf]
	outBuf[2] = udwHex.HexTableLower[rgba8.R&0xf]
	outBuf[3] = udwHex.HexTableLower[rgba8.G>>4&0xf]
	outBuf[4] = udwHex.HexTableLower[rgba8.G&0xf]
	outBuf[5] = udwHex.HexTableLower[rgba8.B>>4&0xf]
	outBuf[6] = udwHex.HexTableLower[rgba8.B&0xf]
	return string(outBuf)
}

func ColorToRGBAHex(c color.Color) string {
	rgba8 := color.RGBAModel.Convert(c).(color.RGBA)
	outBuf := make([]byte, 9)
	outBuf[0] = '#'
	outBuf[1] = udwHex.HexTableLower[rgba8.R>>4&0xf]
	outBuf[2] = udwHex.HexTableLower[rgba8.R&0xf]
	outBuf[3] = udwHex.HexTableLower[rgba8.G>>4&0xf]
	outBuf[4] = udwHex.HexTableLower[rgba8.G&0xf]
	outBuf[5] = udwHex.HexTableLower[rgba8.B>>4&0xf]
	outBuf[6] = udwHex.HexTableLower[rgba8.B&0xf]
	outBuf[7] = udwHex.HexTableLower[rgba8.A>>4&0xf]
	outBuf[8] = udwHex.HexTableLower[rgba8.A&0xf]
	return string(outBuf)
}

func MustColorRGBAHexToObj(hex string) (c color.Color) {
	if len(hex) != 9 {
		panic("[MustRGBAHexToColor] g347a7dadj")
	}
	if hex[0] != '#' {
		panic("[MustRGBAHexToColor] q6e42s95a4")
	}
	r1, ok := udwHex.FromHexChar(hex[1])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	r2, ok := udwHex.FromHexChar(hex[2])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	g1, ok := udwHex.FromHexChar(hex[3])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	g2, ok := udwHex.FromHexChar(hex[4])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	b1, ok := udwHex.FromHexChar(hex[5])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	b2, ok := udwHex.FromHexChar(hex[6])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	a1, ok := udwHex.FromHexChar(hex[7])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	a2, ok := udwHex.FromHexChar(hex[8])
	if ok == false {
		panic("[MustRGBAHexToColor] kcnftwjf92")
	}
	return color.RGBA{(r1<<4 | r2), (g1<<4 | g2), (b1<<4 | b2), (a1<<4 | a2)}
}

func GetColorBlack() (c color.Color) {
	return color.Black
}

func GetColorWhite() (c color.Color) {
	return color.White
}

func GetColorRed() (c color.Color) {
	return color.RGBA{R: 0xff, A: 0xff}
}

func GetColorGreen() (c color.Color) {
	return color.RGBA{G: 0xff, A: 0xff}
}
func GetColorBlue() (c color.Color) {
	return color.RGBA{B: 0xff, A: 0xff}
}
