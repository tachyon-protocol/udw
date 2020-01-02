package udwCryptoEncryptV3

import (
	"bytes"
	"crypto/cipher"
	"github.com/tachyon-protocol/udw/udwIo"
	"github.com/tachyon-protocol/udw/udwNet/udwNetTestV2"
	"github.com/tachyon-protocol/udw/udwRand"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"sync"
	"testing"
	"time"
)

func TestSymmetryConn(t *testing.T) {
	buf := &bytes.Buffer{}
	oconn := udwIo.StructWriterReaderCloser{
		Writer: buf,
		Reader: buf,
		Closer: udwIo.NopCloser,
	}
	key := Get32PskFromString("a1")
	conn1 := MustSymmetryConn(oconn, key)
	conn2 := MustSymmetryConn(oconn, key)
	content := bytes.Repeat([]byte{0}, 1024)
	_, err := conn1.Write(content)
	if err != nil {
		t.Error(err)
	}
	bufByte := make([]byte, len(content))
	_, err = conn2.Read(bufByte)
	if err != nil {
		t.Error(err)
	}
	_, err = conn2.Write(content)
	if err != nil {
		t.Error(err)
	}
	_, err = conn1.Read(bufByte)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(bufByte, content) {
		t.Error("!bytes.Equal(buf,content)")
	}
	conn1.Close()
	conn2.Close()

	conn1.Close()
	conn2.Close()

	sleepRandTimeFn := func() {
		time.Sleep(udwRand.TimeDurationBetween(0, 3*time.Millisecond))
	}

	wg := sync.WaitGroup{}
	block := GetAesBlockFrom32Psk(key)
	for i := 0; i < 1e2; i++ {
		wg.Add(1)
		go func() {
			conns := udwNetTestV2.NewPipeConns(4)
			conn1 := PoolGetSymmetryConnWithBlock(conns.Conn1(), block)
			conn2 := PoolGetSymmetryConnWithBlock(conns.Conn2(), block)
			bufByte := make([]byte, len(content))
			wg2 := sync.WaitGroup{}
			wg2.Add(4)
			go func() {
				sleepRandTimeFn()
				conn1.Read(bufByte)
				sleepRandTimeFn()
				conn1.Read(bufByte)
				wg2.Done()
			}()
			go func() {
				sleepRandTimeFn()
				conn2.Write(content)
				sleepRandTimeFn()
				conn2.Write(content)
				wg2.Done()
			}()
			go func() {
				sleepRandTimeFn()
				conn1.Close()
				wg2.Done()
			}()
			go func() {
				sleepRandTimeFn()
				conn2.Close()
				wg2.Done()
			}()
			wg2.Wait()
			PoolPutSymmetryConnAndClose(conn1)
			PoolPutSymmetryConnAndClose(conn2)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestSymmetryConnBenchmark(t *testing.T) {
	const benchNum = 1000
	const perRunSize = 511
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(benchNum)
		udwTest.BenchmarkSetBytePerRun(perRunSize * 2)
		buf := &bytes.Buffer{}
		oconn := udwIo.StructWriterReaderCloser{
			Writer: buf,
			Reader: buf,
			Closer: udwIo.NopCloser,
		}
		key := Get32PskFromString("a1")
		conn1 := MustSymmetryConn(oconn, key)
		conn2 := MustSymmetryConn(oconn, key)
		content := make([]byte, perRunSize)
		bufByte := make([]byte, len(content))
		for i := 0; i < benchNum; i++ {
			conn1.Write(content)
			conn2.Read(bufByte)
		}
	})
}

func TestSymmetryNewConn(t *testing.T) {
	const benchNum = 1e3
	const perRunSize = 511

	buf := &bytes.Buffer{}
	oconn := udwIo.StructWriterReaderCloser{
		Writer: buf,
		Reader: buf,
		Closer: udwIo.NopCloser,
	}
	oconn2 := (io.ReadWriteCloser)(oconn)
	key := Get32PskFromString("a1")
	block := GetAesBlockFrom32Psk(key)
	content := make([]byte, perRunSize)
	bufByte := make([]byte, len(content))
	runOnce := func() {
		conn1 := PoolGetSymmetryConnWithBlock(oconn2, block)
		conn2 := PoolGetSymmetryConnWithBlock(oconn2, block)
		conn1.Write(content)
		conn2.Read(bufByte)
		PoolPutSymmetryConnAndClose(conn1)
		PoolPutSymmetryConnAndClose(conn2)
	}
	runOnce()
	runOnce()

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(benchNum)
		udwTest.BenchmarkSetBytePerRun(perRunSize * 2)
		for i := 0; i < benchNum; i++ {
			runOnce()
		}
	})
}

func BenchmarkSymmetryConn(ot *testing.B) {
	ot.StopTimer()
	buf := &bytes.Buffer{}
	oconn := udwIo.StructWriterReaderCloser{
		Writer: buf,
		Reader: buf,
		Closer: udwIo.NopCloser,
	}
	key := Get32PskFromString("a1")
	conn1 := MustSymmetryConn(oconn, key)
	conn2 := MustSymmetryConn(oconn, key)
	ot.SetBytes(int64(4096 * 2))
	content := make([]byte, 4096)
	bufByte := make([]byte, len(content))
	ot.StartTimer()
	for i := 0; i < ot.N; i++ {
		conn1.Write(content)
		conn2.Read(bufByte)
	}
}

type NopCipherBlock struct {
}

func (b NopCipherBlock) BlockSize() int {
	return 16
}

func (b NopCipherBlock) Encrypt(dst, src []byte) {
	copy(dst, src)
}

func (b NopCipherBlock) Decrypt(dst, src []byte) {
	copy(dst, src)
}

func NopCipherConn(conn io.ReadWriteCloser) (outConn io.ReadWriteCloser, err error) {
	return &symmetryConn{
		rwc:   conn,
		block: NopCipherBlock{},
	}, nil
}

func BenchmarkSymmetryConn1(ot *testing.B) {
	ot.StopTimer()
	buf := &bytes.Buffer{}
	oconn := udwIo.StructWriterReaderCloser{
		Writer: buf,
		Reader: buf,
		Closer: udwIo.NopCloser,
	}

	conn1, _ := NopCipherConn(oconn)
	conn2, _ := NopCipherConn(oconn)
	ot.SetBytes(int64(4096 * 2))
	content := make([]byte, 4096)
	bufByte := make([]byte, len(content))
	ot.StartTimer()
	for i := 0; i < ot.N; i++ {
		conn1.Write(content)
		conn2.Read(bufByte)
	}
}

func BenchmarkSymmetryConn2(ot *testing.B) {
	ot.StopTimer()
	ot.SetBytes(int64(4096))
	content := make([]byte, 4096)
	ctr := cipher.NewCTR(NopCipherBlock{}, content[:16])
	ot.StartTimer()
	for i := 0; i < ot.N; i++ {
		ctr.XORKeyStream(content, content)
	}
}

func BenchmarkCopy(ot *testing.B) {
	ot.StopTimer()
	ot.SetBytes(int64(4096))
	content := make([]byte, 4096)
	bufByte := make([]byte, len(content))
	ot.StartTimer()
	for i := 0; i < ot.N; i++ {
		copy(bufByte, content)
	}
}

func BenchmarkMake(ot *testing.B) {
	ot.StopTimer()
	ot.SetBytes(int64(4096))
	ot.StartTimer()
	for i := 0; i < ot.N; i++ {
		a := make([]byte, 4096)
		a[1] = 1
		_ = a
	}
}
