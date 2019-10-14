package udwBitwise

func Uint16SetPartWithMask(origin uint16, toSet uint16, mask uint16) uint16 {
	return origin ^ ((origin ^ toSet) & mask)
}

func Uint16Tail(origin uint16, n uint) uint16 {
	return origin & (1<<n - 1)
}
