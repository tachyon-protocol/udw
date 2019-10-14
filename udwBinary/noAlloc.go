// +build !js

package udwBinary

import (
	"encoding/binary"
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

func WriteByteSliceWithUint32LenNoAlloc(w io.Writer, b []byte, tmpB []byte) (err error) {
	size := uint32(len(b))
	afterSize := 4 + len(b)
	if len(tmpB) < afterSize {
		tmpB = make([]byte, afterSize)
	}
	tmpB[0] = byte(size)
	tmpB[1] = byte(size >> 8)
	tmpB[2] = byte(size >> 16)
	tmpB[3] = byte(size >> 24)
	copy(tmpB[4:afterSize], b)
	_, err = w.Write(tmpB[:afterSize])
	return err
}

func WriteByteSliceWithUint32LenNoAllocV2(w io.Writer, b []byte) (err error) {
	const tmpStackBufSize = 4

	var tmpStackBuf [tmpStackBufSize]byte
	var tmpB []byte
	tmpStackBufp := uintptr(unsafe.Pointer(&tmpStackBuf))
	tmpStackBufp2 := *(*uintptr)(unsafe.Pointer(&tmpStackBufp))
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&tmpB))
	bx.Data = uintptr(tmpStackBufp2)
	bx.Len = tmpStackBufSize
	bx.Cap = tmpStackBufSize

	size := uint32(len(b))
	tmpB[0] = byte(size)
	tmpB[1] = byte(size >> 8)
	tmpB[2] = byte(size >> 16)
	tmpB[3] = byte(size >> 24)
	_, err = w.Write(tmpB[:4])
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func WriteByteSliceWithUint32LenNoAllocV3(w io.Writer, b []byte, tmpB *udwBytes.BufWriter) (err error) {
	tmpB.Reset()

	tmpB.WriteLittleEndUint32(uint32(len(b)))
	tmpB.Write(b)
	_, err = w.Write(tmpB.GetBytes())
	return err
}

func ReadByteSliceWithUint32LenNoAllocLimitMaxSize(r io.Reader, tmpB []byte, maxSize uint32) (outB []byte, err error) {
	size, err := ReadUint32NoAlloc(r, tmpB)
	if err != nil {
		return nil, err
	}
	if size < 0 || size > maxSize {
		return nil, errors.New("size < 0 || size > maxSize, size=[" + strconv.Itoa(int(size)) + "]")
	}
	outB, err = readFromReaderNoAllocWithSize(r, tmpB, int(size))
	if err != nil {
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return outB, nil
}

func ReadByteSliceWithUint32LenNoAlloc(r io.Reader, tmpB []byte) (outB []byte, err error) {
	size, err := ReadUint32NoAlloc(r, tmpB)
	if err != nil {
		return nil, err
	}
	outB, err = readFromReaderNoAllocWithSize(r, tmpB, int(size))
	if err != nil {
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return outB, nil
}

func ReadUint32NoAlloc(r io.Reader, tmpB []byte) (out uint32, err error) {
	if len(tmpB) < 4 {
		tmpB = make([]byte, 4)
	}
	_, err = io.ReadFull(r, tmpB[:4])
	if err != nil {
		return 0, err
	}
	out = binary.LittleEndian.Uint32(tmpB[:4])
	return out, nil
}

func ReadUint32NoAllocV2(r io.Reader) (out uint32, err error) {
	var tmpB [4]byte
	_, err = io.ReadFull(r, tmpB[:])
	if err != nil {
		return 0, err
	}
	out = binary.LittleEndian.Uint32(tmpB[:])
	return out, nil
}

func readFromReaderNoAllocWithSize(r io.Reader, tmpB []byte, size int) (out []byte, err error) {
	if len(tmpB) >= size {
		_, err = io.ReadFull(r, tmpB[:size])
		if err != nil {
			return nil, err
		}
		return tmpB[:size], nil
	}
	_, err = io.ReadFull(r, tmpB)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, size)
	copy(buf, tmpB)
	_, err = io.ReadFull(r, buf[len(tmpB):size])
	if err != nil {
		return nil, err
	}
	return buf[:size], nil
}
