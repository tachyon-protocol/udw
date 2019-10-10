package udwNet_test

import (
	"bytes"
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwNet/udwNetTestV2"
	"github.com/tachyon-protocol/udw/udwTest"
	"io"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConnTwoWayCopy(ot *testing.T) {
	listerner := udwNetTestV2.NewInmemoryListener()
	closer := udwNet.TcpNewListenerFromExistListener(listerner, func(conn net.Conn) {
		defer conn.Close()
		err := conn.SetDeadline(time.Now().Add(time.Second * 4))
		if err != nil {
			if udwNet.IsSocketCloseError(err) {
				return
			}
			fmt.Println("err conn.SetDeadline", err.Error())
			return
		}

		buf := make([]byte, 1024*32)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if udwNet.IsSocketCloseError(err) {
					return
				}
				fmt.Println("err conn.Read", err.Error())
				return
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				if udwNet.IsSocketCloseError(err) {
					return
				}
				fmt.Println("err conn.Write", err.Error())
				return
			}
		}
	})
	defer closer()
	listerner2 := udwNetTestV2.NewInmemoryListener()
	closer2 := udwNet.TcpNewListenerFromExistListener(listerner2, func(conn net.Conn) {
		defer conn.Close()
		nextConn, err := listerner.Dial()
		udwErr.PanicIfError(err)
		udwNet.ConnTwoWayCopy(udwNet.ConnTwoWayCopyRequest{
			FromHopConn: conn,
			NextHopConn: nextConn,
			AddTimeout:  time.Second,
		})
	})
	defer closer2()
	func() {
		conn1, err := listerner2.Dial()
		udwErr.PanicIfError(err)
		defer conn1.Close()
		buf := bytes.Repeat([]byte{1}, 1024)
		for i := 0; i < 2; i++ {
			_, err = conn1.Write(buf)
			udwErr.PanicIfError(err)
			buf2 := make([]byte, 1024)
			_, err = io.ReadFull(conn1, buf2)
			udwErr.PanicIfError(err)
			udwTest.Equal(buf, buf2)
		}
		conn1.Close()
	}()

	func() {
		conn1, err := listerner2.Dial()
		udwErr.PanicIfError(err)
		buf := bytes.Repeat([]byte{1}, 1024*32+1)
		buf2 := make([]byte, 1024*32+1)
		for i := 0; i < 10; i++ {
			conn1.Write(buf)
			io.ReadFull(conn1, buf2)
		}

		udwTest.Benchmark(func() {
			const num = 1000
			udwTest.BenchmarkSetNum(num)
			udwTest.BenchmarkSetBytePerRun(1024*32 + 1)
			for i := 0; i < num; i++ {
				conn1.Write(buf)

				io.ReadFull(conn1, buf2)

			}
		})
	}()

}

func TestConnTwoWayCopy2(ot *testing.T) {
	closer := udwNet.TcpNewListener("127.0.0.1:34567", func(conn net.Conn) {
		defer conn.Close()
		err := conn.SetDeadline(time.Now().Add(time.Second * 4))
		if err != nil {
			if udwNet.IsSocketCloseError(err) {
				return
			}
			fmt.Println("err conn.SetDeadline", err.Error())
			return
		}

		buf := make([]byte, 1024*32)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if udwNet.IsSocketCloseError(err) {
					return
				}
				fmt.Println("err conn.Read", err.Error())
				return
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				if udwNet.IsSocketCloseError(err) {
					return
				}
				fmt.Println("err conn.Write", err.Error())
				return
			}
		}
	})
	defer closer()
	closer2 := udwNet.TcpNewListener("127.0.0.1:34568", func(conn net.Conn) {
		defer conn.Close()
		nextConn, err := net.DialTimeout("tcp", "127.0.0.1:34567", time.Second)
		udwErr.PanicIfError(err)
		udwNet.ConnTwoWayCopy(udwNet.ConnTwoWayCopyRequest{
			FromHopConn: conn,
			NextHopConn: nextConn,
			AddTimeout:  time.Millisecond * 100,
		})
	})
	defer closer2()

	func() {
		conn1, err := net.DialTimeout("tcp", "127.0.0.1:34568", time.Second)
		udwErr.PanicIfError(err)
		defer conn1.Close()
		buf := bytes.Repeat([]byte{1}, 1024)
		for i := 0; i < 2; i++ {
			_, err = conn1.Write(buf)
			udwErr.PanicIfError(err)
			buf2 := make([]byte, 1024)
			_, err = io.ReadFull(conn1, buf2)
			udwErr.PanicIfError(err)
			udwTest.Equal(buf, buf2)
		}
		conn1.Close()
	}()
	func() {
		conn1, err := net.DialTimeout("tcp", "127.0.0.1:34568", time.Second)
		udwErr.PanicIfError(err)
		const bufSize = 1024*32 + 1
		buf := bytes.Repeat([]byte{1}, bufSize)
		buf2 := make([]byte, bufSize)
		for i := 0; i < 10; i++ {
			conn1.Write(buf)
			io.ReadFull(conn1, buf2)
		}

		udwTest.Benchmark(func() {
			const num = 1000
			udwTest.BenchmarkSetNum(num)
			udwTest.BenchmarkSetBytePerRun(bufSize)
			for i := 0; i < num; i++ {
				conn1.Write(buf)
				io.ReadFull(conn1, buf2)
			}
		})
	}()
}

func TestConnTwoWayCopy3(ot *testing.T) {
	testFn := func(p1 func(conn net.Conn), p2 func(conn net.Conn)) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		closer := udwNet.TcpNewListener("127.0.0.1:34567", func(conn net.Conn) {
			p1(conn)
			wg.Done()
		})
		defer closer()
		closer2 := udwNet.TcpNewListener("127.0.0.1:34568", func(conn net.Conn) {
			defer conn.Close()
			nextConn, err := net.DialTimeout("tcp", "127.0.0.1:34567", time.Second)
			udwErr.PanicIfError(err)
			udwNet.ConnTwoWayCopy(udwNet.ConnTwoWayCopyRequest{
				FromHopConn: conn,
				NextHopConn: nextConn,
				AddTimeout:  time.Millisecond * 100,
			})
		})
		defer closer2()

		wg.Add(1)
		func() {
			conn1, err := net.DialTimeout("tcp", "127.0.0.1:34568", time.Second)
			udwErr.PanicIfError(err)
			p2(conn1)
			wg.Done()
		}()
		wg.Wait()
	}
	const t1Size = 256 * 1024
	writePartFn := func(conn net.Conn) {
		err := conn.SetDeadline(time.Now().Add(time.Second))
		udwErr.PanicIfError(err)
		buf := bytes.Repeat([]byte{2}, t1Size)
		_, err = conn.Write(buf)
		udwErr.PanicIfError(err)
	}
	readPartFn := func(conn net.Conn) {
		err := conn.SetDeadline(time.Now().Add(time.Second))
		udwErr.PanicIfError(err)
		buf := make([]byte, t1Size)
		_, err = io.ReadFull(conn, buf)
		for i := 0; i < t1Size; i++ {
			if buf[i] != 2 {
				panic("fail " + strconv.Itoa(i))
			}
		}
		udwErr.PanicIfError(err)
	}
	p1 := func(conn net.Conn) {
		defer conn.Close()
		for i := 0; i < 5; i++ {
			writePartFn(conn)
			readPartFn(conn)
		}
	}
	p2 := func(conn net.Conn) {
		defer conn.Close()
		for index := 0; index < 5; index++ {
			readPartFn(conn)
			writePartFn(conn)
		}
	}
	testFn(p1, p2)
	testFn(p2, p1)
}
