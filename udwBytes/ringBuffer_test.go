package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestNewRingBuffer(ot *testing.T) {
	buf := QueueByteSlice{}
	buf.Init()
	buf.AddOne([]byte("123"))
	udwTest.Equal(buf.HasData(), true)
	udwTest.Equal(buf.GetOne(), []byte("123"))
	buf.RemoveOne()
	udwTest.Equal(buf.HasData(), false)
	buf.AddOne([]byte("1234"))
	buf.AddOne([]byte("4321"))
	udwTest.Equal(buf.HasData(), true)
	udwTest.Equal(buf.GetOne(), []byte("1234"))
	buf.RemoveOne()
	udwTest.Equal(buf.GetOne(), []byte("4321"))
	buf.RemoveOne()
	udwTest.Equal(buf.HasData(), false)
}

func TestQueueInt(ot *testing.T) {
	buf := QueueInt{}
	buf.Init()
	buf.AddOne(1)
	udwTest.Equal(buf.HasData(), true)
	udwTest.Equal(buf.GetOne(), 1)
	buf.RemoveOne()
	udwTest.Equal(buf.HasData(), false)
	buf.AddOne(2)
	buf.AddOne(3)
	udwTest.Equal(buf.HasData(), true)
	udwTest.Equal(buf.GetOne(), 2)
	buf.RemoveOne()
	udwTest.Equal(buf.GetOne(), 3)
	buf.RemoveOne()
	udwTest.Equal(buf.HasData(), false)
}
