// +build !js

package udwNet

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwMap"
	"net"
	"strconv"
)

type NetDevice struct {
	index     int
	name      string
	ipList    []net.IP
	ipNetList []net.IPNet
	firstIpv4 net.IP

	isPointToPoint bool
	isUp           bool
	isBroadcast    bool
	isLoopback     bool
	isMulticast    bool
}

func (i *NetDevice) GetIndex() int {
	return i.index
}

func (i *NetDevice) GetName() string {
	return i.name
}

func (i *NetDevice) GetFirstIpv4IP() net.IP {
	return i.firstIpv4
}

func (i *NetDevice) GetIpList() []net.IP {
	return i.ipList
}
func (i *NetDevice) HasIpv4Addr() bool {
	return i.firstIpv4 != nil
}
func (i *NetDevice) GetNonLinkLocalIpList() []net.IP {
	if len(i.ipList) == 0 {
		return nil
	}
	out := make([]net.IP, 0, len(i.ipList))
	for _, ip := range i.ipList {
		if !IsInLinkLocalNetwork(ip) {
			out = append(out, ip)
		}
	}
	return out
}

func (i *NetDevice) IsPointToPoint() bool {
	return i.isPointToPoint
}

func (i *NetDevice) IsUp() bool {
	return i.isUp
}

func (i *NetDevice) CanConnect() bool {
	return i.isUp && len(i.GetNonLinkLocalIpList()) > 0
}

func (i *NetDevice) String() string {
	s := "udwNet.Interface "
	s += "index:" + strconv.Itoa(i.index) + " "
	s += "name:" + i.name + " "
	if len(i.ipList) > 0 {
		s += "ipList:"
		for _, ip := range i.ipList {
			s += ip.String() + ","
		}
		s += " "
	}
	if i.isPointToPoint {
		s += "PointToPoint "
	}
	if i.isUp {
		s += "Up "
	}
	if i.isLoopback {
		s += "Loopback "
	}
	if i.isBroadcast {
		s += "Broadcast "
	}
	if i.isMulticast {
		s += "Multicast "
	}
	return s
}

func MustGetNetDeviceList() []*NetDevice {
	return mustGetNetDeviceList()
}

func GetNetDeviceList() (devlist []*NetDevice, err error) {
	err = udwErr.PanicToError(func() {
		devlist = MustGetNetDeviceList()
	})
	return devlist, err
}

func GetNetDeviceByIndex(index int) (dev *NetDevice, err error) {
	devlist, err := GetNetDeviceList()
	if err != nil {
		return nil, err
	}
	for _, dev := range devlist {
		if dev.GetIndex() == index {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("[GetNetDeviceByIndex] can not found net device with index %d", index)
}

func GetNetDeviceByName(name string) (dev *NetDevice, err error) {
	devlist, err := GetNetDeviceList()
	if err != nil {
		return nil, err
	}
	for _, dev := range devlist {
		if dev.GetName() == name {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("[GetNetDeviceByIndex] can not found net device with name %s", name)
}

func MustGetIpList() (ipList []net.IP) {
	netDevList := MustGetNetDeviceList()
	for _, dev := range netDevList {
		for _, ip := range dev.GetIpList() {
			if IsInLinkLocalNetwork(ip) || IsInLoopBackNetwork(ip) {
				continue
			}
			ipList = append(ipList, ip)
		}
	}
	return ipList

}

func MustGetIpListToString() string {
	_buf := bytes.Buffer{}
	for _, ip := range MustGetIpList() {
		_buf.WriteString(ip.String())
		_buf.WriteByte(',')
	}
	return _buf.String()
}

func MustGetCurrentIpWithPortList(port uint16) (sList []string) {
	netDevList := MustGetNetDeviceList()
	sPort := strconv.Itoa(int(port))
	ipSet := map[string]struct{}{}
	for _, dev := range netDevList {

		if dev.IsPointToPoint() {
			continue
		}
		if len(dev.ipList) > 0 {

			for _, ip := range dev.ipList {
				ipv4 := ip.To4()
				if ipv4 != nil {
					ipSet[ipv4.String()] = struct{}{}
				}
			}
		}

	}
	ipList := udwMap.SetStringToStringListAes(ipSet)
	for _, ipS := range ipList {
		sList = append(sList, net.JoinHostPort(ipS, sPort))
	}
	return sList

}
