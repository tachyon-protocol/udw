package udwRandFast

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwCryptoSha3"
	"io"
	"sync"
	"time"
)

var gFastRandReaderV3 fastRandReaderV3

func MustRead(buf []byte) {
	gFastRandReaderV3.mustRead(buf)
}

func MustGetByteWithNum(n int) []byte {
	buf := make([]byte, n)
	MustRead(buf)
	return buf
}

const randBlockSizeV3 = 256
const addTimeNonceByteMaxValue = 1024
const reSeedNonceByteMaxValue = 1024 * 1024

type fastRandReaderV3 struct {
	stream           cipher.Stream
	buf              [randBlockSizeV3]byte
	lock             sync.Mutex
	addTimeNonceByte int
	reSeedNonceByte  int
	hasher           udwCryptoSha3.ShakeHash
	hashState        [64]byte
}

func (r *fastRandReaderV3) mustRead(dst []byte) {
	r.lock.Lock()
	r.init()
	pos := 0
	for {
		remainSize := len(dst) - pos
		if remainSize == 0 {
			break
		}
		thisBlockSize := randBlockSizeV3
		if remainSize < randBlockSizeV3 {
			thisBlockSize = remainSize
		}
		r.reSeedNonceByte -= thisBlockSize
		r.addTimeNonceByte -= thisBlockSize
		if r.addTimeNonceByte <= 0 {
			binary.BigEndian.PutUint64(r.hashState[:8], uint64(time.Now().UnixNano()))
			r.hasher.Write(r.hashState[:8])
			r.addTimeNonceByte = addTimeNonceByteMaxValue
			if r.reSeedNonceByte <= 0 {
				_, err := io.ReadFull(rand.Reader, r.hashState[:64])
				if err != nil {
					panic(err)
				}
				r.hasher.Write(r.hashState[:])
				r.hasher.Read(r.hashState[:])
				r.hasher.Reset()
				r.hasher.Write(r.hashState[:])
				block, err := AesCtr.NewCipher(r.hashState[:32])
				if err != nil {
					panic(err)
				}
				r.stream = cipher.NewCTR(block, r.hashState[32:48])
				r.reSeedNonceByte = reSeedNonceByteMaxValue
			}
		}
		r.stream.XORKeyStream(dst[pos:pos+thisBlockSize], r.buf[:thisBlockSize])
		pos += thisBlockSize
	}
	r.lock.Unlock()
	return
}

func (r *fastRandReaderV3) init() {
	if r.stream != nil {
		return
	}
	r.hasher = udwCryptoSha3.NewShake256()
	_, err := io.ReadFull(rand.Reader, r.hashState[:64])
	if err != nil {
		panic(err)
	}
	r.hasher.Write(r.hashState[:])
	binary.BigEndian.PutUint64(r.hashState[:8], uint64(time.Now().UnixNano()))
	r.hasher.Write(r.hashState[:8])
	r.hasher.Read(r.hashState[:])
	r.hasher.Reset()
	r.hasher.Write(r.hashState[:])
	block, err := AesCtr.NewCipher(r.hashState[:32])
	if err != nil {
		panic(err)
	}
	r.stream = cipher.NewCTR(block, r.hashState[32:48])
	_, err = io.ReadFull(rand.Reader, r.buf[:])
	if err != nil {
		panic(err)
	}
	r.addTimeNonceByte = addTimeNonceByteMaxValue
	r.reSeedNonceByte = reSeedNonceByteMaxValue
}
