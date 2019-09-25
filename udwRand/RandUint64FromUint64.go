package udwRand

const (
	m1 = 16877499708836156737
	m2 = 2820277070424839065
	m3 = 9497967016996688599
)

func RandUint64FromUint64(srcPort uint64) uint64 {
	h := uint64(m1)
	h ^= uint64(srcPort)
	h = rotl_31(h*m2) + m3
	return h
}

func rotl_31(x uint64) uint64 {
	return (x << 31) | (x >> (64 - 31))
}
