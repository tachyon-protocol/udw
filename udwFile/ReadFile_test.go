package udwFile_test

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestReadFileToBufW(t *testing.T) {
	udwFile.MustDelete("zzzig_testFile")
	defer udwFile.MustDelete("zzzig_testFile")
	udwFile.MustWriteFileWithMkdir("zzzig_testFile/0.txt", []byte{})
	bw := &udwBytes.BufWriter{}
	udwFile.MustReadFileToBufW("zzzig_testFile/0.txt", bw)
	udwTest.Equal(len(bw.GetBytes()), 0)

	for _, size := range []int{
		1024,
		32*1024 - 1,
		32 * 1024,
		32*1024 + 1,
		1024 * 1024,
	} {
		content1 := bytes.Repeat([]byte("1"), size)
		udwFile.MustWriteFileWithMkdir("zzzig_testFile/1.txt", content1)
		bw = &udwBytes.BufWriter{}
		udwFile.MustReadFileToBufW("zzzig_testFile/1.txt", bw)
		udwTest.Equal(bw.GetBytes(), content1)
		udwTest.BenchmarkWithRepeatNum(1e2, func() {
			bw.Reset()
			udwFile.MustReadFileToBufW("zzzig_testFile/1.txt", bw)
		})
	}

}
