package udwNetTestV2

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"strconv"
	"sync"
)

func RunTestOnPipe(newConn func() (io.ReadWriteCloser, io.ReadWriteCloser)) {
	runTestOnPipeCase1(newConn)
	runTestOnPipeCase2(newConn)
}

func runTestOnPipeCase1(newConn func() (io.ReadWriteCloser, io.ReadWriteCloser)) {
	sConn, cConn := newConn()
	defer sConn.Close()
	defer cConn.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			buf := make([]byte, 4)
			n, err := sConn.Read(buf)
			udwTest.Equal(err, nil)
			udwTest.Equal(n, 4)
		}
		for i := 0; i < 10; i++ {
			n, err := sConn.Write([]byte("123" + strconv.Itoa(i)))
			udwTest.Equal(err, nil)
			udwTest.Equal(n, 4)
		}
		sConn.Close()
		wg.Done()
	}()
	for i := 0; i < 10; i++ {
		n, err := cConn.Write([]byte("123" + strconv.Itoa(i)))
		udwTest.Equal(err, nil)
		udwTest.Equal(n, 4)
	}
	for i := 0; i < 10; i++ {
		buf := make([]byte, 4)
		n, err := cConn.Read(buf)
		udwTest.Equal(err, nil)
		udwTest.Equal(n, 4)
	}
	cConn.Close()
	wg.Wait()
}

func runTestOnPipeCase2(newConn func() (io.ReadWriteCloser, io.ReadWriteCloser)) {
	sConn, cConn := newConn()
	defer sConn.Close()
	defer cConn.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			n, err := sConn.Write([]byte("123" + strconv.Itoa(i)))
			udwTest.Equal(err, nil)
			udwTest.Equal(n, 4)
		}
		for i := 0; i < 10; i++ {
			buf := make([]byte, 4)
			n, err := sConn.Read(buf)
			udwTest.Equal(err, nil)
			udwTest.Equal(n, 4)
		}
		sConn.Close()
		wg.Done()
	}()
	for i := 0; i < 10; i++ {
		buf := make([]byte, 4)
		n, err := cConn.Read(buf)
		udwTest.Equal(err, nil)
		udwTest.Equal(n, 4)
	}
	for i := 0; i < 10; i++ {
		n, err := cConn.Write([]byte("123" + strconv.Itoa(i)))
		udwTest.Equal(err, nil)
		udwTest.Equal(n, 4)
	}
	cConn.Close()
	wg.Wait()
}

func RunTestTwoRwc(c io.ReadWriteCloser, s io.ReadWriteCloser) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		buf := make([]byte, 4096)
		nr, err := s.Read(buf)
		udwErr.PanicIfError(err)
		udwTest.Equal(buf[:nr], []byte{1})
		wg.Done()
	}()
	_, err := c.Write([]byte{1})
	udwErr.PanicIfError(err)
	wg.Wait()
	for i := 0; i < 10; i++ {
		_, err := c.Write([]byte{1})
		udwErr.PanicIfError(err)
		buf := make([]byte, 4096)
		nr, err := s.Read(buf)
		udwErr.PanicIfError(err)
		udwTest.Equal(buf[:nr], []byte{1})
	}
	c.Close()
	s.Close()
}
