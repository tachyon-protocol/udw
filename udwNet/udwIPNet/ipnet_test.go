package udwIPNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestMustParseIPNet(ot *testing.T) {
	ipnet := MustParseIPNet("127.0.0.1/8")
	udwTest.Ok(ipnet.ip.Equal(net.ParseIP("127.0.0.1")))
	udwTest.Equal(ipnet.prefix, 8)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("127.1.0.1")))
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("128.1.0.1")) == false)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("1.1.0.1")) == false)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("126.1.0.1")) == false)

	udwTest.Ok(MustParseIPNet("0.0.0.0/0").ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(MustParseIPNet("0.0.0.0/0").ContainIPNet(MustParseIPNet("0.0.0.0/0")))
	udwTest.Ok(MustParseIPNet("0.0.0.0/0").ContainIPNet(MustParseIPNet("127.0.0.1/32")))

	udwTest.Ok(MustParseIPNet("0.0.0.0/1").ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(MustParseIPNet("0.0.0.0/1").ContainIP(net.ParseIP("129.0.0.1")) == false)

	set := NewAllPassIPNet()
	set.RemoveIPNet(MustParseIPNet("127.0.0.1/32"))

	mustNotHasOverlap(set)
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.1")) == false)
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.2")))

	set.RemoveIPNet(MustParseIPNet("127.0.0.2/32"))
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.2")) == false)
	udwTest.Ok(set.ContainIP(net.ParseIP("0.0.0.0")))
	udwTest.Ok(set.ContainIP(net.ParseIP("0.0.0.1")))

	mustNotHasOverlap(set)
	udwTest.Equal(set.String(), `0.0.0.0/2
64.0.0.0/3
96.0.0.0/4
112.0.0.0/5
120.0.0.0/6
124.0.0.0/7
126.0.0.0/8
127.0.0.0/32
127.0.0.3/32
127.0.0.4/30
127.0.0.8/29
127.0.0.16/28
127.0.0.32/27
127.0.0.64/26
127.0.0.128/25
127.0.1.0/24
127.0.2.0/23
127.0.4.0/22
127.0.8.0/21
127.0.16.0/20
127.0.32.0/19
127.0.64.0/18
127.0.128.0/17
127.1.0.0/16
127.2.0.0/15
127.4.0.0/14
127.8.0.0/13
127.16.0.0/12
127.32.0.0/11
127.64.0.0/10
127.128.0.0/9
128.0.0.0/1`)
	set.MergeIPNet(MustParseIPNet("127.0.0.1/32"))
	set.MergeIPNet(MustParseIPNet("127.0.0.2/32"))
	mustNotHasOverlap(set)
	set.SimpleIpNet()

	udwTest.Equal(set.String(), `0.0.0.0/0`)
}

func mustNotHasOverlap(set *IPNetSet) {
	list := set.GetIPNetList()
	for i := range list {
		for j := range list {
			if i == j {
				continue
			}
			if list[i].ContainIPNet(list[j]) {
				panic("[mustNotHasOverlap] " + list[i].String() + " " + list[j].String())
			}
		}
	}
}
