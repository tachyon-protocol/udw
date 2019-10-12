package udwFlate

import (
	"github.com/tachyon-protocol/udw/udwCompress/kkcflate"
	"github.com/tachyon-protocol/udw/udwIo"
	"io"
)

func CompressMustNewWriter(w io.Writer) io.WriteCloser {
	flateW, err := kkcflate.NewWriter(w, 4)
	if err != nil {
		panic(err)
	}
	return flateW
}

func CompressMustNewReader(r1 io.Reader) (r2 io.Reader) {
	return udwIo.NewWriterWrapToReaderWrap(CompressMustNewWriter, r1)
}

func UncompressMustNewReader(r io.Reader) (r2 io.ReadCloser) {
	return kkcflate.NewReader(r)
}

func NewReader(r io.Reader) (r2 io.ReadCloser) {
	return kkcflate.NewReader(r)
}
