// +build !amd64

package udwIpCountryV2Map

func (r *Reader) ReadNode(nodeNumber uint32, index uint8) uint32 {
	baseOffset := nodeNumber * 6
	offset := baseOffset + uint32(index*3)
	_ = r.Buf[offset+2]
	return uint32(r.Buf[offset]) | uint32(r.Buf[offset+1])<<8 | uint32(r.Buf[offset+2])<<16
}
