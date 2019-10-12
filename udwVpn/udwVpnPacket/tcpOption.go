package udwVpnPacket

import (
	"strconv"
)

const (
	maxOptionLen                               = 40
	tcpOptionEnd                               = 0x00
	tcpOptionNop                               = 0x01
	tcpOptionMss                               = 0x02
	tcpOptionWindowScale                       = 0x03
	tcpOptionSelectiveAcknowledgementPermitted = 0x04
	tcpOptionSelectiveAcknowledgement          = 0x05
	tcpOptionTimestamp                         = 0x08
)

func checkTcpOptionCorrect(tcpOptionBuf []byte) (errMsg string) {
	i := 0
	size := len(tcpOptionBuf)
	for {
		if i >= size {
			break
		}
		if tcpOptionBuf[i] == tcpOptionEnd || tcpOptionBuf[i] == tcpOptionNop {
			i++
			continue
		}
		if i+1 >= size {
			return "tcp option EOF1"
		}
		optionSize := int(tcpOptionBuf[i+1])
		if optionSize < 2 {
			return "tcp optionSize<2 " + strconv.Itoa(optionSize)
		}
		if i+optionSize > size {
			return "tcp option EOF2"
		}
		i += optionSize
		continue
	}
	return ""
}
