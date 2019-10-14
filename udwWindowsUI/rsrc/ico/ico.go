package ico

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"sort"
)

const (
	BI_RGB = 0
)

type ICONDIR struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

type IconDirEntryCommon struct {
	Width      byte
	Height     byte
	ColorCount byte
	Reserved   byte
	Planes     uint16
	BitCount   uint16
	BytesInRes uint32
}

type ICONDIRENTRY struct {
	IconDirEntryCommon
	ImageOffset uint32
}

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type RGBQUAD struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

func skip(r io.Reader, n int64) error {
	_, err := io.CopyN(ioutil.Discard, r, n)
	return err
}

type icoOffset struct {
	n      int
	offset uint32
}

type rawico struct {
	icoinfo ICONDIRENTRY
	bmpinfo *BITMAPINFOHEADER
	idx     int
	data    []byte
}

type byOffsets []rawico

func (o byOffsets) Len() int           { return len(o) }
func (o byOffsets) Less(i, j int) bool { return o[i].icoinfo.ImageOffset < o[j].icoinfo.ImageOffset }
func (o byOffsets) Swap(i, j int) {
	tmp := o[i]
	o[i] = o[j]
	o[j] = tmp
}

type ICO struct {
	image.Image
}

func DecodeHeaders(r io.Reader) ([]ICONDIRENTRY, error) {
	var hdr ICONDIR
	err := binary.Read(r, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}
	if hdr.Reserved != 0 || hdr.Type != 1 {
		return nil, fmt.Errorf("bad magic number")
	}

	entries := make([]ICONDIRENTRY, hdr.Count)
	for i := 0; i < len(entries); i++ {
		err = binary.Read(r, binary.LittleEndian, &entries[i])
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

func unused_decodeAll(r io.Reader) ([]*ICO, error) {
	var hdr ICONDIR
	err := binary.Read(r, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}
	if hdr.Reserved != 0 || hdr.Type != 1 {
		return nil, fmt.Errorf("bad magic number")
	}

	raws := make([]rawico, hdr.Count)
	for i := 0; i < len(raws); i++ {
		err = binary.Read(r, binary.LittleEndian, &raws[i].icoinfo)
		if err != nil {
			return nil, err
		}
		raws[i].idx = i
	}

	sort.Sort(byOffsets(raws))

	offset := uint32(binary.Size(&hdr) + len(raws)*binary.Size(ICONDIRENTRY{}))
	for i := 0; i < len(raws); i++ {
		err = skip(r, int64(raws[i].icoinfo.ImageOffset-offset))
		if err != nil {
			return nil, err
		}
		offset = raws[i].icoinfo.ImageOffset

		raws[i].bmpinfo = &BITMAPINFOHEADER{}
		err = binary.Read(r, binary.LittleEndian, raws[i].bmpinfo)
		if err != nil {
			return nil, err
		}

		err = skip(r, int64(raws[i].bmpinfo.Size-uint32(binary.Size(BITMAPINFOHEADER{}))))
		if err != nil {
			return nil, err
		}
		raws[i].data = make([]byte, raws[i].icoinfo.BytesInRes-raws[i].bmpinfo.Size)
		_, err = io.ReadFull(r, raws[i].data)
		if err != nil {
			return nil, err
		}
	}

	icos := make([]*ICO, len(raws))
	for i := 0; i < len(raws); i++ {
		fmt.Println(i)
		icos[raws[i].idx], err = decode(raws[i].bmpinfo, &raws[i].icoinfo, raws[i].data)
		if err != nil {
			return nil, err
		}
	}
	return icos, nil
}

func decode(info *BITMAPINFOHEADER, icoinfo *ICONDIRENTRY, data []byte) (*ICO, error) {
	if info.Compression != BI_RGB {
		return nil, fmt.Errorf("ICO compression not supported (got %d)", info.Compression)
	}

	r := bytes.NewBuffer(data)

	bottomup := info.Height > 0
	if !bottomup {
		info.Height = -info.Height
	}

	switch info.BitCount {
	case 8:
		ncol := int(icoinfo.ColorCount)
		if ncol == 0 {
			ncol = 256
		}

		pal := make(color.Palette, ncol)
		for i := 0; i < ncol; i++ {
			var rgb RGBQUAD
			err := binary.Read(r, binary.LittleEndian, &rgb)
			if err != nil {
				return nil, err
			}
			pal[i] = color.NRGBA{R: rgb.Red, G: rgb.Green, B: rgb.Blue, A: 0xff}
		}
		fmt.Println(pal)

		fmt.Println(info.SizeImage, len(data)-binary.Size(RGBQUAD{})*len(pal), info.Width, info.Height)

	default:
		return nil, fmt.Errorf("unsupported ICO bit depth (BitCount) %d", info.BitCount)
	}

	return nil, nil
}
