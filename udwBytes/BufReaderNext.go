package udwBytes

import "bytes"

func (r *BufReader) ReadToLineEnd() []byte {
	startPos := r.pos
	for {
		if r.pos >= len(r.buf) {
			break
		}
		b := r.buf[r.pos]
		r.pos++
		if b == '\n' {
			break
		}
	}
	return r.buf[startPos:r.pos]
}

func (r *BufReader) ReadUtilFlag(flag []byte) []byte {
	startPos := r.pos
	if r.pos >= len(r.buf) {
		return nil
	}
	pos := bytes.Index(r.buf[r.pos:], flag)
	if pos == -1 {
		r.pos = len(r.buf)
		return r.buf[startPos:r.pos]
	}
	r.pos += pos + len(flag)
	return r.buf[startPos:r.pos]
}

func (r *BufReader) IsMatchPrefix(toMatch []byte) bool {
	if r.GetRemainSize() < len(toMatch) {
		return false
	}
	return bytes.Equal(r.buf[r.pos:r.pos+len(toMatch)], toMatch)
}

func (r *BufReader) GetNextLine() []byte {
	thisLine := r.ReadToLineEnd()
	r.pos = r.pos - len(thisLine)
	return thisLine
}

func (r *BufReader) IsNextLineContains(toMatch []byte) bool {
	nextLine := r.GetNextLine()
	return bytes.Contains(nextLine, toMatch)
}

func (r *BufReader) ReadSpace() []byte {
	startPos := r.pos
	for {
		if r.pos >= len(r.buf) {
			break
		}
		b := r.buf[r.pos]
		isSpace := b == ' ' || b == '\t' || b == '\n'
		if isSpace == false {
			break
		}
		r.pos++
	}
	return r.buf[startPos:r.pos]
}

func (r *BufReader) ReadToEof() []byte {
	startPos := r.pos
	r.pos = len(r.buf)
	return r.buf[startPos:]
}
