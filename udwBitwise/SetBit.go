package udwBitwise

func Uint8SetBit0(inNum uint8, pos uint8) uint8 {
	return inNum & (^(1 << pos))
}

func Uint16SetBit0(inNum uint16, pos uint8) uint16 {
	return inNum & (^(1 << pos))
}
func Uint32SetBit0(inNum uint32, pos uint8) uint32 {
	return inNum & (^(1 << pos))
}
func Uint64SetBit0(inNum uint64, pos uint8) uint64 {
	return inNum & (^(1 << pos))
}

func Uint8SetBit1(inNum uint8, pos uint8) uint8 {
	return inNum | (1 << pos)
}

func Uint16SetBit1(inNum uint16, pos uint8) uint16 {
	return inNum | (1 << pos)
}
func Uint32SetBit1(inNum uint32, pos uint8) uint32 {
	return inNum | (1 << pos)
}
func Uint64SetBit1(inNum uint64, pos uint8) uint64 {
	return inNum | (1 << pos)
}

func Uint16SetBitWithBool(inNum uint16, pos uint8, is bool) uint16 {
	if is {
		return Uint16SetBit1(inNum, pos)
	} else {
		return Uint16SetBit0(inNum, pos)
	}
}
func Uint64SetBitWithBool(inNum uint64, pos uint8, is bool) uint64 {
	if is {
		return Uint64SetBit1(inNum, pos)
	} else {
		return Uint64SetBit0(inNum, pos)
	}
}

func Uint8GetBitToBool(inNum uint8, pos uint8) bool {
	return (inNum & (1 << pos)) != 0
}
func Uint16GetBitToBool(inNum uint16, pos uint8) bool {
	return (inNum & (1 << pos)) != 0
}
func Uint64GetBitToBool(inNum uint64, pos uint8) bool {
	return (inNum & (1 << pos)) != 0
}

func Uint16GetBitTo1Or0Uint8(inNum uint16, pos uint8) uint8 {
	if Uint16GetBitToBool(inNum, pos) {
		return 1
	} else {
		return 0
	}
}

func Uint32GetBitTo1Or0Uint8(inNum uint32, pos int) uint8 {
	return uint8((inNum >> uint8(pos)) & 1)
}
func Uint32GetBit(inNum uint32, pos int) uint8 {
	return uint8((inNum >> uint8(pos)) & 1)
}

func Uint32SetBit(inNum uint32, pos int, value uint8) uint32 {
	if value == 1 {
		return Uint32SetBit1(inNum, uint8(pos))
	} else {
		return Uint32SetBit0(inNum, uint8(pos))
	}
}

func BoolTrueToUint8_1(b bool) uint8 {
	if b {
		return 1
	} else {
		return 0
	}
}
