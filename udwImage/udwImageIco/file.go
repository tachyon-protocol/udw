package udwImageIco

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwCache"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwImage"
	"github.com/tachyon-protocol/udw/udwImage/udwImageResize"
)

type MustImageToIcoFileToFileRequest struct {
	InPath               string
	OutPath              string
	ResizeHeightAndWidth int
}

func MustEncodeOneImageFromFileToFileWithCache(req MustImageToIcoFileToFileRequest) {
	udwCache.MustMd5FileChangeCache("udwImageIco_"+req.OutPath, []string{
		req.OutPath,
		req.InPath,
		"src/github.com/tachyon-protocol/udw/udwImage/udwImageIco",
	}, func() {
		im := udwImage.MustDecodeImageFromFile(req.InPath)
		if req.ResizeHeightAndWidth > 0 {
			im = udwImageResize.MustResizeGoImageToGoImage(udwImageResize.MustResizeGoImageToGoImageRequest{
				Height:  req.ResizeHeightAndWidth,
				Width:   req.ResizeHeightAndWidth,
				InImage: im,
			})
		}
		_buf := &bytes.Buffer{}
		err := EncodeOneImageToWriter(_buf, im)
		if err != nil {
			panic(err)
		}
		udwFile.MustWriteFileWithMkdir(req.OutPath, _buf.Bytes())
	})

}

func MustEncodePngContentToIcoContent(pngContent []byte) (icoContent []byte) {
	im := udwImage.MustDecodeImageFromByte(pngContent)
	_buf := &bytes.Buffer{}
	err := EncodeOneImageToWriter(_buf, im)
	if err != nil {
		panic(err)
	}
	return _buf.Bytes()
}
