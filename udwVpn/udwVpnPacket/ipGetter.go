package udwVpnPacket

import (
	"encoding/binary"
	"net"
)

func (packet *IpPacket) GetConnAddrTcpUdp() ConnAddrTcpUdp {
	return ConnAddrTcpUdp{
		SrcIp:   packet.GetSrcIp(),
		SrcPort: packet.GetSrcPort(),
		DstIp:   packet.GetDstIp(),
		DstPort: packet.GetDstPort(),
	}
}

func (packet *IpPacket) GetSrcIp() net.IP {
	return net.IP(packet.buf[12:16])
}

func (packet *IpPacket) GetSrcPort() uint16 {
	if !packet.IsTcp() && !packet.IsUdp() {
		return 0
	}
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint16(packet.buf[ipLen : ipLen+2])
}
func (packet *IpPacket) SetSrcPort__NoRecomputeChecksum(port uint16) {
	if !packet.IsTcp() && !packet.IsUdp() {
		return
	}
	ipLen := packet.GetIpHeaderLen()
	binary.BigEndian.PutUint16(packet.buf[ipLen:ipLen+2], port)
}

func (packet *IpPacket) GetDstIp() net.IP {
	return net.IP(packet.buf[16:20])
}

func (packet *IpPacket) GetDstPort() uint16 {
	if !packet.IsTcp() && !packet.IsUdp() {
		return 0
	}
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint16(packet.buf[ipLen+2 : ipLen+4])
}
func (packet *IpPacket) SetDstPort__NoRecomputeChecksum(port uint16) {
	if !packet.IsTcp() && !packet.IsUdp() {
		return
	}
	ipLen := packet.GetIpHeaderLen()
	binary.BigEndian.PutUint16(packet.buf[ipLen+2:ipLen+4], port)
}

func (packet *IpPacket) GetIpVersion() int {
	return int((packet.buf[0] & 0xf0) >> 4)
}

func (packet *IpPacket) GetIpHeaderLen() int {
	return int((packet.buf[0] & 0x0f) * 4)
}

func (packet *IpPacket) GetIpDscpAndEcn() byte {
	return packet.buf[1]
}

func (packet *IpPacket) GetIpTotalLen() uint16 {
	return binary.BigEndian.Uint16(packet.buf[2:4])
}

func (packet *IpPacket) GetIpIdentification() uint16 {
	return binary.BigEndian.Uint16(packet.buf[4:6])
}

func (packet *IpPacket) GetIpFlagsAndFragmentOffset() uint16 {
	return binary.BigEndian.Uint16(packet.buf[6:8])
}

func (packet *IpPacket) GetIpTimeToLive() byte {
	return packet.buf[8]
}

func (packet *IpPacket) GetIpProtocol() byte {
	return packet.buf[9]
}

func (packet *IpPacket) GetIpHeaderChecksum() uint16 {
	return binary.BigEndian.Uint16(packet.buf[10:12])
}

func (packet *IpPacket) GetIpOption() []byte {
	ihl := packet.GetIpHeaderLen()
	if ihl > 20 {
		return packet.buf[20:ihl]
	}
	return nil
}
