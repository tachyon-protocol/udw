package udwBytes

type BufWriter struct{
	buf []byte
}

func (w *BufWriter) Write_(p []byte){
	w.buf = append(w.buf,p...)
}

func (w *BufWriter) Write(p []byte)(n int,err error){
	w.buf = append(w.buf,p...)
	return len(p),nil
}

func (w *BufWriter) WriteString_(s string){
	w.buf = append(w.buf,s...)
}

func (w *BufWriter) WriteString(p string)(n int,err error){
	w.buf = append(w.buf,p...)
	return len(p),nil
}

func (w *BufWriter) WriteByte_(b uint8) {
	w.buf = append(w.buf,b)
}

func (w *BufWriter) WriteByte(b uint8) error{
	w.buf = append(w.buf,b)
	return nil
}

func (bw *BufWriter) WriteUvarint(x uint64) {
	for x >= 0x80 {
		bw.WriteByte_(byte(x) | 0x80)
		x >>= 7
	}
	bw.WriteByte_(byte(x))
}

func (bw *BufWriter) WriteStringLenUvarint(s string){
	bw.WriteUvarint(uint64(len(s)))
	bw.WriteString_(s)
}

func (w *BufWriter) GetBytes() []byte{
	return w.buf
}