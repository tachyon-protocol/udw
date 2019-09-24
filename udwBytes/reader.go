package udwBytes

type BufReader struct{
	buf []byte
	pos int
}

func NewBufReader(buf []byte) *BufReader {
	return &BufReader{
		buf: buf,
	}
}

func (r *BufReader) ReadByteOrEof() (b byte,isRead bool){
	if r.pos>=len(r.buf){
		return 0,false
	}
	b = r.buf[r.pos]
	r.pos+=1
	return b,true
}

func (br *BufReader) ReadUvarint() (x uint64,isOk bool){
	var s uint
	i:=0
	for{
		b,isRead:=br.ReadByteOrEof()
		if isRead==false{
			return 0,false // read eof
		}
		if b < 0x80{
			if i > 9 || i == 9 && b > 1 {
				return 0, false // overflow
			}
			return x | uint64(b)<<s, true
		}
		x |= uint64(b&0x7f) << s
		s += 7
		i++
	}
}
func (br *BufReader) ReadStringLenUvarint()(s string,isOk bool){
	x,isOk:=br.ReadUvarint()
	if isOk==false{
		return "",false
	}
	buf,isOk:=br.ReadByteNumOrEof(int(x))
	if isOk==false{
		return "",false
	}
	return string(buf),true
}
func (r *BufReader) ReadByteNumOrEof(num int) (b []byte,isOk bool){
	startPos:=r.pos
	if startPos>=len(r.buf){
		return nil,false
	}else if r.pos+num<len(r.buf){
		r.pos+=num
	}else{
		r.pos = len(r.buf)
	}
	return r.buf[startPos:r.pos],true
}