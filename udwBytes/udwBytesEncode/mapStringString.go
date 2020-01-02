package udwBytesEncode

import "github.com/tachyon-protocol/udw/udwBytes"

func MapStringStringMarshal(m map[string]string) []byte {
	if m == nil || len(m) == 0 {
		return []byte{0}
	}
	_buf := udwBytes.BufWriter{}
	_buf.WriteUvarint(uint64(len(m)))
	for k, v := range m {
		_buf.WriteStringLenUvarint(k)
		_buf.WriteStringLenUvarint(v)
	}
	return _buf.GetBytes()
}

func MapStringStringUnmarshal(b []byte) (m map[string]string, ok bool) {
	_buf := udwBytes.NewBufReaderWithOk(b)
	l := _buf.ReadUvarint()
	if _buf.IsOk() == false {
		return nil, false
	}
	if l == 0 {
		return map[string]string{}, true
	}
	if l > uint64(_buf.GetRemainSize()/2) {
		return nil, false
	}
	m = make(map[string]string, int(l))
	for i := 0; i < int(l); i++ {
		k := _buf.ReadStringLenUvarint()
		v := _buf.ReadStringLenUvarint()
		if _buf.IsOk() == false {
			return nil, false
		}
		m[k] = v
	}
	return m, true
}
