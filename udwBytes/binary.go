package udwBytes

import "encoding/binary"

func Uint64ToBigEndSlice(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}
