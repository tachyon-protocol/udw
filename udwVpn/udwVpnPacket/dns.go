package udwVpnPacket

import (
	"github.com/tachyon-protocol/udw/udwNet/udwDns/udwDnsPacket"
)

func (packet *IpPacket) IsDnsRequest() bool {
	return packet.IsUdp() && packet.GetDstPort() == udwDnsPacket.DnsPort
}

func (packet *IpPacket) GetDnsQueryName() string {
	body := packet.GetUdpBody()
	dnsPacket, errMsg := udwDnsPacket.PacketReadFromByteSlice(body)
	if errMsg != "" {
		return ""
	}
	return dnsPacket.QuestionName
}
