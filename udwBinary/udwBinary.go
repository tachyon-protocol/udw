package udwBinary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
	"math"
)

func MustWriteString255(buf *bytes.Buffer, s string) {
	if len(s) > 255 {
		panic(fmt.Errorf("[MustWriteString255] len(s)[%d]>255", len(s)))
	}
	buf.WriteByte(byte(len(s)))
	buf.WriteString(s)
}

func WriteUint16(buf *bytes.Buffer, i uint16) {
	buf.WriteByte(byte(i))
	buf.WriteByte(byte(i >> 8))
}

func WriteUint32(buf *bytes.Buffer, i uint32) {
	buf.WriteByte(byte(i))
	buf.WriteByte(byte(i >> 8))
	buf.WriteByte(byte(i >> 16))
	buf.WriteByte(byte(i >> 24))
}

func WriteUint64(buf *bytes.Buffer, i uint64) {

	buf.WriteByte(byte(i))
	buf.WriteByte(byte(i >> 8))
	buf.WriteByte(byte(i >> 16))
	buf.WriteByte(byte(i >> 24))
	buf.WriteByte(byte(i >> 32))
	buf.WriteByte(byte(i >> 40))
	buf.WriteByte(byte(i >> 48))
	buf.WriteByte(byte(i >> 56))
}

func WriteBool(buf *bytes.Buffer, b bool) {
	b1 := byte(0)
	if b == true {
		b1 = 1
	}
	buf.WriteByte(b1)
}
func ReadBoolWithByte(b byte) bool {
	return b == 1
}

func WriteFloat64(buf *bytes.Buffer, f float64) {
	WriteUint64(buf, math.Float64bits(f))
}
func ReadFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b))
}

func ReadString255(r io.Reader) (s string, err error) {
	buf := make([]byte, 256)
	_, err = io.ReadFull(r, buf[:1])
	if err != nil {
		return "", err
	}
	len := int(buf[0])
	_, err = io.ReadFull(r, buf[:len])
	if err != nil {
		return "", err
	}
	return string(buf[:len]), nil
}

func ReadString255WithByteSlice(b []byte) (s string, err error) {
	if len(b) == 0 {
		return "", io.ErrShortBuffer
	}
	thisLen := int(b[0])
	if len(b) < thisLen+1 {
		return "", io.ErrShortBuffer
	}
	return string(b[1 : thisLen+1]), nil
}

const readByteSliceWithUint32LenChunkSize = 32 * 1024

func ReadByteSliceWithUint32Len(r io.Reader) (b []byte, err error) {
	bufW := udwBytes.BufWriter{}
	err = ReadByteSliceWithUint32LenToBufW(r, &bufW)
	if err != nil {
		return nil, err
	}
	return bufW.GetBytes(), nil

}

func ReadByteSliceWithUint32LenToBufW(r io.Reader, bufW *udwBytes.BufWriter) (err error) {
	sizeBuf := bufW.GetHeadBuffer(4)
	rn, err := io.ReadFull(r, sizeBuf)
	if err != nil {
		if err == io.ErrUnexpectedEOF && rn == 0 {
			return io.EOF
		}
		return errors.New("24fm8ux7ba " + err.Error())
	}
	size := binary.LittleEndian.Uint32(sizeBuf)
	if size == 0 {
		return nil
	}
	if size <= readByteSliceWithUint32LenChunkSize {
		data := bufW.GetHeadBuffer(int(size))
		_, err = io.ReadFull(r, data)
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return errors.New("a4e7ju2ase " + err.Error())
		}
		bufW.AddPos(int(size))
		return nil
	}

	readedPos := 0

	for {
		tmpBufUsedSize := int(size) - readedPos
		if tmpBufUsedSize > readByteSliceWithUint32LenChunkSize {
			tmpBufUsedSize = readByteSliceWithUint32LenChunkSize
		}
		tmpBuf := bufW.GetHeadBuffer(tmpBufUsedSize)
		thisNs, err := r.Read(tmpBuf)
		readedPos += thisNs
		bufW.AddPos(thisNs)
		if readedPos == int(size) {
			return nil
		}
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
	}
}

func WriteByteSliceWithUint32Len(w io.Writer, b []byte) (err error) {
	var buf bytes.Buffer
	WriteUint32(&buf, uint32(len(b)))
	buf.Write(b)
	_, err = w.Write(buf.Bytes())
	return err
}
