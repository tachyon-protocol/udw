package udwNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestIPNet(ot *testing.T) {
	udwTest.Equal(IPNetIsZero(net.IPNet{}), true)
}

func TestIsLocalNetwork(t *testing.T) {
	udwTest.Ok(IsLocalNetwork(net.ParseIP("127.23.0.1")))
}

func TestIPNetGetGenMaskString(ot *testing.T) {
	udwTest.Equal(IPNetGetGenMaskString(MustParseIPNet("0.0.0.0/8")), "255.0.0.0")
}

func TestIsPrivateNetwork(t *testing.T) {

	udwTest.Equal(IsPrivateNetwork(net.ParseIP("0.0.0.0")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("0.0.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("0.255.255.255")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("1.0.0.0")), false)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("10.0.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("127.0.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("172.16.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("172.33.0.1")), false)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("169.254.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("192.168.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("198.18.0.1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("100.67.164.96")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("::1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("::2")), false)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("fc00::1")), true)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("fe00::1")), false)
	udwTest.Equal(IsPrivateNetwork(net.ParseIP("fe80::1")), true)

	const testNum = 1e7
	ip := net.ParseIP("1.0.0.0")
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(testNum)
		for i := 0; i < testNum; i++ {
			IsPrivateNetwork(ip)
		}
	})
}
