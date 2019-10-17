package udwIpPacket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwBinary"
	"github.com/tachyon-protocol/udw/udwHex"
	"github.com/tachyon-protocol/udw/udwNet/udwDns/udwDnsPacket"
	"strconv"
	"sync/atomic"
)

type IpPacket struct {
	buf []byte
}

func NewIpv4PacketFromBuf(buf []byte) (ipPacket IpPacket, errMsg string) {
	errMsg = CheckIPPacket(buf)
	if errMsg != "" {
		return ipPacket, errMsg
	}
	ipPacket.buf = buf
	if !ipPacket.IsVpnSupport() {
		return ipPacket, "vqw396d9w5 udwVpnPacket NewIpv4PacketFromBuf vpn not support"
	}
	return ipPacket, ""
}

func NewIpPacketFromBuf(buf []byte) (ipPacket IpPacket, errMsg string) {
	errMsg = CheckIPPacket(buf)
	if errMsg != "" {
		return ipPacket, errMsg
	}
	return IpPacket{buf: buf}, ""
}

func (packet *IpPacket) ResetIpPacketFromBuf(buf []byte) (err error) {
	errMsg := CheckIPPacket(buf)
	if errMsg != "" {
		return errors.New(errMsg)
	}
	packet.buf = buf
	return nil
}

func (packet *IpPacket) SerializeToBuf() []byte {
	return packet.buf
}

func (packet *IpPacket) Clone() IpPacket {
	buf := make([]byte, len(packet.buf))
	copy(buf, packet.buf)
	return IpPacket{
		buf: buf,
	}
}

func (packet *IpPacket) String() string {
	if packet == nil || len(packet.buf) == 0 {
		return "<nil>"
	}
	_buf := &bytes.Buffer{}
	if packet.IsIpv4() {
		_buf.WriteString("ipv4 ")
		if packet.Ipv4HasMoreFragments() {

			_buf.WriteString(" hasMoreFramgments ")
			_buf.WriteString(packet.GetSrcIp().String())
			_buf.WriteString("-")
			_buf.WriteString(packet.GetDstIp().String())
			return _buf.String()
		}
		if packet.IsUdp() {
			_buf.WriteString("udp ")
			_buf.WriteString(packet.GetSrcDstAddrPeerString())
			_buf.WriteString(fmt.Sprintf(" len: %d ", len(packet.GetUdpBody())))
			if packet.GetDstPort() == 53 || packet.GetSrcPort() == 53 {
				dnsPacket, errMsg := udwDnsPacket.PacketReadFromByteSlice(packet.GetUdpBody())
				if errMsg == "" {
					_buf.WriteString(" dns ")
					_buf.WriteString(dnsPacket.String())
				}
			}
			return _buf.String()
		} else if packet.IsTcp() {
			_buf.WriteString("tcp ")
			_buf.WriteString(packet.GetSrcDstAddrPeerString())
			_buf.WriteString(" ")
			packet.dumpTcpFlag(_buf)
			_buf.WriteString(fmt.Sprintf("seq: %08x ", packet.GetTcpSequenceNumber()))
			_buf.WriteString(fmt.Sprintf("ack: %08x ", packet.GetTcpAcknowledgmentNumber()))
			_buf.WriteString(fmt.Sprintf("win: %04x ", packet.GetTcpWindowSize()))
			_buf.WriteString(fmt.Sprintf("len: %d ", packet.GetTcpBodyLen()))
			mss := packet.TcpGetMss()
			if mss > 0 {
				_buf.WriteString(fmt.Sprintf("mss: %d ", mss))
			}

			return _buf.String()
		} else if packet.IsICMP() {
			_buf.WriteString("ICMP ")
			_buf.WriteString(packet.GetSrcIp().String())
			_buf.WriteString("-")
			_buf.WriteString(packet.GetDstIp().String())
			_buf.WriteString(" ")
			_buf.WriteString(fmt.Sprintf("Type: %d ", packet.ICMPGetType()))
			_buf.WriteString(fmt.Sprintf("Code: %d ", packet.ICMPGetCode()))
			return _buf.String()
		} else {
			_buf.WriteString(fmt.Sprintf("ipProto: %d", packet.buf[9]))
			_buf.WriteString(packet.GetSrcIp().String())
			_buf.WriteString("-")
			_buf.WriteString(packet.GetDstIp().String())
			_buf.WriteString(" ")
			return _buf.String()
		}
	} else {
		_buf.WriteString("ipv6 ")
		return _buf.String()
	}
}

func (packet *IpPacket) HexDump() string {
	return udwHex.EncodeBytesToString(packet.buf)
}

func (packet *IpPacket) IsICMP() bool {
	return packet.buf[9] == 1
}
func (packet *IpPacket) IsIpv6() bool {
	return packet.buf[0]&0xf0 == 0x60
}
func (packet *IpPacket) IsIpv4() bool {
	return packet.buf[0]&0xf0 == 0x40
}

func (packet *IpPacket) Ipv4HasMoreFragments() bool {
	return packet.buf[6]&0x20 == 0x20
}

func (packet *IpPacket) ICMPGetType() uint8 {
	ipLen := packet.GetIpHeaderLen()
	return packet.buf[ipLen]
}
func (packet *IpPacket) ICMPGetCode() uint8 {
	ipLen := packet.GetIpHeaderLen()
	return packet.buf[ipLen+1]
}

