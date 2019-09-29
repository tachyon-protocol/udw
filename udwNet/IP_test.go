package udwNet

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestIp(ot *testing.T) {
	udwTest.Equal(IpGetBit(net.ParseIP("128.0.0.0").To4(), 0), uint8(1))
	udwTest.Equal(IpGetBit(net.ParseIP("128.0.0.0").To4(), 1), uint8(0))
	udwTest.Equal(IpGetBit(net.ParseIP("128.0.0.0").To4(), 31), uint8(0))

	udwTest.Ok(IpSetBit(net.ParseIP("128.0.0.0").To4(), 0, 0).Equal(net.ParseIP("0.0.0.0").To4()))
	udwTest.Ok(IpSetBit(net.ParseIP("128.0.0.0").To4(), 0, 1).Equal(net.ParseIP("128.0.0.0").To4()))

	udwTest.Ok(IpSetBit(net.ParseIP("128.0.0.0").To4(), 1, 0).Equal(net.ParseIP("128.0.0.0").To4()))
	udwTest.Ok(IpSetBit(net.ParseIP("128.0.0.0").To4(), 1, 1).Equal(net.ParseIP("192.0.0.0").To4()))

}

func TestMustIpv4ToWindowsDword(t *testing.T) {
	udwTest.Equal(MustIpv4ToWindowsDword(net.ParseIP("172.21.0.1")), uint32(0xac150001))
}

func TestIpv4AddAndCopyWithBuffer(t *testing.T) {
	bufW := udwBytes.BufWriter{}
	udwTest.Ok(Ipv4AddAndCopyWithBuffer(net.ParseIP("128.0.0.0"), 1, nil).Equal(net.ParseIP("128.0.0.1")))
	udwTest.Ok(Ipv4AddAndCopyWithBuffer(net.ParseIP("128.0.0.0"), 65536, &bufW).Equal(net.ParseIP("128.1.0.0")))
	udwTest.Ok(Ipv4AddAndCopyWithBuffer(net.ParseIP("172.21.3.65"), 65536, &bufW).Equal(net.ParseIP("172.22.3.65")))

}
