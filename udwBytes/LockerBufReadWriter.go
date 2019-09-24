package udwBytes

import "sync"

type LockerBufReadWriter struct {
	locker         sync.Mutex
	buf            []byte
	bufStartOffset int
	readPos        int
	writePos       int
}

func (buf *LockerBufReadWriter) Write(p []byte) (n int, err error) {
	buf.locker.Lock()
	buf.tryGrow__NOLOCK(len(p))
	copy(buf.buf[buf.writePos-buf.bufStartOffset:], p)
	buf.writePos += len(p)
	buf.locker.Unlock()
	return len(p), nil
}

func (buf *LockerBufReadWriter) tryGrow__NOLOCK(toWrite int) {
	needSize := buf.writePos - buf.bufStartOffset + toWrite
	if needSize > len(buf.buf) {
		newBuf := make([]byte, len(buf.buf)*2+toWrite)
		copy(newBuf, buf.buf[buf.readPos-buf.bufStartOffset:buf.writePos-buf.bufStartOffset])
		buf.buf = newBuf
		buf.bufStartOffset = buf.readPos
	}
}

type ReadCopyFromPosResponse struct {
	ReadPos int
	Content []byte
}

func (buf *LockerBufReadWriter) TryReadCopyFromPos(readPos int) (resp ReadCopyFromPosResponse) {
	buf.locker.Lock()
	resp.ReadPos = readPos
	if resp.ReadPos < buf.readPos {
		resp.ReadPos = buf.readPos
	}
	if resp.ReadPos > buf.writePos {
		resp.ReadPos = buf.writePos
	}
	resp.Content = make([]byte, buf.writePos-resp.ReadPos)
	copy(resp.Content, buf.buf[resp.ReadPos-buf.bufStartOffset:buf.writePos-buf.bufStartOffset])
	buf.locker.Unlock()
	return resp
}

func (buf *LockerBufReadWriter) SetReadPos(toCutReadPos int) (errMsg string) {
	buf.locker.Lock()
	if buf.readPos > toCutReadPos {
		buf.locker.Unlock()
		return "chferyzhm4"
	}
	if buf.writePos < toCutReadPos {
		buf.locker.Unlock()
		return "7j5zp8gp4b"
	}
	buf.readPos = toCutReadPos
	buf.locker.Unlock()
	return ""
}
