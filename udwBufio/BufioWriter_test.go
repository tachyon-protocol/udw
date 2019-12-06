package udwBufio

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBufioWriter(ot *testing.T) {
	{
		w := &udwBytes.BufWriter{}
		w1 := NewBufioWriter(w, 0)
		w1.Write_(bytes.Repeat([]byte{1}, 1024*33))
		errMsg := w1.SoftFlush()
		udwErr.PanicIfErrorMsg(errMsg)
		udwTest.Equal(w.GetLen(), 1024*33)
	}
	{
		w := &udwBytes.BufWriter{}
		w1 := NewBufioWriter(w, 0)
		w1.WriteByte_(1)
		errMsg := w1.SoftFlush()
		udwErr.PanicIfErrorMsg(errMsg)
		udwTest.Equal(w.GetLen(), 0)
		errMsg = w1.Flush()
		udwErr.PanicIfErrorMsg(errMsg)
		udwTest.Equal(w.GetLen(), 1)
		udwTest.Equal(w.GetBytes(), []byte{1})
	}
}
