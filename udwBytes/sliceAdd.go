package udwBytes

func SliceAddLittleEndianFixLen(b []byte, toAdd uint64) {
	var c uint8
	for i := 0; i < len(b); i++ {
		thisOut := uint16(b[i]) + uint16(toAdd&0xff) + uint16(c)
		c = uint8(thisOut >> 8 & 0xff)
		b[i] = uint8(thisOut & 0xff)
		toAdd = toAdd >> 8
		if toAdd == 0 && c == 0 {
			break
		}
	}
	return
}
