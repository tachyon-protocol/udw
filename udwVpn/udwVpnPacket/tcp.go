package udwVpnPacket

import (
	"bytes"
)

func (packet *IpPacket) IsTcp() bool {
	return packet.buf[9] == 6
}

func (packet *IpPacket) dumpTcpFlag(_buf *bytes.Buffer) {
	if packet.GetTcpFlagSyn() {
		_buf.WriteString("SYN ")
	}
	if packet.GetTcpFlagAck() {
		_buf.WriteString("ACK ")
	}
	if packet.GetTcpFlagPsh() {
		_buf.WriteString("PSH ")
	}
	if packet.GetTcpFlagFin() {
		_buf.WriteString("FIN ")
	}
	if packet.GetTcpFlagRst() {
		_buf.WriteString("RST ")
	}
	if packet.GetTcpBodyLen() > 0 {
		_buf.WriteString("DATA ")
	}
}
