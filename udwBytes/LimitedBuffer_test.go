package udwBytes

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestLog(t *testing.T) {
	logBuffer := NewLimitedBuffer(5)
	logBuffer.Add("hello")
	logBuffer.Add("world")
	udwTest.Equal(logBuffer.ToString(), "world")
	logBuffer.Reset()
	logBuffer.Add("hello")
	logBuffer.Add("world!!!")
	udwTest.Equal(logBuffer.ToString(), "ld!!!")

	logBuffer.Reset()
	udwTest.Equal(logBuffer.ToString(), "")
	logBuffer.Add("hell")
	logBuffer.Add("o")
	logBuffer.Add("world")
	udwTest.Equal(logBuffer.ToString(), "world")

	logBuffer = NewLimitedBuffer(14)
	logBuffer.Add("hello")
	logBuffer.AddLine("world")
	udwTest.Equal(logBuffer.ToString(), "helloworld\n")
}
