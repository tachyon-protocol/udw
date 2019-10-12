package udwBinary

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"runtime"
	"testing"
)

func TestReadByteSliceWithUint32LenNoAlloc(t *testing.T) {
	buf := make([]byte, 1024)
	out, err := ReadByteSliceWithUint32LenNoAlloc(bytes.NewReader([]byte{1, 0, 0, 0, 0x30}), buf)
	udwTest.Equal(err, nil)
	udwTest.Equal(out, []byte{0x30})
	buf2 := bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0x30})
	var before runtime.MemStats
	var after runtime.MemStats
	runtime.ReadMemStats(&before)
	out, err = ReadByteSliceWithUint32LenNoAlloc(buf2, buf)
	runtime.ReadMemStats(&after)
	udwTest.Equal(err, io.ErrUnexpectedEOF)
	udwTest.Ok(after.TotalAlloc-before.TotalAlloc < 1024*1024*1024)
}

func TestReadUint32NoAllocV2(t *testing.T) {
	out, err := ReadUint32NoAllocV2(bytes.NewBuffer([]byte{1, 0, 0, 0, 0x30}))
	udwTest.Equal(err, nil)
	udwTest.Equal(out, uint32(1))

	buf := bytes.NewBuffer([]byte{100, 0, 0, 0})
	bufReader := io.Reader(buf)
	out, err = ReadUint32NoAllocV2(bufReader)
	udwTest.Equal(err, nil)
	udwTest.Equal(out, uint32(100))
}

func TestWriteByteSliceWithUint32LenNoAlloc(t *testing.T) {
	buf := &bytes.Buffer{}
	content := []byte{1, 2, 3, 4, 5, 6}
	resultContent := []byte{6, 0, 0, 0, 1, 2, 3, 4, 5, 6}
	err := WriteByteSliceWithUint32LenNoAllocV2(buf, content)
	udwTest.Equal(err, nil)
	udwTest.Equal(buf.Bytes(), resultContent)

	tmpB := make([]byte, 16)
	buf.Reset()
	err = WriteByteSliceWithUint32LenNoAlloc(buf, content, tmpB)
	udwTest.Equal(err, nil)
	udwTest.Equal(buf.Bytes(), resultContent)

	tmpB2 := udwBytes.NewBufWriter(tmpB)
	tmpB2.Reset()
	buf.Reset()
	err = WriteByteSliceWithUint32LenNoAlloc(buf, content, tmpB)
	udwTest.Equal(err, nil)
	udwTest.Equal(buf.Bytes(), resultContent)

	const num = 1e8

	udwTest.Benchmark(func() {
		udwTest.BenchmarkSetNum(num)
		udwTest.BenchmarkSetBytePerRun(len(content))
		for i := 0; i < num; i++ {
			buf.Reset()
			WriteByteSliceWithUint32LenNoAllocV3(buf, content, tmpB2)
		}
	})
}
