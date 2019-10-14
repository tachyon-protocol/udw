package udwBinary

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestReadByteSliceWithUint32Len(t *testing.T) {
	buf := &bytes.Buffer{}

	err := WriteByteSliceWithUint32Len(buf, []byte{})
	udwTest.Equal(err, nil)

	b, err := ReadByteSliceWithUint32Len(buf)
	udwTest.Equal(err, nil)
	udwTest.Equal(len(b), 0)

	for _, size := range []int{
		32 * 1024,
		48 * 1024,
		100 * 1024,
	} {
		buf = &bytes.Buffer{}
		inB := bytes.Repeat([]byte{1}, size)
		err = WriteByteSliceWithUint32Len(buf, inB)
		udwTest.Equal(err, nil)

		b, err = ReadByteSliceWithUint32Len(buf)
		udwTest.Equal(err, nil)
		udwTest.Equal(len(b), size)
		udwTest.Equal(b, inB)
	}

	buf = bytes.NewBuffer([]byte{255, 255, 255, 255, 0})
	b, err = ReadByteSliceWithUint32Len(buf)
	udwTest.Ok(err != nil)

}
