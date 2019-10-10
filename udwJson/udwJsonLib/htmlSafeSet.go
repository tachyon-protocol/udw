package udwJsonLib

var ghtmlSafeSetArray = [2]uint64{0xafffffbb00000000, 0xffffffffefffffff}

func getHtmlSafeSetV2(b uint8) bool {

	return (ghtmlSafeSetArray[b>>6] & (1 << (b & 0x3f))) != 0
}
