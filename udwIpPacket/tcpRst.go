package udwIpPacket

import (
	"github.com/tachyon-protocol/udw/udwNet/udwDns/udwDnsPacket"
)

func TcpRstSameWay(ipPacket IpPacket, tmpBuf []byte) IpPacket {
	ipLen := ipPacket.GetIpHeaderLen()
	srcPortBuf := ipPacket.buf[ipLen+0 : ipLen+2]
	dstPortBuf := ipPacket.buf[ipLen+2 : ipLen+4]
	reqTcpSeqNumberBuf := ipPacket.buf[ipLen+4 : ipLen+8]
	reqTcpAckNumberBuf := ipPacket.buf[ipLen+8 : ipLen+12]
	tmpBuf = append(tmpBuf[:0], gTcpRstTemplate...)
	copy(tmpBuf[12:16], ipPacket.GetSrcIp())
	copy(tmpBuf[16:20], ipPacket.GetDstIp())
	copy(tmpBuf[20:22], srcPortBuf)
	copy(tmpBuf[22:24], dstPortBuf)
	copy(tmpBuf[24:28], reqTcpSeqNumberBuf)
	copy(tmpBuf[28:32], reqTcpAckNumberBuf)
	return IpPacket{
		buf: tmpBuf,
	}
}

func TcpRstAnotherWay(ipPacket IpPacket, tmpBuf []byte) IpPacket {
	ipLen := ipPacket.GetIpHeaderLen()
	srcPortBuf := ipPacket.buf[ipLen+0 : ipLen+2]
	dstPortBuf := ipPacket.buf[ipLen+2 : ipLen+4]
	reqTcpSeqNumberBuf := ipPacket.buf[ipLen+4 : ipLen+8]
	reqTcpAckNumberBuf := ipPacket.buf[ipLen+8 : ipLen+12]
	if ipPacket.GetTcpFlagSyn() && ipPacket.GetTcpFlagAck() == false {

		for i := 3; i >= 0; i-- {
			reqTcpSeqNumberBuf[i]++
			if reqTcpSeqNumberBuf[i] != 0 {
				break
			}
		}
	}
	tmpBuf = append(tmpBuf[:0], gTcpRstTemplate...)
	copy(tmpBuf[12:16], ipPacket.GetDstIp())
	copy(tmpBuf[16:20], ipPacket.GetSrcIp())
	copy(tmpBuf[20:22], dstPortBuf)
	copy(tmpBuf[22:24], srcPortBuf)
	copy(tmpBuf[24:28], reqTcpAckNumberBuf)
	copy(tmpBuf[28:32], reqTcpSeqNumberBuf)
	outIpPacket := IpPacket{
		buf: tmpBuf,
	}
	return outIpPacket
}

var gTcpRstTemplate = []byte{
	0x45, 0x00, 0x00, 0x28, 0x00, 0x00, 0x40, 0x00, 0x40, 0x06, 0, 0, 0xac, 0x15, 0x00, 0x01,
	0x68, 0xc7, 0xe7, 0x82, 0xc5, 0xb0, 0x4e, 0x3c, 0x4f, 0x59, 0x37, 0xd2, 0x00, 0x00, 0x00, 0x00,
	0x50, 0x14, 0x00, 0x00, 0, 0, 0x00, 0x00,
}

func NotifyBlockIpPacketToClient(ipPacket IpPacket) (outIpPacket IpPacket, ok bool) {
	if ipPacket.IsDnsRequest() {

		dnsPacket, errMsg := udwDnsPacket.PacketReadFromByteSlice(ipPacket.GetUdpBody())
		if errMsg != "" {
			return outIpPacket, false
		}
		dnsPacket.AnswerList = nil
		dnsPacket.SetIsResponse(true)
		dnsPacket.SetRcode(udwDnsPacket.RcodeREFUSED)
		outBuf, errMsg := udwDnsPacket.PacketWriteToByteSlice(dnsPacket, nil)
		if errMsg != "" {
			return outIpPacket, false
		}
		ipPacket2 := NewUdpIpPacket(NewUdpIpPacketRequest{
			SrcIp:   ipPacket.GetDstIp(),
			SrcPort: ipPacket.GetDstPort(),
			DstIp:   ipPacket.GetSrcIp(),
			DstPort: ipPacket.GetSrcPort(),
			Body:    outBuf,
		})
		return ipPacket2, true
	}
	if ipPacket.IsTcp() && ipPacket.GetTcpFlagRst() == false {
		ipPacket2 := TcpRstAnotherWay(ipPacket, nil)
		return ipPacket2, true
	}
	return outIpPacket, false
}
