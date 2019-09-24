package udwNet

import (
	"encoding/binary"
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"net"
	"strings"
	"sync"
)

func MustParseIPNet(s string) *net.IPNet {
	return MustParseCIDR(s)
}

func IPNetGetGenMaskString(ipNet *net.IPNet) (s string) {
	return net.IP(ipNet.Mask).String()
}

func MustParseCIDR(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ipnet
}

var privateAddrMaskListOnce sync.Once
var privateAddrMaskList []*net.IPNet

func getPrivateAddrMaskList() []*net.IPNet {
	privateAddrMaskListOnce.Do(func() {
		privateAddrMaskList = []*net.IPNet{
			MustParseCIDR("0.0.0.0/8"),
			MustParseCIDR("10.0.0.0/8"),
			MustParseCIDR("127.0.0.0/8"),
			MustParseCIDR("172.16.0.0/12"),
			MustParseCIDR("169.254.0.0/16"),
			MustParseCIDR("192.168.0.0/16"),
			MustParseCIDR("198.18.0.0/15"),
			MustParseCIDR("100.64.0.0/10"),
			MustParseCIDR("::1/128"),
			MustParseCIDR("fc00::/7"),
			MustParseCIDR("fe80::/10"),
		}
	})
	return privateAddrMaskList
}

func IsPrivateNetwork(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return isPrivateNetworkFromIpNetList(ip)
	}
	v4 := binary.BigEndian.Uint32(ip4[:])

	return v4 == 0xffffffff ||
		v4 <= 0xffffff ||
		((v4 >= 0xa000000) && (v4 <= 0xaffffff)) ||
		((v4 >= 0x64400000) && (v4 <= 0x647fffff)) ||
		((v4 >= 0x7f000000) && (v4 <= 0x7fffffff)) ||
		((v4 >= 0xa9fe0000) && (v4 <= 0xa9feffff)) ||
		((v4 >= 0xac100000) && (v4 <= 0xac1fffff)) ||
		((v4 >= 0xc0a80000) && (v4 <= 0xc0a8ffff)) ||
		((v4 >= 0xc6120000) && (v4 <= 0xc613ffff)) ||
		((v4 >= 0xe0000000) && (v4 <= 0xefffffff))
}

func isPrivateNetworkFromIpNetList(ip net.IP) bool {
	for _, ipnet := range getPrivateAddrMaskList() {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func IsLocalNetwork(ip net.IP) bool {
	return MustParseCIDR("127.0.0.0/8").Contains(ip) || MustParseCIDR("::1/128").Contains(ip)
}

var linkLocalAddrMaskListOnce sync.Once
var linkLocalAddrMaskList []*net.IPNet

func IsInLinkLocalNetwork(ip net.IP) bool {
	linkLocalAddrMaskListOnce.Do(func() {
		linkLocalAddrMaskList = []*net.IPNet{
			MustParseCIDR("127.0.0.0/8"),
			MustParseCIDR("::1/128"),
		}
	})
	for _, ipnet := range linkLocalAddrMaskList {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

var loopBackAddrMaskListOnce sync.Once
var loopBackAddrMaskList []*net.IPNet

func IsInLoopBackNetwork(ip net.IP) bool {
	loopBackAddrMaskListOnce.Do(func() {
		loopBackAddrMaskList = []*net.IPNet{
			MustParseCIDR("169.254.0.0/16"),
			MustParseCIDR("fe80::/10"),
		}
	})
	for _, ipnet := range loopBackAddrMaskList {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func ParseGenmask(genMask string) (mask net.IPMask, err error) {
	ip := net.ParseIP(genMask)
	if ip == nil {
		return nil, fmt.Errorf("[ParseGenmask] ip==nil genmask %s", genMask)
	}
	ip = ip.To4()
	if ip == nil {
		return nil, fmt.Errorf("[ParseGenmask] ip.To4()==nil genmask %s", genMask)
	}
	return net.IPv4Mask(ip[0], ip[1], ip[2], ip[3]), nil
}

func MustParseIpStringToInt64(ipStr string) int64 {
	bits := strings.Split(ipStr, ".")
	b0 := udwStrconv.AtoIDefault0(bits[0])
	b1 := udwStrconv.AtoIDefault0(bits[1])
	b2 := udwStrconv.AtoIDefault0(bits[2])
	b3 := udwStrconv.AtoIDefault0(bits[3])
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

func IPNetGetOnes(ipnet net.IPNet) int {
	ones, _ := ipnet.Mask.Size()
	return ones
}
func IPMaskGetOnes(mask net.IPMask) int {
	ones, _ := mask.Size()
	return ones
}
func IPNetIsZero(ipnet net.IPNet) bool {
	_, bits := ipnet.Mask.Size()
	return bits == 0
}
func IPNetGetMinIp(ipnet *net.IPNet) net.IP {

	maskIp := ipnet.Mask
	minIp := IpClone(ipnet.IP)
	for i := 0; i < len(minIp); i++ {
		if maskIp[i] != 0xff {
			minIp[i] = minIp[i] & maskIp[i]
		}
	}
	return minIp
}
func IPNetGetMaxIp(ipnet *net.IPNet) net.IP {

	maskIp := ipnet.Mask
	minIp := IPNetGetMinIp(ipnet)
	for i := 0; i < len(minIp); i++ {
		notMaskB := ^maskIp[i]
		if notMaskB != 0 {
			minIp[i] += notMaskB
		}
	}
	return minIp
}
