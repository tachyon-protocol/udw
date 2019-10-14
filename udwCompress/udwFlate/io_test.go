package udwFlate

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"io/ioutil"
	"testing"
)

func TestCompressMustNewWriter(t *testing.T) {
	inB := bytes.Repeat([]byte{1}, 1024*1024)
	buf := &udwBytes.BufWriter{}
	w2 := CompressMustNewWriter(buf)
	_, err := w2.Write(inB)
	udwTest.Equal(err, nil)
	_, err = w2.Write(inB)
	udwTest.Equal(err, nil)
	w2.Close()
	udwTest.Equal(MustFlateUnCompress(buf.GetBytes()), bytes.Repeat([]byte{1}, 2*1024*1024))
}

func TestCompressMustNewReader(t *testing.T) {
	size := 1024*1024 + 1
	inB := bytes.Repeat([]byte{1}, size)
	r1 := udwBytes.NewBufReader(inB)

	r2 := CompressMustNewReader(r1)
	b, err := ioutil.ReadAll(r2)
	udwTest.Equal(err, nil)
	udwTest.Equal(MustFlateUnCompress(b), bytes.Repeat([]byte{1}, size))
}
