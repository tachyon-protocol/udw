package udwBufio

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"testing"
)

func TestBufioReader(ot *testing.T) {
	{
		bufioReader := NewBufioReader(udwBytes.NewBufReader([]byte{1, 2, 3, 4}), 0)
		b, errMsg := bufioReader.ReadBySize(1)
		udwTest.Equal(errMsg, "")
		udwTest.Equal(b, []byte{1})

		b, errMsg = bufioReader.ReadBySize(1)
		udwTest.Equal(errMsg, "")
		udwTest.Equal(b, []byte{2})

		bufioReader.AddPos(-1)
		b, errMsg = bufioReader.ReadBySize(1)
		udwTest.Equal(errMsg, "")
		udwTest.Equal(b, []byte{2})

		b, errMsg = bufioReader.ReadBySize(3)
		udwTest.Equal(errMsg, "unexpected EOF")
		udwTest.Equal(b, nil)

	}
}

func TestBufioReader2(ot *testing.T) {
	bufioReader := NewBufioReader(udwBytes.NewBufReader([]byte{1, 2, 3}), 0)
	b, errMsg := bufioReader.ReadBySize(1)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(b, []byte{1})

	buf := make([]byte, 1024)
	nr, err := bufioReader.Read(buf)
	udwTest.Equal(err, nil)
	udwTest.Equal(buf[:2], []byte{2, 3})

	nr, err = bufioReader.Read(buf)
	udwTest.Equal(err, io.EOF)
	udwTest.Equal(nr, 0)
}

func TestBufioReader3(ot *testing.T) {
	totalBuf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		totalBuf[i] = byte(i)
	}
	bufioReader := NewBufioReader(udwBytes.NewBufReader(totalBuf), 0)
	b, errMsg := bufioReader.ReadBySize(1)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(b, []byte{0})

	buf := make([]byte, 700)
	nr, err := bufioReader.Read(buf)
	udwTest.Equal(err, nil)
	udwTest.Equal(nr, 700)
	for i := 0; i < 700; i++ {
		udwTest.Equal(buf[i], byte(i+1))
	}

	nr, err = bufioReader.Read(buf)
	udwTest.Equal(err, nil)
	udwTest.Equal(nr, 323)
	for i := 0; i < 323; i++ {
		udwTest.Equal(buf[i], byte(i+701), i)
	}

	nr, err = bufioReader.Read(buf)
	udwTest.Equal(err, io.EOF)
	udwTest.Equal(nr, 0)
}

func TestBufioReader4(ot *testing.T) {
	content := udwRand.MustCryptoRandBytes(32 * 1024)
	sr := &shortReader{
		r:             udwBytes.NewBufReader(content),
		thisReadLimit: 10,
	}
	bufioReader := NewBufioReader(sr, 0)
	b, errMsg := bufioReader.ReadBySize(1)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(b, content[0:1])
	b, errMsg = bufioReader.ReadBySize(100)
	udwTest.Equal(errMsg, "")
	udwTest.Equal(b, content[1:101])
	sr.shouldErr = errors.New("abc")
	b, errMsg = bufioReader.ReadBySize(100)
	udwTest.Equal(errMsg, "abc")
}

func TestBufioReader5(ot *testing.T) {
	{
		bufioReader := NewBufioReader(udwBytes.NewBufReader([]byte{1, 2, 3, 4}), 0)
		b, errMsg := bufioReader.ReadBySize(1)
		udwTest.Equal(errMsg, "")
		udwTest.Equal(b, []byte{1})
		udwTest.AssertPanicWithErrorMessage(func() {
			bufioReader.AddPos(-100)
		}, "z8xgzab2qr")
	}
}

type shortReader struct {
	r             io.Reader
	thisReadLimit int
	shouldErr     error
}

func (r *shortReader) Read(buf []byte) (nr int, err error) {
	if r.shouldErr != nil {
		return 0, r.shouldErr
	}
	if len(buf) > r.thisReadLimit {
		buf = buf[:r.thisReadLimit]
	}
	return r.r.Read(buf)
}
