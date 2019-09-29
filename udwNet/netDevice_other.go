// +build !windows,!js

package udwNet

import "net"

func mustGetNetDeviceList() []*NetDevice {
	goNetDevList, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	out := make([]*NetDevice, len(goNetDevList))
	for i, goNetDev := range goNetDevList {
		addrList, err := goNetDev.Addrs()
		if err != nil {
			panic(err)
		}
		out[i] = &NetDevice{
			index:     goNetDev.Index,
			name:      goNetDev.Name,
			firstIpv4: mustGetFirstIpv4IpFromGoNetDev(goNetDev, addrList),
		}
		for _, addrI := range addrList {
			ip := MustGetIpFromNetAddr(addrI)
			out[i].ipList = append(out[i].ipList, ip)
		}
		out[i].isPointToPoint = (goNetDev.Flags&net.FlagPointToPoint == net.FlagPointToPoint)
		out[i].isUp = (goNetDev.Flags&net.FlagUp == net.FlagUp)
		out[i].isBroadcast = (goNetDev.Flags&net.FlagBroadcast == net.FlagBroadcast)
		out[i].isMulticast = (goNetDev.Flags&net.FlagMulticast == net.FlagMulticast)
		out[i].isLoopback = (goNetDev.Flags&net.FlagLoopback == net.FlagLoopback)
	}
	return out
}

func mustGetFirstIpv4IpFromGoNetDev(goNetDev net.Interface, addrList []net.Addr) net.IP {

	for _, addrI := range addrList {
		ip := MustGetIpFromNetAddr(addrI)
		ip4 := ip.To4()
		if ip4 != nil {
			return ip4
		}
	}
	return nil
}
