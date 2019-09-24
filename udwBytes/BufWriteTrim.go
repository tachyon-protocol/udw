package udwBytes

func (w *BufWriter) TrimSuffixSpace() {
	if len(w.buf) == 0 {
		return
	}
	endPos := -1
	for i := len(w.buf) - 1; i >= 0; i-- {
		b := w.buf[i]
		isSpace := (b == ' ' || b == '\t' || b == '\n')
		if isSpace == false {
			endPos = i
			break
		}
	}
	w.buf = w.buf[:endPos+1]
}

func (w *BufWriter) TrimSuffixOneByte(b byte) {
	b1, ok := w.GetLastByte()
	if ok && b1 == b {
		w.buf = w.buf[:len(w.buf)-1]
	}
}
