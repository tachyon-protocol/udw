package udwNetTestV2

import (
	"io"
	"testing"
)

func TestRunTestOnPipe(ot *testing.T) {
	RunTestOnPipe(func() (io.ReadWriteCloser, io.ReadWriteCloser) {
		pipe := NewPipeConns(4)
		return pipe.Conn1(), pipe.Conn2()
	})
}
