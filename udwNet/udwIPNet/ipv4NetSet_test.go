package udwIPNet

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
	"testing"
)

func TestIpv4Net(ot *testing.T) {
	ipnet := MustParseIpv4Net("127.0.0.1/8")
	udwTest.Ok(ipnet.GetIp().Equal(net.ParseIP("127.0.0.1")))
	udwTest.Equal(ipnet.GetPrefix(), 8)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("127.1.0.1")))
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("128.1.0.1")) == false)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("1.1.0.1")) == false)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("126.0.0.1")) == false)
	udwTest.Ok(ipnet.ContainIP(net.ParseIP("126.1.0.1")) == false)

	udwTest.Ok(MustParseIpv4Net("0.0.0.0/0").ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(MustParseIpv4Net("0.0.0.0/0").ContainIPNet(MustParseIpv4Net("0.0.0.0/0")))
	udwTest.Ok(MustParseIpv4Net("0.0.0.0/0").ContainIPNet(MustParseIpv4Net("127.0.0.1/32")))

	udwTest.Ok(MustParseIpv4Net("0.0.0.0/1").ContainIP(net.ParseIP("127.0.0.1")))
	udwTest.Ok(MustParseIpv4Net("0.0.0.0/1").ContainIP(net.ParseIP("129.0.0.1")) == false)

	set := NewAllPassIpv4Net()
	set.RemoveIPNet(MustParseIpv4Net("127.0.0.1/32"))
	udwTest.Equal(set.String(), `0.0.0.0/2
64.0.0.0/3
96.0.0.0/4
112.0.0.0/5
120.0.0.0/6
124.0.0.0/7
126.0.0.0/8
127.0.0.0/32
127.0.0.2/31
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
	ipv4MustNotHasOverlap(set)
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.1")) == false)
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.2")))

	set.RemoveIPNet(MustParseIpv4Net("127.0.0.2/32"))
	udwTest.Ok(set.ContainIP(net.ParseIP("127.0.0.2")) == false)
	udwTest.Ok(set.ContainIP(net.ParseIP("0.0.0.0")))
	udwTest.Ok(set.ContainIP(net.ParseIP("0.0.0.1")))

	ipv4MustNotHasOverlap(set)
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
}

func TestIpv4Net2(ot *testing.T) {
	set := NewAllPassIpv4Net()
	udwTest.Equal(set.String(), `0.0.0.0/0`)
	udwTest.Equal(len(set.GetIpv4NetList()), 1)
	set.RemoveIpString("116.93.97.59")
	udwTest.Equal(len(set.GetIpv4NetList()), 32)
	udwTest.Equal(set.String(), `0.0.0.0/2
64.0.0.0/3
96.0.0.0/4
112.0.0.0/6
116.0.0.0/10
116.64.0.0/12
116.80.0.0/13
116.88.0.0/14
116.92.0.0/16
116.93.0.0/18
116.93.64.0/19
116.93.96.0/24
116.93.97.0/27
116.93.97.32/28
116.93.97.48/29
116.93.97.56/31
116.93.97.58/32
116.93.97.60/30
116.93.97.64/26
116.93.97.128/25
116.93.98.0/23
116.93.100.0/22
116.93.104.0/21
116.93.112.0/20
116.93.128.0/17
116.94.0.0/15
116.96.0.0/11
116.128.0.0/9
117.0.0.0/8
118.0.0.0/7
120.0.0.0/5
128.0.0.0/1`)
}

func ipv4MustNotHasOverlap(set *Ipv4NetSet) {
	list := set.List
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
