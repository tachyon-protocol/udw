package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"testing"
)

func TestUvarint(ot *testing.T) {
	buf := make([]byte, 10)
	type testCase struct {
		number uint64
		size   int
	}
	bufW := &BufWriter{}
	for _, cas := range []testCase{
		{0, 1},
		{1, 1},
		{1<<7 - 1, 1},
		{1 << 7, 2},
		{1<<7 + 1, 2},
		{1<<14 - 1, 2},
		{1 << 14, 3},
		{1<<21 - 1, 3},
		{1 << 21, 4},
		{1<<28 - 1, 4},
		{1 << 28, 5},
		{1<<35 - 1, 5},
		{1 << 35, 6},
		{1<<42 - 1, 6},
		{1 << 42, 7},
		{1<<49 - 1, 7},
		{1 << 49, 8},
		{1<<56 - 1, 8},
		{1 << 56, 9},
		{1<<63 - 1, 9},
		{1 << 63, 10},
		{math.MaxUint64, 10},
	} {
		n := WriteUvarint(buf, cas.number)
		udwTest.Equal(n, cas.size)
		udwTest.Equal(GetUvarintOutputSize(cas.number), cas.size)
		bufW.Reset()
		bufW.WriteUvarint(cas.number)
		udwTest.Equal(bufW.GetLen(), cas.size)
		udwTest.Equal(buf[:n], bufW.GetBytes())
		bufR := NewBufReader(bufW.GetBytes())
		v, ok := bufR.ReadUvarint()
		udwTest.Equal(ok, true)
		udwTest.Equal(v, uint64(cas.number))
	}
	{
		type testCase struct {
			number uint64
			result []byte
		}
		for _, cas := range []testCase{
			{uint64(0x0), []byte{0x0}},
			{uint64(0x80), []byte{0x80, 0x01}},
			{uint64(0xff), []byte{0xff, 0x01}},
			{uint64(0x1000), []byte{0x80, 0x20}},
			{uint64(0x10000), []byte{0x80, 0x80, 0x04}},
			{uint64(0x100000), []byte{0x80, 0x80, 0x40}},
			{uint64(0x1000000), []byte{0x80, 0x80, 0x80, 0x08}},
		} {
			bufW := &BufWriter{}
			bufW.WriteUvarint(cas.number)
			udwTest.Equal(bufW.GetBytes(), cas.result)

		}

	}
}
