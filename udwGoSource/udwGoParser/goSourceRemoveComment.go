package udwGoParser

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
)

type noGrowBuf struct {
	buf []byte
	pos int
}

func (buf *noGrowBuf) WriteByte(b byte) {
	buf.buf[buf.pos] = b
	buf.pos++
}
func (buf *noGrowBuf) Write(b []byte) {
	copy(buf.buf[buf.pos:], b)
	buf.pos += len(b)
}
func GoSourceRemoveCommentWithByte(in []byte) (out []byte) {
	posFile := udwGoReader.NewPosFile("", in)
	return goSourceRemoveComment(in, posFile)
}
func goSourceRemoveComment(in []byte, filePos *udwGoReader.FilePos) (out []byte) {
	r := udwGoReader.NewReader(in, filePos)
	buf := &noGrowBuf{
		buf: make([]byte, len(in)),
		pos: 0,
	}
	for {
		if r.IsEof() {
			return buf.buf
		}
		b := r.ReadByte()
		switch b {
		case '/':
			r.UnreadByte()
			if r.IsMatchAfter(tokenSlashStar) {
				thisBuf := r.ReadUntilString(tokenStarSlash)

				commentReader := udwGoReader.NewReader(thisBuf, nil)
				for {
					if commentReader.IsEof() {
						break
					}
					b := commentReader.ReadByte()
					if b == '\n' {
						buf.WriteByte('\n')
					} else {
						buf.WriteByte(' ')
					}
				}
			} else if r.IsMatchAfter(tokenDoubleSlash) {
				thisBuf := r.ReadUntilByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
				buf.WriteByte('\n')
			} else {
				buf.WriteByte(r.ReadByte())
			}
		case '"', '`':
			r.UnreadByte()
			startPos := r.Pos()

			mustReadGoString(r)
			buf.Write(r.BufToCurrent(startPos))
		case '\'':
			r.UnreadByte()
			startPos := r.Pos()
			mustReadGoChar(r)
			buf.Write(r.BufToCurrent(startPos))
		default:
			buf.WriteByte(b)
		}

	}
}
