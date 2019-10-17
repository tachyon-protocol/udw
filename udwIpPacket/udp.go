package udwIpPacket

import (
	"fmt"
	"net"
)

type NewUdpIpPacketRequest struct {
	SrcIp   net.IP
	SrcPort uint16
	DstIp   net.IP
	DstPort uint16
	Body    []byte
}

func NewUdpIpPacket(req NewUdpIpPacketRequest) IpPacket {
	if len(req.Body) > 65507 {
		panic(fmt.Errorf("[NewUdpIpPacket] udp body %d too big", len(req.Body)))
	}
	totalLen := 20 + 8 + len(req.Body)
	buf := make([]byte, totalLen)

	buf[0] = 0x45

	buf[2] = uint8(totalLen >> 8)
	buf[3] = uint8(totalLen)

	id := nextIpHeaderId()
	buf[4] = uint8(id >> 8)
	buf[5] = uint8(id)
	buf[8] = uint8(0xff)
	buf[9] = 17
	copy(buf[12:16], []byte(req.SrcIp.To4()))
	copy(buf[16:20], []byte(req.DstIp.To4()))

	buf[20] = uint8(req.SrcPort >> 8)
	buf[21] = uint8(req.SrcPort)

	buf[22] = uint8(req.DstPort >> 8)
	buf[23] = uint8(req.DstPort)
	udpLen := len(req.Body) + 8
	buf[24] = uint8(udpLen >> 8)
	buf[25] = uint8(udpLen)

	copy(buf[28:], req.Body)
	packet, errMsg := NewIpv4PacketFromBuf(buf)
	if errMsg != "" {

		panic(errMsg)
	}
	packet.RecomputeCheckSum()
	return packet
}

func (packet *IpPacket) getUdpCheckSum() uint16 {
	return getUdpCheckSum(packet.buf)
}

func getUdpCheckSum(buf []byte) uint16 {
	ipLen := int((buf[0] & 0x0f) * 4)
	var csum uint64
	csum += uint64(buf[12]) << 8
	csum += uint64(buf[13])
	csum += uint64(buf[14]) << 8
	csum += uint64(buf[15])

	csum += uint64(buf[16]) << 8
	csum += uint64(buf[17])
	csum += uint64(buf[18]) << 8
	csum += uint64(buf[19])

	csum += uint64(17)

	csum += uint64(buf[ipLen+4]) << 8
	csum += uint64(buf[ipLen+5])

	for i := ipLen; i < ipLen+6; i += 2 {
		csum += uint64(buf[i])<<8 + uint64(buf[i+1])
	}

	for i := ipLen + 8; i < len(buf); i += 2 {
		csum += uint64(buf[i]) << 8
		if i+1 < len(buf) {
			csum += uint64(buf[i+1])
		}
	}
	checksum := ^uint16((csum >> 16) + csum)
	if checksum == 0 {
		checksum = 0xffff
	}
	return checksum
}

func (packet *IpPacket) IsUdp() bool {
	return packet.buf[9] == 17
}

func (packet *IpPacket) GetUdpBody() []byte {
	ipLen := packet.GetIpHeaderLen()
	return packet.buf[ipLen+8:]
}
