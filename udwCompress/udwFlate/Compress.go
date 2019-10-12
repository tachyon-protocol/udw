package udwFlate

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCompress/kkcflate"
	"sync"
)

const MinUsefullInputSize = 33

func TryCompress(inData []byte) (outData []byte) {
	bufW := udwBytes.BufWriter{}
	TryCompressToBufW(inData, &bufW)
	return bufW.GetBytes()
}

func TryCompressToBufW(inData []byte, bufW *udwBytes.BufWriter) {
	if len(inData) < MinUsefullInputSize {

		bufW.WriteByte(0)
		bufW.Write(inData)
		return
	}

	bufW.WriteByte(1)
	startPos := bufW.GetPos()
	MustFlateCompressWithBufferToBufW(inData, bufW)
	compressedSize := bufW.GetPos() - startPos
	if compressedSize < len(inData) {
		return
	}
	bufW.SetPos(startPos - 1)
	bufW.WriteByte(0)
	bufW.Write(inData)
	return
}

func FlateMustCompress(inb []byte) (outb []byte) {
	bufW := udwBytes.BufWriter{}
	MustFlateCompressWithBufferToBufW(inb, &bufW)
	return bufW.GetBytes()
}

var gCompressPool sync.Pool

func MustFlateCompressWithBufferToBufW(inb []byte, bufW *udwBytes.BufWriter) {
	var flateW *kkcflate.Writer
	obj := gCompressPool.Get()
	if obj == nil {
		var err error
		flateW, err = kkcflate.NewWriter(bufW, 4)
		if err != nil {
			panic(err)
		}
	} else {
		flateW = obj.(*kkcflate.Writer)
		flateW.Reset(bufW)
	}
	_, err := flateW.Write(inb)
	if err != nil {
		flateW.Close()
		flateW.Reset(nil)
		gCompressPool.Put(flateW)
		panic(err)
	}
	err = flateW.Close()
	flateW.Reset(nil)
	gCompressPool.Put(flateW)
	if err != nil {
		panic(err)
	}
}
