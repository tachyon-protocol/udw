package udwIo

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTime"
	"io"
	"sync/atomic"
	"time"
)

type debugRwc struct {
	io.ReadWriteCloser
	Name     string
	showData bool
}

func NewDebugRwc(rwc io.ReadWriteCloser, name string) debugRwc {
	return debugRwc{
		ReadWriteCloser: rwc,
		Name:            name,
		showData:        true,
	}
}

func (c debugRwc) Write(b []byte) (n int, err error) {
	fmt.Printf("[debugConn] [%s] %s Write Start len: %d\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b))
	n, err = c.ReadWriteCloser.Write(b)
	if err != nil {
		if c.showData {
			fmt.Printf("[debugConn] [%s] %s Write finish len: %d err: %s content: %#v\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, n, err, b)
		} else {
			fmt.Printf("[debugConn] [%s] %s Write finish len: %d err: %s\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, n, err)
		}
	} else {
		if c.showData {
			fmt.Printf("[debugConn] [%s] %s Write finish len: %d content: %#v\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, n, b)
		} else {
			fmt.Printf("[debugConn] [%s] %s Write finish len: %d\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, n)
		}
	}
	return n, err
}

func (c debugRwc) Read(b []byte) (n int, err error) {
	fmt.Printf("[debugConn] [%s] %s Read Start len: %d\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b))
	n, err = c.ReadWriteCloser.Read(b)
	if err != nil {
		if c.showData {
			fmt.Printf("[debugConn] [%s] %s Read finish iLen: %d oLen: %d err: %s content: %#v\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b), n, err, b[:n])
		} else {
			fmt.Printf("[debugConn] [%s] %s Read finish iLen: %d oLen: %d err: %s\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b), n, err)
		}
	} else {
		if c.showData {
			fmt.Printf("[debugConn] [%s] %s Read finish iLen: %d oLen: %d content: %#v\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b), n, b[:n])
		} else {
			fmt.Printf("[debugConn] [%s] %s Read finish iLen: %d oLen: %d\n", udwTime.MysqlUsNowFromDefaultNower(), c.Name, len(b), n)
		}
	}
	return n, err
}

func (c debugRwc) Close() (err error) {
	fmt.Println("[debugConn]", "["+udwTime.MysqlUsNowFromDefaultNower()+"]", c.Name, "Close start err:", err)
	err = c.ReadWriteCloser.Close()
	fmt.Println("[debugConn]", "["+udwTime.MysqlUsNowFromDefaultNower()+"]", c.Name, "Close finish err:", err)
	return err
}

func NewDebugRwcNoData(rwc io.ReadWriteCloser, name string) debugRwc {
	return debugRwc{
		ReadWriteCloser: rwc,
		Name:            name,
		showData:        false,
	}
}

type sumSizeRwc struct {
	io.ReadWriteCloser
	Name       string
	startTime  time.Time
	writeBytes uint64
	readBytes  uint64
	readNum    uint64
	writeNum   uint64
	hasClose   bool
}

func (c *sumSizeRwc) Write(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Write(b)
	if n > 0 {
		atomic.AddUint64(&c.writeBytes, uint64(n))
		atomic.AddUint64(&c.writeNum, 1)
	}
	return n, err
}

func (c *sumSizeRwc) Read(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Read(b)
	if n > 0 {
		atomic.AddUint64(&c.readBytes, uint64(n))
		atomic.AddUint64(&c.readNum, 1)
	}
	return n, err
}

func (c *sumSizeRwc) Close() (err error) {
	err = c.ReadWriteCloser.Close()
	if !c.hasClose {
		fmt.Printf("[sumSizeRwc] [%s] read[bytes:%s num:%d] write[bytes:%s num:%d] duration:%s\n",
			c.Name, udwStrconv.GbFromFloat64(float64(c.readBytes)), c.readNum, udwStrconv.GbFromFloat64(float64(c.writeBytes)), c.writeNum,
			time.Since(c.startTime))
		c.hasClose = true
	}
	return err
}

func NewSumSizeRwc(rwc io.ReadWriteCloser, name string) io.ReadWriteCloser {
	return &sumSizeRwc{
		ReadWriteCloser: rwc,
		Name:            name,
		startTime:       time.Now(),
	}
}
