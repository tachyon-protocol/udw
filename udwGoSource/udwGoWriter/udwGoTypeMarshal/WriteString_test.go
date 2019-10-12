package udwGoTypeMarshal_test

import (
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal/udwWriteStringToGolangFuzzLib"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwTest"
	"go/format"
	"testing"
	"testing/quick"
)

const SlowWriteStringCheck = false

func TestWriteStringToGolang(ot *testing.T) {
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolang("1"), "`1`")
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolang("`"), `"`+"`"+`"`)
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolang("\xef\xbb\xbf"), `"\xef\xbb\xbf"`)
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolang(string([]byte("\xef\xbb\xbf"))), `"\xef\xbb\xbf"`)
	writeStringCheck("\xef\xbb\xbf")

	for i := 0; i < 256; i++ {
		writeStringCheck(string([]byte{byte(i)}))
	}

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			writeStringCheck(string([]byte{byte(i), byte(j)}))
		}
	}
	if SlowWriteStringCheck {

		for i := 128; i < 256; i++ {
			for j := 128; j < 256; j++ {
				for k := 128; k < 256; k++ {
					writeStringCheck(string([]byte{byte(i), byte(j), byte(k)}))
				}
			}
		}
	}
	err := quick.Check(func(s []byte) bool {
		writeStringCheck(string(s))
		return true
	}, &quick.Config{
		MaxCount: 10000,
	})
	if err != nil {
		panic(err)
	}
}

func writeStringCheck(s string) {
	out := udwGoTypeMarshal.WriteStringToGolang(s)
	outS := udwGoTypeMarshal.MustReadGoStringFromString(out)
	udwTest.Equal([]byte(outS), []byte(s))
	goFile := `package main

var a = ` + out + `
`
	_, err := format.Source([]byte(goFile))
	if err != nil {
		panic("fail at s:[" + hex.Dump([]byte(s)) + "] [" + out + "] err [" + err.Error() + "]")
	}
}

func TestWriteStringToGolangASCII(ot *testing.T) {
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolangASCII("1"), "`1`")
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolangASCII("`"), `"`+"`"+`"`)
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolangASCII("\xef\xbb\xbf"), `"\xef\xbb\xbf"`)
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolangASCII("中文"), `"\xe4\xb8\xad\xe6\x96\x87"`)
	udwTest.Equal(udwGoTypeMarshal.WriteStringToGolangASCII(string([]byte("\xef\xbb\xbf"))), `"\xef\xbb\xbf"`)
	writeStringCheckASCII("\xef\xbb\xbf")
}

func writeStringCheckASCII(s string) {
	out := udwGoTypeMarshal.WriteStringToGolangASCII(s)
	outS := udwGoTypeMarshal.MustReadGoStringFromString(out)
	udwTest.Equal([]byte(outS), []byte(s))
	for i := 0; i < len(out); i++ {
		b := out[i]
		if !(b >= 0x20 && b <= 0x7e) {
			panic(fmt.Errorf("[writeStringCheckASCII] fail at [%s] is not ACSII safe",
				udwHex.EncodeStringToString(s)))
		}
	}
	goFile := `package main

var a = ` + out + `
`
	_, err := format.Source([]byte(goFile))
	if err != nil {
		panic(fmt.Errorf("[writeStringCheckASCII] fail at %s err %s", out, err.Error()))
	}
}

func TestFuzz(ot *testing.T) {
	udwWriteStringToGolangFuzzLib.Fuzz([]byte("1"))
	udwWriteStringToGolangFuzzLib.Fuzz([]byte("0\r"))
	udwWriteStringToGolangFuzzLib.Fuzz([]byte("\n\r"))
}
