package udwIPNet

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBitwise"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwSort"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"net"
	"strconv"
	"strings"
)

type Ipv4NetSet struct {
	list  []Ipv4Net
	listB []Ipv4Net
}

func NewAllPassIpv4Net() *Ipv4NetSet {
	set := &Ipv4NetSet{
		list: []Ipv4Net{0},
	}

	return set
}

func (set *Ipv4NetSet) RemoveIPNet(ipnet Ipv4Net) {
	newList := set.listB[:0]
	for i := range set.list {
		if set.list[i].Equal(ipnet) {

			continue
		}
		if ipnet.ContainIPNet(set.list[i]) {

			continue
		}
		if set.list[i].ContainIPNet(ipnet) {

			thisIp := uint32(0)
			ipNetIpUint32 := ipnet.GetUint32Ip()
			for j := 0; j < set.list[i].GetPrefix(); j++ {
				bit := udwBitwise.Uint32GetBit(ipNetIpUint32, 31-j)
				thisIp = udwBitwise.Uint32SetBit(thisIp, 31-j, bit)
			}
			for j := set.list[i].GetPrefix(); j < ipnet.GetPrefix(); j++ {
				bit := udwBitwise.Uint32GetBit(ipNetIpUint32, 31-j)

				thisIp = udwBitwise.Uint32SetBit(thisIp, 31-j, bit)
				writeIp := udwBitwise.Uint32SetBit(thisIp, 31-j, 1^bit)

				thisIpNet := NewIpv4NetFromUint32IpAndPrefix(writeIp, j+1)
				newList = append(newList, thisIpNet)
			}

			continue
		}
		newList = append(newList, set.list[i])
	}
	set.listB = set.list
	set.list = newList
}

func (set *Ipv4NetSet) ContainIP(ip net.IP) bool {
	for i := range set.list {
		if set.list[i].ContainIP(ip) {
			return true
		}
	}
	return false
}

func (set *Ipv4NetSet) RemoveIp(ip net.IP) {
	ipv4 := ip.To4()
	if ipv4 == nil {
		panic("ghvms9vm43")
	}
	set.RemoveIPNet(NewIpv4NetFromUint32IpAndPrefix(udwNet.Ipv4ToUint32OrZero(ipv4), 32))
}
func (set *Ipv4NetSet) RemoveIpString(ipString string) {

	ipv4Uint32, ok := udwNet.ParseIpv4ToUint32(ipString)
	if ok == false {
		panic("xt76auvktq ipString:[" + ipString + "]")
	}
	set.RemoveIPNet(NewIpv4NetFromUint32IpAndPrefix(ipv4Uint32, 32))
}

func (set *Ipv4NetSet) String() string {
	_buf := bytes.Buffer{}
	list := set.GetIpv4NetList()
	for i := range list {
		_buf.WriteString(list[i].String())
		if i != len(list)-1 {
			_buf.WriteByte('\n')
		}
	}
	return _buf.String()
}
func (set *Ipv4NetSet) GetIpv4NetList() []Ipv4Net {
	if set == nil {
		return nil
	}
	set.Sort()
	return set.list
}
func (set *Ipv4NetSet) Sort() {
	udwSort.InterfaceCallbackSortWithIndexLess(set.list, func(a int, b int) bool {
		aIp := set.list[a].GetUint32Ip()
		bIp := set.list[b].GetUint32Ip()
		if aIp != bIp {
			return aIp < bIp
		}
		return set.list[a].GetPrefix() < set.list[b].GetPrefix()
	})
}

type Ipv4Net uint64

func NewIpv4NetFromUint32IpAndPrefix(ip uint32, prefix int) Ipv4Net {
	return Ipv4Net(uint64(ip) | uint64((uint64(prefix)&0xff)<<32))
}

func MustParseIpv4Net(s string) Ipv4Net {
	partList := strings.Split(s, "/")
	if len(partList) != 2 {
		panic("n2864gvqpa len(partList)[" + strconv.Itoa(len(partList)) + "]!=2 [" + s + "]")
	}
	ip := net.ParseIP(partList[0])
	if ip == nil {
		panic("mzxjhn9rmz ip can not parse [" + s + "]")
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		panic("bjk6rv7fdp")
	}
	prefix := udwStrconv.MustParseInt(partList[1])
	return NewIpv4NetFromUint32IpAndPrefix(udwNet.MustIpv4ToUint32(ipv4), prefix)
}
func (thisIpNet Ipv4Net) String() string {
	return udwNet.Uint32ToIpv4(thisIpNet.GetUint32Ip()).String() + "/" + strconv.Itoa(int(thisIpNet.GetPrefix()))
}
func (thisIpNet Ipv4Net) GetPrefix() int {
	return int(uint8(thisIpNet >> 32 & 0xff))
}
func (thisIpNet Ipv4Net) GetUint32Ip() uint32 {
	return uint32(thisIpNet & 0xffffffff)
}
func (thisIpNet Ipv4Net) GetIp() net.IP {
	return udwNet.Uint32ToIpv4(thisIpNet.GetUint32Ip())
}
func (thisIpNet Ipv4Net) Equal(ipnet Ipv4Net) bool {
	return ipnet == thisIpNet
}
func (thisIpnet Ipv4Net) ContainIP(ip net.IP) bool {
	ret := thisIpnet.ContainUint32Ip(udwNet.MustIpv4ToUint32(ip))
	return ret
}
func (thisIpnet Ipv4Net) ContainIPNet(ipnet Ipv4Net) bool {
	ret := ipnet.GetPrefix() >= thisIpnet.GetPrefix() && thisIpnet.ContainUint32Ip(ipnet.GetUint32Ip())
	return ret
}
func (thisIpnet Ipv4Net) ContainUint32Ip(ip uint32) bool {
	thisUint32Ip := thisIpnet.GetUint32Ip()

	for i := 0; i < thisIpnet.GetPrefix(); i++ {
		thisBit1 := udwBitwise.Uint32GetBit(thisUint32Ip, 31-i)
		thisBit2 := udwBitwise.Uint32GetBit(ip, 31-i)

		if thisBit1 != thisBit2 {
			return false
		}
	}
	return true
}

func (thisIpNet Ipv4Net) ToGoIPNet() *net.IPNet {
	_, ipnet, err := net.ParseCIDR(thisIpNet.String())
	if err != nil {
		panic(err)
	}
	return ipnet
}
