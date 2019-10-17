package udwIpPacketFuzzLib

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwIpPacket"
	"github.com/tachyon-protocol/udw/udwTest"
	"net"
)

func Fuzz(data []byte) int {
	data1 := udwBytes.Clone(data)

	ipPacket, errMsg := udwIpPacket.NewIpv4PacketFromBuf(data1)
	if errMsg != "" {
		return 0
	}
	if ipPacket.IsVpnSupport() == false {
		return 0
	}
	checkPacket1NotCrash(ipPacket)
	return 1
}

func checkPacket1NotCrash(packet udwIpPacket.IpPacket) {
	if packet.IsTcp() {
		{
			packet2 := packet.Clone()
			packet2.TcpFixMss(1)
			packet2.GetTcpBody()
			packet2.TcpGetNextSequenceNumber()
			mss := packet2.TcpGetMss()
			udwTest.Ok(mss <= 1)
		}
	} else if packet.IsUdp() {
		packet2 := packet.Clone()
		packet2.GetUdpBody()
	}
	{
		packet2 := packet.Clone()
		packet2.String()
		packet2.CheckSumCorrect()
		packet2.Clone()
		packet2.GetL3ProtocolString()
		packet2.GetDstPort()
		packet2.GetSrcPort()
		packet2.IsIpv4()
		packet2.IsIpv6()
		packet2.IsUdp()
		packet2.IsICMP()
		packet2.Ipv4HasMoreFragments()
		packet2.GetDstAddrString()
		packet2.GetIpHeaderChecksum()
		packet2.GetIpDscpAndEcn()
		packet2.GetIpIdentification()
		packet2.GetIPLen()
		packet2.GetIpHeaderLen()
		packet2.GetIpProtocol()
		packet2.GetIpVersion()
		packet2.SetDstIp(packet2.GetDstIp())
		packet2.SetSrcIp(packet2.GetSrcIp())
		packet2.RecomputeCheckSum()

		packet2.GetSrcDstAddrPeerString()
		packet2.GetDstSrcAddrPeerString()
	}
	if packet.IsUdp() || packet.IsTcp() {
		{
			packet2 := packet.Clone()
			connAddr := packet2.GetConnAddrTcpUdp()
			connAddr.String()
			packet2.SetConnAddrTcpUdp(connAddr.Clone().RevertPeer())
		}
		{
			packet2 := packet.Clone()
			packet2.RecomputeCheckSum()
			udwTest.Equal(packet2.CheckSumCorrect(), "")
			thisBuf := udwBytes.Clone(packet2.SerializeToBuf())

			packet2.RecomputeCheckSum()
			udwTest.Equal(thisBuf, packet2.SerializeToBuf())

		}
	}
	if packet.IsIpv4() {
		{
			packet2 := packet.Clone()
			dstIp := net.IP{127, 0, 0, 1}
			packet2.SetDstIp(dstIp)
			udwTest.Equal(packet2.GetDstIp(), dstIp)
		}
		{
			packet2 := packet.Clone()
			dstIp := net.IP{127, 0, 0, 1}
			packet2.SetSrcIp(dstIp)
			udwTest.Equal(packet2.GetSrcIp(), dstIp)
		}
	}
}
