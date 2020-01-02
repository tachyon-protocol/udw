package udwGoParser

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoReader"
)

func OcSourceRemoveCommentWithByte(in []byte) (out []byte) {
	posFile := udwGoReader.NewPosFile("", in)
	return ocSourceRemoveComment(in, posFile)
}
func ocSourceRemoveComment(in []byte, filePos *udwGoReader.FilePos) (out []byte) {
	r := udwGoReader.NewReader(in, filePos)
	buf := &noGrowBuf{
		buf: make([]byte, len(in)),
		pos: 0,
	}
	currentLine := 0
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
			} else if r.IsMatchAfter(tokenHeadBuildFlag) && currentLine == 0 {
				buf.WriteByte(r.ReadByte())
			} else if r.IsMatchAfter(tokenDoubleSlash) {
				thisBuf := r.ReadUntilByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
				buf.WriteByte('\n')
			} else {
				buf.WriteByte(r.ReadByte())
			}
		case '#':
			r.UnreadByte()
			if r.IsMatchAfter(tokenOcPragmaMark) {
				thisBuf := r.ReadUntilByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
				buf.WriteByte('\n')
			} else {
				buf.WriteByte(r.ReadByte())
			}
		case '\n':
			currentLine++
			buf.WriteByte(b)
		default:
			buf.WriteByte(b)
		}

	}
}
