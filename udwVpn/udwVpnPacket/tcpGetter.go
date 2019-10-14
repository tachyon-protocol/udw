package udwVpnPacket

import (
	"encoding/binary"
)

func (packet *IpPacket) GetTcpSequenceNumber() uint32 {
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint32(packet.buf[ipLen+4 : ipLen+8])
}

func (packet *IpPacket) SetTcpSequenceNumber__NoRecomputeChecksum(seq uint32) {
	ipLen := packet.GetIpHeaderLen()
	binary.BigEndian.PutUint32(packet.buf[ipLen+4:ipLen+8], seq)
}

func (packet *IpPacket) TcpGetNextSequenceNumber() uint32 {
	seqNum := packet.GetTcpSequenceNumber()
	seqNum += uint32(packet.GetTcpBodyLen())
	if packet.GetTcpFlagFin() {
		seqNum++
	}
	if packet.GetTcpFlagSyn() {
		seqNum++
	}
	return seqNum
}

func (packet *IpPacket) GetTcpAcknowledgmentNumber() uint32 {
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint32(packet.buf[ipLen+8 : ipLen+12])
}

func (packet *IpPacket) SetTcpAcknowledgmentNumber__NoRecomputeChecksum(ack uint32) {
	ipLen := packet.GetIpHeaderLen()
	binary.BigEndian.PutUint32(packet.buf[ipLen+8:ipLen+12], ack)
}

func (packet *IpPacket) GetTcpDataOffset() int {
	return int((packet.buf[packet.GetIpHeaderLen()+12] >> 4) * 4)
}

func (packet *IpPacket) GetTcpFlagNs() bool {
	return (packet.buf[packet.GetIpHeaderLen()+12] & 1) == 1
}

func (packet *IpPacket) GetTcpFlagCwr() bool {
	i := byte(0x80)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagEce() bool {
	i := byte(0x40)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagUrg() bool {
	i := byte(0x20)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagAck() bool {
	i := byte(0x10)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagPsh() bool {
	i := byte(8)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagRst() bool {
	i := byte(4)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagSyn() bool {
	i := byte(2)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpFlagFin() bool {
	i := byte(1)
	return (packet.buf[packet.GetIpHeaderLen()+13] & i) == i
}

func (packet *IpPacket) GetTcpWindowSize() uint16 {
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint16(packet.buf[ipLen+14 : ipLen+16])
}

func (packet *IpPacket) GetTcpChecksum() uint16 {
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint16(packet.buf[ipLen+16 : ipLen+18])
}

func (packet *IpPacket) GetTcpUrgentPointer() uint16 {
	ipLen := packet.GetIpHeaderLen()
	return binary.BigEndian.Uint16(packet.buf[ipLen+18 : ipLen+20])
}

func (packet *IpPacket) GetTcpBody() []byte {
	return packet.buf[packet.GetIpHeaderLen()+packet.GetTcpDataOffset():]
}

func (packet *IpPacket) GetTcpBodyLen() int {
	return len(packet.buf) - packet.GetIpHeaderLen() - packet.GetTcpDataOffset()
}
