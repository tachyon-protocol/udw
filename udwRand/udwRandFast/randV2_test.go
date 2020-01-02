package udwRandFast

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"io"
	"sync"
	"testing"
	"time"
)

var gFastRandReaderV2 fastRandReaderV2

func MustReadV2(buf []byte) {
	gFastRandReaderV2.mustRead(buf)
}

type fastRandReaderV2 struct {
	hash       udwCryptoSha3.ShakeHash
	hashStates [64]byte
	timeStates [8]byte
	lock       sync.Mutex
}

func (r *fastRandReaderV2) mustRead(dst []byte) {
	r.lock.Lock()
	r.init()
	r.hash.Reset()
	binary.BigEndian.PutUint64(r.timeStates[:], uint64(time.Now().UnixNano()))
	r.hash.Write(r.timeStates[:])
	r.hash.Write(r.hashStates[:])
	r.hash.Read(dst)
	r.hash.Read(r.hashStates[:])
	r.lock.Unlock()
	return
}

func (r *fastRandReaderV2) init() {
	if r.hash != nil {
		return
	}
	_, err := io.ReadFull(rand.Reader, r.hashStates[:])
	if err != nil {
		panic(err)
	}
	r.hash = udwCryptoSha3.NewShake256()
	r.hash.Write([]byte(time.Now().String()))
	r.hash.Write(r.hashStates[:])
	r.hash.Read(r.hashStates[:])

}

func TestMustReadV2(ot *testing.T) {
	DoTestRead(MustReadV2)
	DoTestSpeed(MustReadV2, "MustReadV2")
}
