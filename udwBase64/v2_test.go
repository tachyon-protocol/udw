package udwBase64

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestV2(ot *testing.T) {
	for _, testCase := range []struct {
		in  string
		out string
	}{
		{"", ""},
		{"A", "QQ"},
		{"\x00", "AA"},
		{"\x00\x00", "AAA"},
		{"\x00\x00\x00", "AAAA"},
		{"\xff\xff\xff", "____"},
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l-"},
		{"\x14\xfb\x9c\x03\xd9", "FPucA9k"},
		{"\x14\xfb\x9c\x03", "FPucAw"},

		{"", ""},
		{"f", "Zg"},
		{"fo", "Zm8"},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg"},
		{"fooba", "Zm9vYmE"},
		{"foobar", "Zm9vYmFy"},

		{"sure.", "c3VyZS4"},
		{"sure", "c3VyZQ"},
		{"sur", "c3Vy"},
		{"su", "c3U"},
		{"leasure.", "bGVhc3VyZS4"},
		{"easure.", "ZWFzdXJlLg"},
		{"asure.", "YXN1cmUu"},
		{"sure.", "c3VyZS4"},
		{"Twas brillig, and the slithy toves", "VHdhcyBicmlsbGlnLCBhbmQgdGhlIHNsaXRoeSB0b3Zlcw"},
	} {
		udwTest.Equal(EncodeByteToStringV2([]byte(testCase.in)), testCase.out)
		thisIn, err := DecodeStringToByteV2(testCase.out)
		udwTest.Equal(err, nil)
		udwTest.Equal(thisIn, []byte(testCase.in))
	}
}

func BenchmarkEncodeToStringV2(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		EncodeByteToStringV2(data)
	}
}

func BenchmarkDecodeStringV2(b *testing.B) {
	data := EncodeByteToStringV2(make([]byte, 8192))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		DecodeStringToByteV2(data)
	}
}
