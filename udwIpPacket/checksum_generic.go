// +build !amd64

package udwIpPacket

func (packet *IpPacket) RecomputeCheckSum() {

	ipLen := packet.GetIpHeaderLen()
	var csum uint64
	for i := 0; i < 10; i += 2 {

		csum += uint64(packet.buf[i])<<8 + uint64(packet.buf[i+1])
	}

	for i := 12; i < ipLen; i += 2 {

		csum += uint64(packet.buf[i])<<8 + uint64(packet.buf[i+1])
	}
	csum = uint64(csum>>16) + uint64(csum&0xffff)
	csum = uint64(uint64(csum>>16) + uint64(csum&0xffff)&0xffff)
	checksum := ^uint16(csum)
	packet.buf[10] = byte((checksum >> 8))
	packet.buf[11] = byte(checksum)

	if packet.IsUdp() {
		checksum := packet.getUdpCheckSum()
		packet.buf[ipLen+6] = byte((checksum >> 8) & 0xff)
		packet.buf[ipLen+7] = byte(checksum & 0xff)
	} else if packet.IsTcp() {
		checksum := packet.getTcpCheckSum()
		packet.buf[ipLen+16] = byte((checksum >> 8) & 0xff)
		packet.buf[ipLen+17] = byte(checksum & 0xff)
	}
}
