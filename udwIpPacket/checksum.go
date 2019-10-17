package udwIpPacket

func (packet *IpPacket) CheckSumCorrect() (errMsg string) {

	ipLen := packet.GetIpHeaderLen()
	var csum uint32
	for i := 0; i < ipLen; i += 2 {
		csum += uint32(packet.buf[i]) << 8
		csum += uint32(packet.buf[i+1])
	}
	if uint16((csum>>16)+csum) != 0xffff {
		return "ip checksum fail"
	}

	if packet.IsUdp() {
		if packet.buf[ipLen+6] == 0 && packet.buf[ipLen+7] == 0 {

			return ""
		}
		checksum := packet.getUdpCheckSum()
		if !(byte(checksum>>8&0xff) == packet.buf[ipLen+6] && byte(checksum&0xff) == packet.buf[ipLen+7]) {
			return "udp checksum fail"
		}
	} else if packet.IsTcp() {
		checksum := packet.getTcpCheckSum()
		if !(byte(checksum>>8&0xff) == packet.buf[ipLen+16] && byte(checksum&0xff) == packet.buf[ipLen+17]) {
			return "tcp checksum fail"
		}
	}

	return ""
}

func (packet *IpPacket) getTcpCheckSum() uint16 {
	ipLen := packet.GetIpHeaderLen()
	var sum uint64
	sourceIp := packet.GetSrcIp()
	sum += uint64(sourceIp[0])<<8 + uint64(sourceIp[1])
	sum += uint64(sourceIp[2])<<8 + uint64(sourceIp[3])

	dstIp := packet.GetDstIp()
	sum += uint64(dstIp[0])<<8 + uint64(dstIp[1])
	sum += uint64(dstIp[2])<<8 + uint64(dstIp[3])

	sum += uint64(6)

	tcpLen := len(packet.buf) - ipLen
	sum += uint64(tcpLen & 0xffff)
	for i := ipLen; i < ipLen+16; i += 2 {

		sum += uint64(packet.buf[i])<<8 + uint64(packet.buf[i+1])
	}

	endPos := len(packet.buf) - len(packet.buf)%2
	for i := ipLen + 18; i < endPos; i += 2 {

		sum += uint64(packet.buf[i])<<8 + uint64(packet.buf[i+1])
	}

	if len(packet.buf)%2 != 0 {
		sum += uint64(packet.buf[len(packet.buf)-1]) << 8
	}

	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	return ^uint16(sum + (sum >> 16))
}
