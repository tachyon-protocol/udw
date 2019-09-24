package udwNet

import "net"

func MustLookupDomainInAddrString(addr string) string {
	addrTcp, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	return addrTcp.String()
}
