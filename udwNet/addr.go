package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwStrings"
	"net"
	"strconv"
)

func MustGetIpFromNetAddr(addri net.Addr) (ip net.IP) {
	switch addr := addri.(type) {
	case *net.TCPAddr:
		return addr.IP
	case *net.UDPAddr:
		return addr.IP
	case *net.IPAddr:
		return addr.IP
	case *net.IPNet:
		return addr.IP
	}
	s := addri.String()
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		panic(fmt.Errorf("[MustGetIpFromAddr] %s addr.String()[%s]", err, addri.String()))
	}
	ip = net.ParseIP(host)
	if ip == nil {
		panic(fmt.Errorf("[MustGetIpFromAddr] net.ParseIP fail host:[%s]", host))
	}
	return ip
}

func GetIpStringFromtNetAddrIgnoreNotExist(addri net.Addr) (s string) {
	switch addr := addri.(type) {
	case *net.TCPAddr:
		return addr.IP.String()
	case *net.UDPAddr:
		return addr.IP.String()
	case *net.IPAddr:
		return addr.IP.String()
	case *net.IPNet:
		return addr.IP.String()
	}
	addrS := addri.String()
	host, _, err := net.SplitHostPort(addrS)
	if err != nil {
		return ""
	}
	return host
}

func MustGetPortFromNetAddr(addri net.Addr) int {
	switch addr := addri.(type) {
	case *net.TCPAddr:
		return addr.Port
	case *net.UDPAddr:
		return addr.Port
	case *net.IPAddr:
		panic("[MustGetPortFromNetAddr] *net.IPAddr do not have port")
	case *net.IPNet:
		panic("[MustGetPortFromNetAddr] *net.IPNet do not have port")
	}
	s := addri.String()
	_, port, err := net.SplitHostPort(s)
	if err != nil {
		panic(fmt.Errorf("[MustGetPortFromNetAddr] %s addr.String()[%s]", err, addri.String()))
	}
	portI := udwStrconv.MustParseInt(port)
	return portI
}

func MustGetPortFromNetAddrIgnoreNotFound(addri net.Addr) int {
	switch addr := addri.(type) {
	case *net.TCPAddr:
		return addr.Port
	case *net.UDPAddr:
		return addr.Port
	case *net.IPAddr:
		return 0
	case *net.IPNet:
		return 0
	}
	s := addri.String()
	_, port, err := net.SplitHostPort(s)
	if err != nil {
		return 0
	}
	portI := udwStrconv.MustParseInt(port)
	return portI
}

func MustSplitIpPort(addr string) (ip net.IP, port uint16) {
	ip, port, err := SplitIpPort(addr)
	if err != nil {
		panic(err)
	}
	return ip, port
}

func SplitIpPort(addr string) (ip net.IP, port uint16, err error) {
	hostS, portS, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("[MustSplitIpPort] net.SplitHostPort fail ip [%s] [%s]", hostS, err.Error())
	}
	ip = net.ParseIP(hostS)
	if ip == nil {
		return nil, 0, fmt.Errorf("[MustSplitIpPort] host is not ip [%s]", hostS)
	}
	portI, err := strconv.Atoi(portS)
	if err != nil {
		return nil, 0, fmt.Errorf("[MustSplitIpPort] port is not int [%s] [%s]", hostS, err.Error())
	}
	return ip, uint16(portI), nil
}
func GetIpStringFromAddrStringIgnoreNotFound(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}
	return host
}

func GetIpAndPortFromNetAddr(addri net.Addr) (ip net.IP, port uint16, errMsg string) {
	if addri == nil {
		return ip, port, "2r5ydz94we"
	}
	switch addr := addri.(type) {
	case *net.TCPAddr:
		return addr.IP, uint16(addr.Port), ""
	case *net.UDPAddr:
		return addr.IP, uint16(addr.Port), ""
	case *net.IPAddr:
		return ip, port, "[MustGetIpAndPortFromNetAddr] *net.IPAddr do not have port"
	case *net.IPNet:
		return ip, port, "[MustGetIpAndPortFromNetAddr] *net.IPNet do not have port"
	}
	s := addri.String()
	host, portS, err := net.SplitHostPort(s)
	if err != nil {
		return ip, port, "[MustGetIpAndPortFromNetAddr] net.SplitHostPort fail " + err.Error() + " " + addri.String()
	}
	ip = net.ParseIP(host)
	if ip == nil {
		return ip, port, "[MustGetIpAndPortFromNetAddr] net.ParseIP fail " + host
	}
	portI, err := udwStrconv.ParseInt(portS)
	if err != nil {
		return ip, port, "[MustGetIpAndPortFromNetAddr] udwStrconv.ParseInt fail " + err.Error() + " " + portS
	}
	return ip, uint16(portI), ""
}

func MustSplitIpPortForListener(addr string) (ip net.IP, port uint16) {
	hostS, portS, err := net.SplitHostPort(addr)
	if err != nil {
		panic("[MustSplitIpPortForListener] net.SplitHostPort fail ip [" + hostS + "] [" + err.Error() + "]")
	}
	if hostS != "" {
		ip = net.ParseIP(hostS)
		if ip == nil {
			panic("[MustSplitIpPortForListener] host is not ip [" + hostS + "]")
		}
	}
	portI, err := strconv.Atoi(portS)
	if err != nil {
		panic("[MustSplitIpPort] port is not int [" + portS + "] [" + err.Error() + "]")
	}
	return ip, uint16(portI)
}

func GetIpStringNoPortOrInput(addr string) string {
	return udwStrings.StringBeforeFirstSubStringOrInput(addr, ":")
}
