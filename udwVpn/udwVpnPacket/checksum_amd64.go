package udwVpnPacket

import (
	"math/bits"
	"unsafe"
)

func (packet *IpPacket) RecomputeCheckSum() {
	recomputeCheckSum(packet.buf)
}

func recomputeCheckSum(buf []byte) {

	_ = buf[19]
	ipLen := int((buf[0] & 0x0f) * 4)
	var csum uint64
	csum += uint64(*(*uint32)(unsafe.Pointer(&buf[0])))
	csum += uint64(*(*uint32)(unsafe.Pointer(&buf[4])))
	csum += uint64(*(*uint16)(unsafe.Pointer(&buf[8])))
	csum += uint64(*(*uint32)(unsafe.Pointer(&buf[12])))
	csum += uint64(*(*uint32)(unsafe.Pointer(&buf[16])))
	for i := 20; i < ipLen; i += 4 {
		csum += uint64(*(*uint32)(unsafe.Pointer(&buf[i])))

	}

	csum = (csum >> 16) + (csum & 0xffff)
	csum = (csum >> 16) + (csum & 0xffff)
	csum = (csum >> 16) + (csum & 0xffff)

	checksum := ^uint16(csum)
	*(*uint16)(unsafe.Pointer(&buf[10])) = checksum

	if buf[9] == 17 {
		checksum := getUdpCheckSum(buf)
		buf[ipLen+6] = byte((checksum >> 8) & 0xff)
		buf[ipLen+7] = byte(checksum & 0xff)
	} else if buf[9] == 6 {

		_ = buf[ipLen+19]
		var sum uint64

		sum = uint64(6) << 8

		sum += uint64(*(*uint32)(unsafe.Pointer(&buf[12])))

		sum += uint64(*(*uint32)(unsafe.Pointer(&buf[16])))

		tcpLen := len(buf) - ipLen
		sum += uint64(bits.ReverseBytes16(uint16(tcpLen & 0xffff)))
		bufs := uintptr(unsafe.Pointer(&buf[0]))
		bufe := bufs + uintptr(len(buf))
		bufp := bufs + uintptr(ipLen)
		sum += uint64(*(*uint32)(unsafe.Pointer(bufp)))
		bufp += 4
		sum += uint64(*(*uint32)(unsafe.Pointer(bufp)))
		bufp += 4
		sum += uint64(*(*uint32)(unsafe.Pointer(bufp)))
		bufp += 4
		sum += uint64(*(*uint32)(unsafe.Pointer(bufp)))
		bufp += 6

		sum += uint64(*(*uint16)(unsafe.Pointer(bufp)))
		bufp += 2
		if tcpLen > 20 {

			endPos := bufs + uintptr(len(buf)-len(buf)%4)

			if (endPos-bufp)/4 > 0 {

				sum += asmBulkAddSum(bufp, (endPos-bufp)/4)

			}

			bufp = endPos

			remain := bufe - bufp
			switch remain {
			case 3:
				sum += uint64(*(*uint16)(unsafe.Pointer(bufp)))
				bufp += 2
				fallthrough
			case 1:
				sum += uint64(*(*uint8)(unsafe.Pointer(bufp)))
			case 2:
				sum += uint64(*(*uint16)(unsafe.Pointer(bufp)))
			}
		}

		sum = (sum >> 16) + (sum & 0xffff)
		sum = (sum >> 16) + (sum & 0xffff)
		sum = (sum >> 16) + (sum & 0xffff)

		sum2 := ^uint16(sum)
		*((*uint16)(unsafe.Pointer(bufs + uintptr(ipLen) + 16))) = sum2

	}
}

func asmBulkAddSum(bufp uintptr, loopNum uintptr) uint64
