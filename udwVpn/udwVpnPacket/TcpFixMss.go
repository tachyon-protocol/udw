package udwVpnPacket

import (
	"encoding/binary"
)

func (packet *IpPacket) TcpFixMss(mss uint16) {
	hasChange := packet.TcpFixMss__NoRecomputeCheckSum(mss)
	if hasChange {
		packet.RecomputeCheckSum()
	}
}

func (packet *IpPacket) TcpFixMss__NoRecomputeCheckSum(mss uint16) bool {
	if !packet.IsTcp() {
		return false
	}

	dataOffset := packet.GetTcpDataOffset()
	if dataOffset <= 20 {
		return false
	}
	ihl := packet.GetIpHeaderLen()
	if dataOffset+ihl > len(packet.buf) {
		return false
	}

	tcpOptionBuf := packet.buf[ihl+20 : ihl+dataOffset]
	size := len(tcpOptionBuf)
	i := 0
	hasChange := false
	for {
		if i >= size {
			break
		}

		if tcpOptionBuf[i] == tcpOptionEnd || tcpOptionBuf[i] == tcpOptionNop {
			i++
			continue
		}
		optionSize := int(tcpOptionBuf[i+1])
		if i+optionSize > size {
			break
		}
		if i+2 > size {
			break
		}
		if tcpOptionBuf[i] != tcpOptionMss || optionSize < 4 {
			i += optionSize
			continue
		}
		oldMss := binary.BigEndian.Uint16(tcpOptionBuf[i+2 : i+optionSize])
		if oldMss <= mss {
			i += optionSize
			continue
		}
		hasChange = true
		binary.BigEndian.PutUint16(tcpOptionBuf[i+2:i+optionSize], mss)

		i += optionSize
		continue
	}
	return hasChange
}

func (packet *IpPacket) TcpGetMss() uint16 {
	if !packet.IsTcp() {
		return 0
	}

	dataOffset := packet.GetTcpDataOffset()
	if dataOffset <= 20 {
		return 0
	}
	ihl := packet.GetIpHeaderLen()
	if dataOffset+ihl > len(packet.buf) {
		return 0
	}

	tcpOptionBuf := packet.buf[ihl+20 : ihl+dataOffset]
	size := len(tcpOptionBuf)
	i := 0
	for {
		if i >= size {
			break
		}

		if tcpOptionBuf[i] == tcpOptionEnd || tcpOptionBuf[i] == tcpOptionNop {
			i++
			continue
		}
		optionSize := int(tcpOptionBuf[i+1])
		if i+optionSize > size {
			break
		}
		if i+2 > size {
			break
		}
		if tcpOptionBuf[i] != tcpOptionMss || optionSize < 4 {
			i += optionSize
			continue
		}
		oldMss := binary.BigEndian.Uint16(tcpOptionBuf[i+2 : i+optionSize])
		return oldMss
	}
	return 0
}
