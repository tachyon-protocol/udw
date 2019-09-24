package udwBytes

func WriteUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

func ReadUvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i, b := range buf {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, -(i + 1)
			}
			return x | uint64(b)<<s, i + 1
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

func GetUvarintOutputSize(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

func GetVarintOutputSize(x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return GetUvarintOutputSize(ux)
}

func (bw *BufWriter) WriteUvarint(x uint64) {
	for x >= 0x80 {
		bw.WriteByte_(byte(x) | 0x80)
		x >>= 7
	}
	bw.WriteByte_(byte(x))
}

func (bw *BufWriter) WriteStringLenUvarint(s string) {
	bw.WriteUvarint(uint64(len(s)))
	bw.WriteString_(s)
}

func (bw *BufWriter) WriteStringListLenUvarint(sList []string) {
	bw.WriteUvarint(uint64(len(sList)))
	for _, s := range sList {
		bw.WriteStringLenUvarint(s)
	}
}

func (bw *BufWriter) WriteVarint(x int64) {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	bw.WriteUvarint(ux)
}

func (br *BufReader) ReadUvarint() (x uint64, isOk bool) {
	var s uint
	i := 0
	for {
		b, isRead := br.ReadByteOrEof()
		if isRead == false {
			return 0, false
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, false
			}
			return x | uint64(b)<<s, true
		}
		x |= uint64(b&0x7f) << s
		s += 7
		i++
	}
}
func (br *BufReader) ReadStringLenUvarint() (s string, isOk bool) {
	x, isOk := br.ReadUvarint()
	if isOk == false {
		return "", false
	}
	buf := br.ReadMaxByteNum(int(x))
	if len(buf) != int(x) {
		return "", false
	}
	return string(buf), true
}

func (br *BufReader) ReadStringListLenUvarint() (sList []string, isOk bool) {
	l, isOk := br.ReadUvarint()
	if isOk == false {
		return nil, false
	}
	if l > uint64(br.GetRemainSize()) {
		return nil, false
	}
	out := make([]string, l)
	for i := 0; i < int(l); i++ {
		out[i], isOk = br.ReadStringLenUvarint()
		if isOk == false {
			return nil, false
		}
	}
	return out, true
}

func (crb *ReadBufferWrap) ReadUvarint() (x uint64, errMsg string) {
	var s uint
	i := 0
	for {
		b, errMsg := crb.ReadByteErrMsg()
		if errMsg != "" {
			return 0, errMsg
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, "2vg5peydjt overflow"
			}
			return x | uint64(b)<<s, ""
		}
		x |= uint64(b&0x7f) << s
		s += 7
		i++
	}
}
func (crb *ReadBufferWrap) ReadStringLenUvarint() (s string, errMsg string) {
	x, errMsg := crb.ReadUvarint()
	if errMsg != "" {
		return "", errMsg
	}
	buf, errMsg := crb.ReadBySize(int(x))
	if errMsg != "" {
		return "", errMsg
	}
	return string(buf), ""
}

func (br *BufReader) ReadVarint() (x int64, isOk bool) {
	ux, isOk := br.ReadUvarint()
	if isOk == false {
		return 0, isOk
	}
	x = int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, true
}

func DecodeStringListLenUvarint(b []byte) ([]string, bool) {
	r := NewBufReader(b)
	return r.ReadStringListLenUvarint()
}

func EncodeStringListLenUvarint(sList []string) (b []byte) {
	bufW := &BufWriter{}
	bufW.WriteStringListLenUvarint(sList)
	return bufW.GetBytes()
}
