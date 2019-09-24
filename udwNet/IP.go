package udwNet

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwBytes"
	"net"
)

func IpSetBit1(ip net.IP, index uint8) net.IP {
	ip = IpClone(ip)
	ip[index>>3] = ip[index>>3] | (1 << (7 - (index % 8)))
	return ip
}

func IpSetBit0(ip net.IP, index uint8) net.IP {
	ip = IpClone(ip)
	ip[index>>3] = ip[index>>3] & ^(1 << (7 - (index % 8)))
	return ip
}

func IpSetBit(ip net.IP, index uint8, value uint8) net.IP {
	if value == 0 {
		return IpSetBit0(ip, index)
	} else {
		return IpSetBit1(ip, index)
	}
}

func IpGetBit(ip net.IP, index uint8) uint8 {
	return uint8(1) & (uint8(ip[index>>3]) >> (7 - (index % 8)))
}

func Ipv4GetBit(ip net.IP, index uint8) uint8 {
	ip = ip.To4()
	return uint8(1) & (uint8(ip[index>>3]) >> (7 - (index % 8)))
}

func IpClone(a net.IP) net.IP {
	b := make(net.IP, len(a))
	copy(b, a)
	return b
}

func IpLess(a net.IP, b net.IP) bool {
	return string(a) < string(b)
}

func MustIpv4ToWindowsDword(a net.IP) uint32 {
	return MustIpv4ToUint32(a)

}

func MustIpv4ToUint32(a net.IP) uint32 {
	ipv4 := a.To4()
	if ipv4 == nil {
		panic("hv5pue2hvt " + a.String())
	}
	return binary.BigEndian.Uint32([]byte(ipv4))
}

func Ipv4ToUint32OrZero(a net.IP) uint32 {
	ipv4 := a.To4()
	if ipv4 == nil {
		return 0
	}
	return binary.BigEndian.Uint32([]byte(ipv4))
}

func Uint32ToIpv4(u uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32([]byte(ip), u)
	return ip
}

func IsStringValidIp(sIp string) bool {
	ip := net.ParseIP(sIp)
	return ip != nil
}

func Ipv4AddAndCopyWithBuffer(ip net.IP, toAdd uint32, bufW *udwBytes.BufWriter) net.IP {
	ipv4 := ip.To4()
	if ipv4 == nil {
		panic("[ipAdd] ip is not ipv4 addr")
	}
	ipInt := binary.BigEndian.Uint32(ipv4)
	ipInt += toAdd
	bufW = udwBytes.ResetBufWriter(bufW)
	bufW.WriteBigEndUint32(ipInt)
	return net.IP(bufW.GetBytes())
}

func MustIpv4Add(ip net.IP, toAdd uint32) net.IP {
	return Ipv4AddAndCopyWithBuffer(ip, toAdd, nil)
}
