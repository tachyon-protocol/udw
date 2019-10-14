package udwVpnPacket

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

var Ipv4TcpOptionTestCase = string([]byte{

	0x45, 0x00, 0x00, 0x40, 0x9f, 0x60, 0x40, 0x00, 0x40, 0x06, 0x34, 0xc4, 0xc0, 0xa8, 0x00, 0x67,
	0x6f, 0x6c, 0x36, 0x18,

	0xc0, 0x65, 0x00, 0x50, 0x46, 0xf9, 0xf6, 0x99, 0x00, 0x00, 0x00, 0x00,
	0xb0, 0x02, 0xff, 0xff, 0x4e, 0x22, 0x00, 0x00,

	0x02, 0x04, 0x05, 0xb4,
	0x01,
	0x03, 0x03, 0x05,
	0x01, 0x01,
	0x08, 0x0a, 0x29, 0xf7, 0x5a, 0x07, 0x00, 0x00, 0x00, 0x00,
	0x04, 0x02,
	0x00, 0x00,
})

func TestIpPacket_SetConnAddrTcpUdp(t *testing.T) {
	buf := []byte(Ipv4TcpOptionTestCase)
	ipPacket, errMsg := NewIpv4PacketFromBuf(buf)
	udwTest.Equal(errMsg, "")
	ipPacket.SetConnAddrTcpUdp(ipPacket.GetConnAddrTcpUdp())
	udwTest.Equal(ipPacket.SerializeToBuf(), []byte(Ipv4TcpOptionTestCase))
	udwTest.Equal(ipPacket.GetSrcDstAddrPeerString(), "192.168.0.103:49253-111.108.54.24:80")
	udwTest.Equal(ipPacket.GetDstSrcAddrPeerString(), "111.108.54.24:80-192.168.0.103:49253")
	udwTest.Equal(ipPacket.GetConnAddrTcpUdp().String(), "192.168.0.103:49253-111.108.54.24:80")
	udwTest.Equal(ipPacket.GetConnAddrTcpUdp().RevertPeer().String(), "111.108.54.24:80-192.168.0.103:49253")
	udwTest.Equal(ipPacket.GetConnAddrTcpUdp().GetIpv4Array(), [12]uint8{192, 168, 0, 103, 192, 101,
		111, 108, 54, 24, 0, 80})
	udwTest.Equal(ipPacket.GetConnAddrTcpUdp().RevertPeer().GetIpv4Array(), [12]uint8{111, 108, 54, 24, 0, 80,
		192, 168, 0, 103, 192, 101})

	addr4 := ipPacket.GetConnAddrTcpUdp()
	udwTest.Benchmark(func() {
		const num = 10000
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			addr4.String()
		}
	})

	udwTest.Benchmark(func() {
		const num = 1e4
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			addr4.GetIpv4Array()
		}
	})
}
