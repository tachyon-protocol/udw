package udwBytesEncode

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMapStringString(ot *testing.T) {
	{
		b := MapStringStringMarshal(nil)
		udwTest.Equal(b, []byte{0})
		m, ok := MapStringStringUnmarshal([]byte{0})
		udwTest.Equal(ok, true)
		udwTest.Equal(len(m), 0)
	}
	{
		b := MapStringStringMarshal(map[string]string{
			"k1": "v1",
		})
		udwTest.Equal(b, []byte{0x01, 0x02, 0x6b, 0x31, 0x02, 0x76, 0x31})
		m, ok := MapStringStringUnmarshal([]byte{0x01, 0x02, 0x6b, 0x31, 0x02, 0x76, 0x31})
		udwTest.Equal(ok, true)
		udwTest.Equal(len(m), 1)
		udwTest.Equal(m["k1"], "v1")
	}
	{
		m, ok := MapStringStringUnmarshal([]byte{10, 0x02, 0x6b, 0x31, 0x02, 0x76, 0x31})
		udwTest.Equal(ok, false)
		udwTest.Equal(m, nil)
	}
	{
		m, ok := MapStringStringUnmarshal([]byte{})
		udwTest.Equal(ok, false)
		udwTest.Equal(m, nil)
	}
	{
		m, ok := MapStringStringUnmarshal([]byte{0x01, 0x02, 0x6b, 0x31, 0x02, 0x76})
		udwTest.Equal(ok, false)
		udwTest.Equal(m, nil)
	}
}
