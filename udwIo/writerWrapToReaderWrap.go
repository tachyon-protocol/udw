package udwIo

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
)

func NewWriterWrapToReaderWrap(newWriter func(w1 io.Writer) (w2 io.WriteCloser), r1 io.Reader) (r2 io.Reader) {
	w := &writerWrapToReaderWrap{
		r1:    r1,
		r1Buf: make([]byte, 32*1024),
	}
	w2 := newWriter(&w.tmpWriter)
	w.tmpBufStartPos = 0
	w.tmpBufEndPos = w.tmpWriter.GetLen()
	w.w2 = w2
	return w
}

type writerWrapToReaderWrap struct {
	r1             io.Reader
	w2             io.WriteCloser
	tmpBufStartPos int
	tmpBufEndPos   int
	r1Buf          []byte
	tmpWriter      udwBytes.BufWriter
	isREof         bool
}

func (r *writerWrapToReaderWrap) reload() (err error) {
	r.tmpWriter.Reset()
	n, err := r.r1.Read(r.r1Buf)
	if err != nil {
		if err == io.EOF {
			r.isREof = true
		} else {
			return err
		}
	}
	_, err = r.w2.Write(r.r1Buf[:n])
	if err != nil {
		return err
	}
	if r.isREof {
		err = r.w2.Close()
		if err != nil {
			return err
		}
	}
	r.tmpBufStartPos = 0
	r.tmpBufEndPos = r.tmpWriter.GetPos()
	return nil
}

func (r *writerWrapToReaderWrap) Read(b []byte) (n int, err error) {
	for {
		if n >= len(b) {
			break
		}
		bufL := r.tmpBufEndPos - r.tmpBufStartPos
		if bufL == 0 {
			if r.isREof {
				return n, io.EOF
			}
			err = r.reload()
			if err != nil {
				return n, err
			}
			continue
		} else if len(b)-n >= bufL {
			copy(b[n:], r.tmpWriter.GetBytes()[r.tmpBufStartPos:r.tmpBufEndPos])
			n += bufL
			r.tmpBufStartPos = r.tmpBufEndPos
		} else {
			endPos := r.tmpBufStartPos + len(b) - n
			copy(b[n:], r.tmpWriter.GetBytes()[r.tmpBufStartPos:endPos])
			r.tmpBufStartPos = endPos
			n = len(b)
			break
		}
	}
	return n, err
}
