package udwRandFast

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwCompress/udwFlate"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"sync"
	"testing"
)

func MustReadGoCrypto(buf []byte) {
	n, err := rand.Read(buf)
	if err != nil {
		panic("[MustReadGoCrypto] " + err.Error())
	}
	if n != len(buf) {
		panic("[MustReadGoCrypto] n!=len(buf)")
	}
}

var gFastRandReader fastRandReader

func MustReadV1(buf []byte) {
	gFastRandReader.mustRead(buf)
}

const randBlockSize = 256

type fastRandReader struct {
	stream  cipher.Stream
	buf     [randBlockSize]byte
	lock    sync.Mutex
	hasInit bool
}

func (r *fastRandReader) mustRead(dst []byte) {
	r.lock.Lock()
	r.init()

	remainSize := len(dst)
	for {
		if remainSize >= randBlockSize {
			r.stream.XORKeyStream(dst[remainSize-randBlockSize:remainSize], r.buf[:])
			remainSize -= randBlockSize
			continue
		}
		r.stream.XORKeyStream(dst[0:remainSize], r.buf[:remainSize])
		break
	}
	r.lock.Unlock()
	return
}

func (r *fastRandReader) init() {
	if r.stream != nil {
		return
	}
	_, err := io.ReadFull(rand.Reader, r.buf[:48])
	if err != nil {
		panic(err)
	}
	block, err := AesCtr.NewCipher(r.buf[:32])
	if err != nil {
		panic(err)
	}
	r.stream = cipher.NewCTR(block, r.buf[32:48])
	_, err = io.ReadFull(rand.Reader, r.buf[:])
	if err != nil {
		panic(err)
	}
}

func DoTestRead(reader func(buf []byte)) {
	seeBufMap := map[[16]byte]struct{}{}
	for i := 0; i < 10; i++ {
		buf := [16]byte{}
		reader(buf[:])
		_, ok := seeBufMap[buf]
		if ok {
			panic("[DoTestRead] rand is not enough 1 " + hex.Dump(buf[:]))
		}
		seeBufMap[buf] = struct{}{}

	}
	for i := 0; i < 10; i++ {

		const totalNum = 256 * 1024 * 4
		buf := make([]byte, totalNum)
		reader(buf)
		byteNumMap := [256]int{}
		for j := 0; j < len(buf); j++ {
			byteNumMap[buf[j]]++
		}
		for j := 0; j < 256; j++ {
			seeRate := float64(byteNumMap[j]) * 256 / totalNum
			if seeRate > 1.1 || seeRate < 0.9 {
				panic("[DoTestRead] rand is not enough 2 " + udwStrconv.FormatFloatPercentPrec4(seeRate))
			}
		}
	}
	for i := 0; i < 10; i++ {
		for _, size := range []int{
			64,
			1024,
			1024 * 4,
			1024 * 32,
			1024 * 256,
		} {

			buf := make([]byte, size)
			reader(buf)
			outB := udwFlate.FlateMustCompress(buf)
			if len(outB) < len(buf) {
				panic("[DoTestRead] rand is not enough 3")
			}
		}
	}
}

func DoTestSpeed(reader func(buf []byte), name string) {
	buf := make([]byte, 1024)
	reader(buf)
	thisName := name + " 1k"
	udwTest.Benchmark(func() {
		const num = 1e4
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(1024)
		udwTest.BenchmarkSetName(thisName)
		for i := 0; i < num; i++ {
			reader(buf)
		}
	})
	buf = make([]byte, 16)
	thisName = name + " 16"
	udwTest.Benchmark(func() {
		const num = 1e5
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(16)
		udwTest.BenchmarkSetName(thisName)
		for i := 0; i < num; i++ {
			reader(buf)
		}
	})
}

func TestMustReadV1(ot *testing.T) {
	DoTestRead(MustReadGoCrypto)
	DoTestRead(MustReadV1)
	DoTestSpeed(MustReadGoCrypto, "MustReadGoCrypto")
	DoTestSpeed(MustReadV1, "MustReadV1")
}
