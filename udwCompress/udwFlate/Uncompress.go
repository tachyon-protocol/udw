package udwFlate

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCompress/kkcflate"
	"io"
	"sync"
)

func TryUncompress(inData []byte) (outData []byte, err error) {
	if len(inData) == 0 {
		return nil, errors.New("[uncopressV2] len(inData)==0")
	}
	if inData[0] == 0 {
		return inData[1:], nil
	}

	tmpBuf := udwBytes.BufWriter{}
	return FlateUnCompressWithTmpBuf(inData[1:], &tmpBuf)
}

func TryUncompressWithTmpBuf(inData []byte, tmpBuf *udwBytes.BufWriter) (outData []byte, err error) {
	if len(inData) == 0 {
		return nil, errors.New("[uncopressV2] len(inData)==0")
	}
	if inData[0] == 0 {
		return inData[1:], nil
	}
	return FlateUnCompressWithTmpBuf(inData[1:], tmpBuf)
}

func MustFlateUnCompress(inb []byte) (outb []byte) {
	tmpBuf := udwBytes.BufWriter{}
	outb, err := FlateUnCompressWithTmpBuf(inb, &tmpBuf)
	if err != nil {
		panic(err)
	}
	return outb
}

func FlateUnCompress(inb []byte) (outb []byte, err error) {
	tmpBuf := udwBytes.BufWriter{}
	return FlateUnCompressWithTmpBuf(inb, &tmpBuf)
}

var gUnCompressPool sync.Pool

func FlateUnCompressWithTmpBuf(inb []byte, tmpBuf *udwBytes.BufWriter) (outb []byte, err error) {
	var entry *flateUncompressCacheEntry
	obj := gUnCompressPool.Get()
	if obj == nil {
		entry = &flateUncompressCacheEntry{}
	} else {
		entry = obj.(*flateUncompressCacheEntry)
	}
	entry.bufR.ResetWithBuffer(inb)
	if entry.zlibR == nil {
		entry.zlibR = kkcflate.NewReader(&entry.bufR)
	} else {
		err = entry.zlibR.(kkcflate.Resetter).Reset(&entry.bufR, nil)
		if err != nil {
			entry.bufR.ResetWithBuffer(nil)
			gUnCompressPool.Put(entry)
			return nil, err
		}
	}
	_, err = tmpBuf.ReadFrom(entry.zlibR)
	if err != nil {
		entry.bufR.ResetWithBuffer(nil)
		gUnCompressPool.Put(entry)
		return nil, err
	}
	err = entry.zlibR.Close()
	if err != nil {
		entry.bufR.ResetWithBuffer(nil)
		gUnCompressPool.Put(entry)
		return nil, err
	}
	outb = tmpBuf.GetBytes()
	entry.bufR.ResetWithBuffer(nil)
	gUnCompressPool.Put(entry)
	return
}

type flateUncompressCacheEntry struct {
	bufR  udwBytes.BufReader
	zlibR io.ReadCloser
}
