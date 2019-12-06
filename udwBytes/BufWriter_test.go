package udwBytes

import (
	"bytes"
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
	"testing"
)

func TestBufWriter(ot *testing.T) {
	bufW := BufWriter{}
	bufW.WriteByte(1)
	udwTest.Equal(bufW.GetPos(), 1)
	udwTest.Equal(bufW.GetBytes(), []byte{1})
	bufW.AddPos(1)
	udwTest.Equal(bufW.GetPos(), 2)
	udwTest.Equal(bufW.GetBytes(), []byte{1, 0})
	bufW.Reset()
	udwTest.Equal(bufW.GetPos(), 0)
	bufW.Write([]byte{2, 3, 4})
	udwTest.Equal(bufW.GetBytes(), []byte{2, 3, 4})
	bufW.Reset()
	inB := bytes.Repeat([]byte{1}, 1024)
	bufW.Write(inB)
	udwTest.Equal(bufW.GetBytes(), inB)
	bufW.Write(inB)
	udwTest.Equal(bufW.GetBytes(), append(inB, inB...))
	bufW.Reset()
	bufW.WriteString("123")
	udwTest.Equal(bufW.GetBytes(), []byte("123"))
	bufW.Reset()
	n, err := bufW.ReadFrom(NewBufReader(inB))
	udwTest.Equal(n, int64(len(inB)))
	udwTest.Equal(err, nil)
	udwTest.Equal(bufW.GetBytes(), inB)
	bufW.Reset()
	bufW.WriteUvarint(1)
	udwTest.Equal(bufW.GetBytes(), []byte{1})
	bufW.Reset()
	bufW.WriteVarint(2)
	udwTest.Equal(bufW.GetBytes(), []byte{4})

	bufW.Reset()
	bufW.WriteLittleEndUint32(1)
	udwTest.Equal(bufW.GetBytes(), []byte{1, 0, 0, 0})

	bufW.Reset()
	bufW.WriteBigEndUint32(1)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 0, 0, 1})
	b1 := bufW.GetBytesClone()
	bufW.Reset()
	bufW.WriteBigEndUint32(256)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 0, 1, 0})
	udwTest.Ok(string(b1) < bufW.GetString())
	udwTest.Equal(b1, []byte{0, 0, 0, 1})

	bufW.Reset()
	bufW.MustWriteString255("123")
	udwTest.Equal(bufW.GetBytes(), []byte{3, '1', '2', '3'})

	bufW.Reset()
	errMsg := udwErr.PanicToErrorMsg(func() {
		bufW.MustWriteString255(strings.Repeat("1", 256))
	})
	udwTest.Ok(strings.Contains(errMsg, "[MustWriteString255]"))

	bufW.Reset()
	bufW.WriteBool(true)
	udwTest.Equal(bufW.GetBytes(), []byte{1})

	bufW.Reset()
	buf := bufW.GetHeadBuffer(1024)
	udwTest.Equal(len(buf), 1024)

	bufW2 := NewBufWriter([]byte{1, 2, 3})
	udwTest.Equal(bufW2.GetLen(), 3)
	bufW.WriteByte(1)

	bufW2 = ResetBufWriter(bufW2)
	udwTest.Equal(bufW2.GetLen(), 0)

	bufW2 = ResetBufWriter(nil)
	udwTest.Ok(bufW2 != nil)
	udwTest.Equal(bufW2.GetLen(), 0)

	bufW2.ResetWithBuffer(nil)
	udwTest.Equal(bufW2.GetLen(), 0)
	bufW2.WriteByte(1)
	udwTest.Equal(bufW2.GetBytes(), []byte{1})
	bufW2.ResetWithBuffer([]byte{1, 2, 3})
	udwTest.Equal(bufW2.GetLen(), 0)
	bufW2.SetPos(1024)
	udwTest.Equal(bufW2.GetLen(), 1024)
	udwTest.Equal(len(bufW2.GetBytes()), 1024)

	bufW = BufWriter{}
	bufW.Reset()
	bufW.WriteLittleEndFloat64(0)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 0, 0, 0, 0, 0, 0, 0})

	bufW.Reset()
	bufW.WriteBigEndUint16(256)
	udwTest.Equal(bufW.GetBytes(), []byte{1, 0})

	bufW.Reset()
	bufW.WriteLittleEndUint16(256)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 1})

	bufW.Reset()
	bufW.WriteBigEndUint64(256)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 0, 0, 0, 0, 0, 1, 0})

	bufW.Reset()
	bufW.WriteLittleEndUint64(256)
	udwTest.Equal(bufW.GetBytes(), []byte{0, 1, 0, 0, 0, 0, 0, 0})

	bufW.Reset()
	n, err = bufW.ReadFrom(&errorReader{})
	udwTest.Equal(n, int64(0))
	udwTest.Ok(err != nil)
	udwTest.Equal(err.Error(), "ur4twn348w")

	{
		bufW := BufWriter{}
		bufW.WriteByte(1)
		bufW.SetPos(1024)
		bufW.WriteByte(1)
		udwTest.Equal(bufW.GetBytes()[0], byte(1))
		udwTest.Equal(bufW.GetBytes()[1024], byte(1))
		bufW.SetPos(1)
		bufW.WriteByte(1)
		udwTest.Equal(bufW.GetBytes()[1], byte(1))
		bufW.SetPos(10240)
		bufW.WriteByte(1)
		udwTest.Equal(bufW.GetBytes()[1024], byte(1), 1024)
		udwTest.Equal(bufW.GetBytes()[10240], byte(1))
	}
}

type errorReader struct {
}

func (r *errorReader) Read(buf []byte) (n int, err error) {
	return 0, errors.New("ur4twn348w")
}

func TestBufWriter_1(ot *testing.T) {
	const num = 10
	udwTest.BenchmarkSetName("string add")
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		s := ""
		for i := 0; i < num; i++ {
			s += "1"
		}
		buf := []byte(s)
		_ = buf
	})
	udwTest.BenchmarkSetName("BufWriter add")
	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		bufW := BufWriter{}
		for i := 0; i < num; i++ {
			bufW.WriteString_("1")
		}
		buf := bufW.GetBytes()
		_ = buf
	})
}