func (packet *IpPacket) GetL3ProtocolString() string {
	if packet.IsIpv6() {
		return "ipv6"
	}
	if packet.IsUdp() {
		return "udp"
	}
	if packet.IsTcp() {
		return "tcp"
	}
	if packet.IsICMP() {
		return "icmp"
	}
	return "unknowProtocol" + strconv.Itoa(int(packet.buf[9]))
}

func (packet *IpPacket) GetIPLen() int {
	return len(packet.buf)
}

func CheckIPPacket(buf []byte) (errMsg string) {
	if len(buf) < 20 {
		return "[checkIPv4Packet] ip check len(buf)[" + strconv.Itoa(len(buf)) + "]<20"
	}

	if buf[0]&0xf0 == 0x40 {
		ipLen := int((buf[0] & 0x0f) * 4)
		if ipLen < 20 {

			return fmt.Sprintf("[checkIPv4Packet] ipv4 check ipLen[%d]<20", len(buf))
		}
		if len(buf) < ipLen {
			return fmt.Sprintf("[checkIPv4Packet] ipv4 check len(buf)[%d]<ipLen[%d]", len(buf), ipLen)
		}
		ipTotalLen := binary.BigEndian.Uint16(buf[2:4])
		if int(ipTotalLen) != len(buf) {
			return fmt.Sprintf("[checkIPv4Packet] ipv4 check ipTotalLen[%d]!=len(buf)[%d]", ipTotalLen, len(buf))
		}

		if buf[6]&0x20 == 0x20 {

			return ""
		}
		if buf[9] == 6 {

			if len(buf) < ipLen+20 {
				return fmt.Sprintf("[checkIPv4Packet] tcp check len(buf)[%d]<ipLen+20[%d]", len(buf), ipLen+20)
			}
			tcpHeaderLen := int((buf[ipLen+12] & 0xf0) / 16 * 4)
			if tcpHeaderLen < 20 {
				return fmt.Sprintf("[checkIPv4Packet] tcp check tcpHeaderLen[%d]<20", tcpHeaderLen)
			}
			if len(buf) < ipLen+tcpHeaderLen {
				return fmt.Sprintf("[checkIPv4Packet] tcp check len(buf)[%d]<ipLen+tcpHeaderLen[%d]", len(buf), ipLen+tcpHeaderLen)
			}
			if tcpHeaderLen > 20 {
				errMsg := checkTcpOptionCorrect(buf[ipLen+20 : ipLen+tcpHeaderLen])
				if errMsg != "" {
					return errMsg
				}
			}
			return ""
		}
		if buf[9] == 17 {
			if len(buf) < ipLen+8 {
				return fmt.Sprintf("[checkIPv4Packet] udp check len(buf)[%d]<ipLen+8[%d]", len(buf), ipLen+8)
			}
			udpLen := int(buf[ipLen+4])*256 + int(buf[ipLen+5])
			if len(buf) != ipLen+udpLen {
				return fmt.Sprintf("[checkIPv4Packet] udp check len(buf)[%d]!=ipLen+udpLen[%d]", len(buf), ipLen+udpLen)
			}
			return ""
		}
		if buf[9] == 1 {
			if len(buf) < ipLen+2 {
				return fmt.Sprintf("[checkIPv4Packet] icmp check len(buf)[%d]<ipLen+2[%d]", len(buf), ipLen+2)
			}
		}

		return ""
	} else {
		if len(buf) < 40 {
			return fmt.Sprintf("[checkIPv4Packet] ipv6 check len(buf)[%d]<40", len(buf))
		}

		return ""
	}
}

var gIpHeaderId uint32 = 82246225

func nextIpHeaderId() uint16 {
	return uint16(atomic.AddUint32(&gIpHeaderId, 1))
}

func (ipPacket *IpPacket) ConnKeyStringDstSrcToBuf(buf *bytes.Buffer) {
	buf.WriteString(ipPacket.GetL3ProtocolString())
	buf.WriteByte('-')
	buf.Write([]byte(ipPacket.GetDstIp()))
	buf.WriteByte(':')
	udwBinary.WriteUint16(buf, ipPacket.GetDstPort())
	buf.WriteByte('-')
	buf.Write([]byte(ipPacket.GetSrcIp()))
	buf.WriteByte(':')
	udwBinary.WriteUint16(buf, ipPacket.GetSrcPort())
	return
}

func (ipPacket *IpPacket) ConnKeyStringSrcDstToBuf(buf *bytes.Buffer) {
	buf.WriteString(ipPacket.GetL3ProtocolString())
	buf.WriteByte('-')
	buf.Write([]byte(ipPacket.GetSrcIp()))
	buf.WriteByte(':')
	udwBinary.WriteUint16(buf, ipPacket.GetSrcPort())
	buf.WriteByte('-')
	buf.Write([]byte(ipPacket.GetDstIp()))
	buf.WriteByte(':')
	udwBinary.WriteUint16(buf, ipPacket.GetDstPort())
	return
}

func (ipPacket *IpPacket) IsVpnSupport() bool {
	return ipPacket.IsIpv4() && !ipPacket.Ipv4HasMoreFragments()
}
