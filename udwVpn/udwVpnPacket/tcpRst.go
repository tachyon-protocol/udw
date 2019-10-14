package udwVpnPacket

func TcpRstSameWay(ipPacket IpPacket, tmpBuf []byte) IpPacket {
	ipLen := ipPacket.GetIpHeaderLen()
	srcPortBuf := ipPacket.buf[ipLen+0 : ipLen+2]
	dstPortBuf := ipPacket.buf[ipLen+2 : ipLen+4]
	reqTcpSeqNumberBuf := ipPacket.buf[ipLen+4 : ipLen+8]

	tmpBuf = append(tmpBuf[:0], gTcpRstTemplate...)
	copy(tmpBuf[12:16], ipPacket.GetSrcIp())
	copy(tmpBuf[16:20], ipPacket.GetDstIp())
	copy(tmpBuf[20:22], srcPortBuf)
	copy(tmpBuf[22:24], dstPortBuf)
	copy(tmpBuf[24:28], reqTcpSeqNumberBuf)

	return IpPacket{
		buf: tmpBuf,
	}
}

func TcpRstAnotherWay(ipPacket IpPacket, tmpBuf []byte) IpPacket {
	ipLen := ipPacket.GetIpHeaderLen()
	srcPortBuf := ipPacket.buf[ipLen+0 : ipLen+2]
	dstPortBuf := ipPacket.buf[ipLen+2 : ipLen+4]

	reqTcpAckNumberBuf := ipPacket.buf[ipLen+8 : ipLen+12]
	tmpBuf = append(tmpBuf[:0], gTcpRstTemplate...)
	copy(tmpBuf[12:16], ipPacket.GetDstIp())
	copy(tmpBuf[16:20], ipPacket.GetSrcIp())
	copy(tmpBuf[20:22], dstPortBuf)
	copy(tmpBuf[22:24], srcPortBuf)
	copy(tmpBuf[24:28], reqTcpAckNumberBuf)

	return IpPacket{
		buf: tmpBuf,
	}
}

var gTcpRstTemplate = []byte{
	0x45, 0x00, 0x00, 0x28, 0x00, 0x00, 0x40, 0x00, 0x40, 0x06, 0, 0, 0xac, 0x15, 0x00, 0x01,
	0x68, 0xc7, 0xe7, 0x82, 0xc5, 0xb0, 0x4e, 0x3c, 0x4f, 0x59, 0x37, 0xd2, 0x00, 0x00, 0x00, 0x00,
	0x50, 0x04, 0x00, 0x00, 0, 0, 0x00, 0x00,
}
