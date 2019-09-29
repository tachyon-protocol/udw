package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestBufWriterTrimSuffixSpace(ot *testing.T) {
	w := BufWriter{}
	w.TrimSuffixSpace()
	udwTest.Equal(w.GetBytes(), nil)
	w.WriteByte(' ')
	w.TrimSuffixSpace()
	udwTest.Equal(w.GetLen(), 0)
	w.WriteString("1 ")
	w.TrimSuffixSpace()
	udwTest.Equal(w.GetLen(), 1)
	w.WriteString(" 1 \t \n")
	w.TrimSuffixSpace()
	udwTest.Equal(w.GetString(), "1 1")

	w = BufWriter{}
	w.WriteString("1")
	w.TrimSuffixSpace()
	udwTest.Equal(w.GetString(), "1")
}
