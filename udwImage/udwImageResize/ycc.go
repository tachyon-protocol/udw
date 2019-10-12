package udwImageResize

import (
	"image"
	"image/color"
)

type ycc struct {
	Pix []uint8

	Stride int

	Rect image.Rectangle

	SubsampleRatio image.YCbCrSubsampleRatio
}

func (p *ycc) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *ycc) Bounds() image.Rectangle {
	return p.Rect
}

func (p *ycc) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *ycc) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}
	i := p.PixOffset(x, y)
	return color.YCbCr{
		p.Pix[i+0],
		p.Pix[i+1],
		p.Pix[i+2],
	}
}

func (p *ycc) Opaque() bool {
	return true
}

func (p *ycc) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &ycc{SubsampleRatio: p.SubsampleRatio}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &ycc{
		Pix:            p.Pix[i:],
		Stride:         p.Stride,
		Rect:           r,
		SubsampleRatio: p.SubsampleRatio,
	}
}

func newYCC(r image.Rectangle, s image.YCbCrSubsampleRatio) *ycc {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 3*w*h)
	return &ycc{Pix: buf, Stride: 3 * w, Rect: r, SubsampleRatio: s}
}

const (
	ycbcrSubsampleRatio444 image.YCbCrSubsampleRatio = iota
	ycbcrSubsampleRatio422
	ycbcrSubsampleRatio420
	ycbcrSubsampleRatio440
	ycbcrSubsampleRatio411
	ycbcrSubsampleRatio410
)

func (p *ycc) YCbCr() *image.YCbCr {
	ycbcr := image.NewYCbCr(p.Rect, p.SubsampleRatio)
	switch ycbcr.SubsampleRatio {
	case ycbcrSubsampleRatio422:
		return p.ycbcr422(ycbcr)
	case ycbcrSubsampleRatio420:
		return p.ycbcr420(ycbcr)
	case ycbcrSubsampleRatio440:
		return p.ycbcr440(ycbcr)
	case ycbcrSubsampleRatio444:
		return p.ycbcr444(ycbcr)
	case ycbcrSubsampleRatio411:
		return p.ycbcr411(ycbcr)
	case ycbcrSubsampleRatio410:
		return p.ycbcr410(ycbcr)
	}
	return ycbcr
}

func imageYCbCrToYCC(in *image.YCbCr) *ycc {
	w, h := in.Rect.Dx(), in.Rect.Dy()
	p := ycc{
		Pix:            make([]uint8, 3*w*h),
		Stride:         3 * w,
		Rect:           image.Rect(0, 0, w, h),
		SubsampleRatio: in.SubsampleRatio,
	}
	switch in.SubsampleRatio {
	case ycbcrSubsampleRatio422:
		return convertToYCC422(in, &p)
	case ycbcrSubsampleRatio420:
		return convertToYCC420(in, &p)
	case ycbcrSubsampleRatio440:
		return convertToYCC440(in, &p)
	case ycbcrSubsampleRatio444:
		return convertToYCC444(in, &p)
	case ycbcrSubsampleRatio411:
		return convertToYCC411(in, &p)
	case ycbcrSubsampleRatio410:
		return convertToYCC410(in, &p)
	}
	return &p
}

func (p *ycc) ycbcr422(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := y * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x/2
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func (p *ycc) ycbcr420(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := (y / 2) * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x/2
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func (p *ycc) ycbcr440(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := (y / 2) * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func (p *ycc) ycbcr444(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := y * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func (p *ycc) ycbcr411(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := y * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x/4
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func (p *ycc) ycbcr410(ycbcr *image.YCbCr) *image.YCbCr {
	var off int
	Pix := p.Pix
	Y := ycbcr.Y
	Cb := ycbcr.Cb
	Cr := ycbcr.Cr
	for y := 0; y < ycbcr.Rect.Max.Y-ycbcr.Rect.Min.Y; y++ {
		yy := y * ycbcr.YStride
		cy := (y / 2) * ycbcr.CStride
		for x := 0; x < ycbcr.Rect.Max.X-ycbcr.Rect.Min.X; x++ {
			ci := cy + x/4
			Y[yy+x] = Pix[off+0]
			Cb[ci] = Pix[off+1]
			Cr[ci] = Pix[off+2]
			off += 3
		}
	}
	return ycbcr
}

func convertToYCC422(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := y * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x/2
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}

func convertToYCC420(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := (y / 2) * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x/2
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}

func convertToYCC440(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := (y / 2) * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}

func convertToYCC444(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := y * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}

func convertToYCC411(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := y * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x/4
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}

func convertToYCC410(in *image.YCbCr, p *ycc) *ycc {
	var off int
	Pix := p.Pix
	Y := in.Y
	Cb := in.Cb
	Cr := in.Cr
	for y := 0; y < in.Rect.Max.Y-in.Rect.Min.Y; y++ {
		yy := y * in.YStride
		cy := (y / 2) * in.CStride
		for x := 0; x < in.Rect.Max.X-in.Rect.Min.X; x++ {
			ci := cy + x/4
			Pix[off+0] = Y[yy+x]
			Pix[off+1] = Cb[ci]
			Pix[off+2] = Cr[ci]
			off += 3
		}
	}
	return p
}
