package udwNetTestV2

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestTcpPipe(t *testing.T) {
	c1, c2, err := TcpPipe()
	udwErr.PanicIfError(err)
	defer c1.Close()
	defer c2.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		buf := make([]byte, 4096)
		nr, err := c2.Read(buf)
		udwErr.PanicIfError(err)
		udwTest.Equal(buf[:nr], []byte{1})
		wg.Done()
	}()
	_, err = c1.Write([]byte{1})
	udwErr.PanicIfError(err)
	wg.Wait()
	for i := 0; i < 10; i++ {
		_, err := c1.Write([]byte{1})
		udwErr.PanicIfError(err)
		buf := make([]byte, 4096)
		nr, err := c2.Read(buf)
		udwErr.PanicIfError(err)
		udwTest.Equal(buf[:nr], []byte{1})

		_, err = c2.Write([]byte{1})
		udwErr.PanicIfError(err)
		nr, err = c1.Read(buf)
		udwErr.PanicIfError(err)
		udwTest.Equal(buf[:nr], []byte{1})
	}
}
