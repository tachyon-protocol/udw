package udwGaocLib

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestGoBuffer(ot *testing.T) {
	buf := &GoBuffer{}
	buf.WriteString("1")
	buf.ResetToRead()
	out := buf.ReadString()
	udwTest.Equal(out, "1")
	buf.ResetToWrite()
	for i := 0; i < 100; i++ {
		buf.WriteString("1")
		buf.WriteInt(1)
		buf.WriteBool(i%2 == 0)
		buf.WriteFloat64(1010)
		buf.WriteByteSlice([]byte{1, 2, 3, 0})
	}
	buf.ResetToRead()
	for i := 0; i < 100; i++ {
		out = buf.ReadString()
		udwTest.Equal(out, "1")
		outI := buf.ReadInt()
		udwTest.Equal(outI, 1)
		outB := buf.ReadBool()
		udwTest.Equal(outB, i%2 == 0)
		outF := buf.ReadFloat64()
		udwTest.Equal(outF, float64(1010))
		outBs := buf.ReadByteSlice()
		udwTest.Equal(outBs, []byte{1, 2, 3, 0})
	}
	buf.FreeFromGo()
}

func TestGoBufferBench(ot *testing.T) {
	return
	buf := &GoBuffer{}
	udwTest.Benchmark(func() {
		num := int(1e6)
		udwTest.BenchmarkSetNum(num)
		for i := 0; i < num; i++ {
			buf.WriteString("1")
			buf.WriteInt(1)
		}
		buf.ResetToRead()
		for i := 0; i < num; i++ {
			out := buf.ReadString()
			if out != "1" {
				udwTest.Equal(out, "1")
			}
			outI := buf.ReadInt()
			if outI != 1 {
				udwTest.Equal(outI, 1)
			}
		}
	})
	buf.FreeFromGo()
}
