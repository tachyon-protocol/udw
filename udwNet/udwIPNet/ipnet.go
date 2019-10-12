package udwIPNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"net"
	"strconv"
	"strings"
)

type IPNet struct {
	ip     net.IP
	prefix int
}

func (net IPNet) String() string {
	if net.prefix >= 64 {
		ip := net.ip.To4()
		if ip != nil {

			return ip.String() + "/" + strconv.Itoa(net.prefix-96)
		}
	}
	return net.ip.String() + "/" + strconv.Itoa(net.prefix)
}

func (net IPNet) ContainIP(ip net.IP) bool {
	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}
	for i := uint8(0); i < uint8(net.prefix); i++ {
		thisBit1 := udwNet.IpGetBit(ip, i)
		thisBit2 := udwNet.IpGetBit(net.ip, i)
		if thisBit1 != thisBit2 {
			return false
		}
	}
	return true
}

func (thisIpnet IPNet) ContainIPNet(ipnet IPNet) bool {
	return ipnet.prefix >= thisIpnet.prefix && thisIpnet.ContainIP(ipnet.ip)
}
func (ipnet IPNet) GetIP() net.IP {
	return ipnet.ip
}

func (ipnet IPNet) GetIPString() string {
	return ipnet.ip.String()
}

func (ipnet IPNet) GetMaskIpString() string {
	return net.IP(net.CIDRMask(ipnet.prefix, len(ipnet.ip)*8)).String()
}
func (ipnet IPNet) GetPrefix() int {
	return ipnet.prefix
}
func (ipnet IPNet) GetMaskPrefixInIpv6() int {
	maskSize := ipnet.prefix
	if ipnet.GetIpByteLen() == 4 && maskSize <= 32 {
		maskSize += 96
	}
	return maskSize
}
func (ipnet IPNet) GetIpByteLen() int {
	if ipnet.ip.To4() != nil {
		return 4
	}
	return 16
}
func (thisIpNet IPNet) Equal(ipnet IPNet) bool {
	return thisIpNet.prefix == ipnet.prefix && thisIpNet.ip.Equal(ipnet.ip)
}
func (thisIpNet IPNet) ToGoIPNet() *net.IPNet {
	_, ipnet, err := net.ParseCIDR(thisIpNet.String())
	if err != nil {
		panic(err)
	}
	return ipnet
}
func (thisIpNet IPNet) IsRootIpNet() bool {
	return thisIpNet.prefix == 0
}
func (thisIpNet IPNet) GetParentIpNet() IPNet {
	thisIpNet.prefix -= 1
	thisIpNet.ip = udwNet.IpClone(thisIpNet.ip)
	ipSize := thisIpNet.GetIpByteLen()
	for i := uint8(thisIpNet.prefix); i < uint8(ipSize); i++ {
		thisBit1 := udwNet.IpGetBit(thisIpNet.ip, i)
		if thisBit1 != 0 {
			thisIpNet.ip = udwNet.IpSetBit0(thisIpNet.ip, i)
		}
	}
	return thisIpNet
}

func ParseIPNetDefault32(s string) (IPNet, error) {
	if !strings.Contains(s, "/") {
		s += "/32"
	}
	var ipNet IPNet
	err := udwErr.PanicToError(func() {
		ipNet = MustParseIPNet(s)
	})
	return ipNet, err
}

func MustParseIPNet(s string) IPNet {
	partList := strings.Split(s, "/")
	if len(partList) != 2 {
		panic(fmt.Errorf("[MustParseIPNet] len(partList)[%d]!=2 [%s]", len(partList), s))
	}
	ip := net.ParseIP(partList[0])
	if ip == nil {
		panic(fmt.Errorf("[MustParseIPNet] ip can not parse [%s]", s))
	}
	ipv4 := ip.To4()
	if ipv4 != nil {
		ip = ipv4
	}
	prefix := udwStrconv.MustParseInt(partList[1])
	return IPNet{
		ip:     ip,
		prefix: prefix,
	}
}

func NewIPNetFromIPAndPrefix(ip net.IP, prefix int) IPNet {
	return IPNet{
		ip:     ip,
		prefix: prefix,
	}
}

func NewIPNetFromIPAndPrefixString(s string) IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	one, _ := ipnet.Mask.Size()
	return IPNet{
		ip:     ipnet.IP,
		prefix: one,
	}
}

func GetIpv4RootIpNet() IPNet {
	return NewIPNetFromIPAndPrefix(net.IP{0, 0, 0, 0}, 0)
}

func GetIpv6RootIpNet() IPNet {
	return NewIPNetFromIPAndPrefix(net.IPv6zero, 0)
}
