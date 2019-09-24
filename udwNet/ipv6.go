package udwNet

import "net"

func IsIpv6Only() bool {
	conn, err := net.Dial("udp", "104.199.139.77:81")
	if err == nil {
		conn.Close()
		return false
	}
	if IsNetworkIsUnreachable(err) {
		conn, err := net.Dial("udp", "[64:ff9b::68c7:8b4d]:81")
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

func TransferIpv4StringToSupportIpv6Only(ipS string) string {
	ip := net.ParseIP(ipS)
	if ip.To4() == nil {
		return ipS
	}

	if IsIpv6Only() {
		return ipS + ".ipv4.tubnetwork.com"
	}
	return ipS
}
