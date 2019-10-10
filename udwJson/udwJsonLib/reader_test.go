package udwJsonLib

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustReadJsonByteSlice(ot *testing.T) {

	for _, cas := range []struct {
		in  string
		out []byte
	}{
		{`"ew=="`, []byte{123}},
		{`"ewE="`, []byte{123, 1}},
		{`"AA=="`, []byte{0}},
		{`"AAA="`, []byte{0, 0}},
		{`"AAAA"`, []byte{0, 0, 0}},
		{`"AAAAAA=="`, []byte{0, 0, 0, 0}},
		{`""`, []byte{}},
		{`null`, nil},
		{`[0,1,2]`, []byte{0, 1, 2}},
	} {
		ctx := NewContextFromBuffer([]byte(cas.in))
		b := MustReadJsonByteSlice(ctx)
		udwTest.Equal(b, cas.out, cas.in)
	}
}
