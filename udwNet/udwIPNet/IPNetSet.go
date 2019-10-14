package udwIPNet

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwNet"
	"net"
	"sort"
)

type IPNetSet struct {
	list  []IPNet
	listB []IPNet
}

func NewAllPassIPNet() *IPNetSet {
	return &IPNetSet{
		list: []IPNet{
			NewIPNetFromIPAndPrefix(net.IP{0, 0, 0, 0}, 0),
		},
	}
}

func NewIpSetFromIpAndPrefixList(ipPrefixList []string) *IPNetSet {
	out := &IPNetSet{}
	for _, ipPrefix := range ipPrefixList {
		out.MergeIPNet(NewIPNetFromIPAndPrefixString(ipPrefix))
	}
	return out
}

func (set *IPNetSet) MergeIPNet(ipnet IPNet) {
	newList := set.listB[:0]
	for i := range set.list {

		if set.list[i].ContainIPNet(ipnet) {
			return
		}

		if ipnet.ContainIPNet(set.list[i]) {
			continue
		}
		newList = append(newList, set.list[i])
	}
	newList = append(newList, ipnet)
	set.listB = set.list
	set.list = newList
}

func (set *IPNetSet) RemoveIPNet(ipnet IPNet) {
	newList := set.listB[:0]
	for i := range set.list {
		if set.list[i].Equal(ipnet) {
			continue
		}
		if ipnet.ContainIPNet(set.list[i]) {
			continue
		}
		if set.list[i].ContainIPNet(ipnet) {

			thisIp := net.IPv4zero.To4()
			for j := 0; j < set.list[i].prefix; j++ {
				bit := udwNet.Ipv4GetBit(ipnet.ip, uint8(j))
				thisIp = udwNet.IpSetBit(thisIp, uint8(j), bit)
			}
			for j := set.list[i].prefix; j < ipnet.prefix; j++ {
				bit := udwNet.Ipv4GetBit(ipnet.ip, uint8(j))
				thisIp = udwNet.IpSetBit(thisIp, uint8(j), bit)
				writeIp := udwNet.IpSetBit(thisIp, uint8(j), 1^bit)
				thisIpNet := NewIPNetFromIPAndPrefix(writeIp, j+1)
				newList = append(newList, thisIpNet)
			}
			continue
		}
		newList = append(newList, set.list[i])
	}
	set.listB = set.list
	set.list = newList
}

func (set *IPNetSet) ContainIP(ip net.IP) bool {
	for i := range set.list {
		if set.list[i].ContainIP(ip) {
			return true
		}
	}
	return false
}

func (set *IPNetSet) String() string {
	_buf := bytes.Buffer{}
	list := set.GetIPNetList()
	for i := range list {
		_buf.WriteString(list[i].String())
		if i != len(list)-1 {
			_buf.WriteByte('\n')
		}
	}
	return _buf.String()
}

type iPNetSlice []IPNet

func (s iPNetSlice) Len() int {
	return len(s)
}
func (s iPNetSlice) Less(i int, j int) bool {
	if s[i].ip.Equal(s[j].ip) {
		return s[i].prefix < s[j].prefix
	}
	return udwNet.IpLess(s[i].ip, s[j].ip)
}
func (s iPNetSlice) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (set *IPNetSet) GetIPNetList() []IPNet {
	if set == nil {
		return nil
	}
	if len(set.list) <= 1 {
		return set.list
	}
	sort.Sort(iPNetSlice(set.list))
	return set.list
}

func (set *IPNetSet) SimpleIpNet() {
	if set == nil || len(set.list) <= 1 {
		return
	}
	sort.Sort(iPNetSlice(set.list))
	for {
		thisDeleteSize := 0
		newList := set.listB[:0]
		newList = append(newList, set.list[0])
		for i := 1; i < len(set.list); i++ {
			lastNew := newList[len(newList)-1]
			thisIpNet := set.list[i]
			if lastNew.IsRootIpNet() {
				break
			}
			if thisIpNet.IsRootIpNet() {
				newList = append(newList, thisIpNet)
				break
			}
			if lastNew.prefix == thisIpNet.prefix {
				parentIpNet := lastNew.GetParentIpNet()
				if parentIpNet.ContainIPNet(lastNew) && parentIpNet.ContainIPNet(thisIpNet) {
					newList[len(newList)-1] = parentIpNet
					thisDeleteSize++
					continue
				}
			}
			newList = append(newList, thisIpNet)
		}
		if thisDeleteSize > 0 {
			set.listB = set.list
			set.list = newList
		} else {
			break
		}
	}
	return
}
