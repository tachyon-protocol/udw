package udwPcapIpPacket

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestReader(ot *testing.T) {
	inItem := Item{T: time.Now().Truncate(time.Microsecond), IpPacket: []byte{0x45, 0x0,
		0x0, 0x3b,
		0xc3, 0xb8,
		0x0, 0x0,
		0xff,
		0x11,
		0x66, 0xfe,
		0xac, 0x15, 0x0, 0x1,
		0x72, 0x72, 0x72, 0x72,

		0xe4, 0xfa,
		0x0, 0x35,
		0x0, 0x27,
		0x4, 0x3b,
		0xe0, 0x2a, 0x1, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd, 0x78, 0x78, 0x63, 0x6c, 0x7a, 0x72, 0x74, 0x61, 0x64, 0x64, 0x6a, 0x77, 0x75, 0x0, 0x0, 0x1, 0x0, 0x1},
	}
	_buf := &udwBytes.BufWriter{}
	MarshalToBuffer(_buf, []Item{
		inItem, inItem})
	num := 0
	errMsg := ReadPcapToIpPacket(_buf.GetBytes(), func(item Item) {
		udwTest.Ok(item.T.Equal(inItem.T), item.T, inItem.T)
		num++
	})
	udwTest.Equal(errMsg, "")
	udwTest.Equal(num, 2)
}
