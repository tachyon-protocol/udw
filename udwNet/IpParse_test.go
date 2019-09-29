package udwNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

var parseIPTests = []struct {
	in  string
	out net.IP
}{
	{"127.0.1.2", net.IP{127, 0, 1, 2}},
	{"127.0.0.1", net.IP{127, 0, 0, 1}},
	{"127.001.002.003", net.IP{127, 1, 2, 3}},
	{"::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"::ffff:127.001.002.003", net.IPv4(127, 1, 2, 3)},
	{"::ffff:7f01:0203", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0:0000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0:000000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},

	{"2001:4860:0:2001::68", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
	{"2001:4860:0000:2001:0000:0000:0000:0068", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},

	{"-0.0.0.0", nil},
	{"0.-1.0.0", nil},
	{"0.0.-2.0", nil},
	{"0.0.0.-3", nil},
	{"127.0.0.256", nil},
	{"abc", nil},
	{"123:", nil},
	{"fe80::1%lo0", nil},
	{"fe80::1%911", nil},
	{"", nil},
	{"a1:a2:a3:a4::b1:b2:b3:b4", nil},
}

func TestParseIpAppendToBuf(t *testing.T) {
	for _, tt := range parseIPTests {
		out := ParseIpToBuf(tt.in, nil)
		udwTest.Equal(out, tt.out, tt.in)
	}
	buf := make([]byte, 16)
	const testNum = 1e3
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(testNum * len(parseIPTests))
		for i := 0; i < testNum; i++ {
			for _, tt := range parseIPTests {
				ParseIpToBuf(tt.in, buf)
			}
		}
	})

}
