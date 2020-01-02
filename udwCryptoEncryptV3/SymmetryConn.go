package udwCryptoEncryptV3

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/tachyon-protocol/udw/AesCtr"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwNet"
	"io"
	"sync"
)

const ErrMsgDecryptKey = "decrypt key error magic buf not match"

var gMagicBuf = []byte{0xc6, 0x1f, 0x2d, 0xae}

func MustSymmetryConn(conn io.ReadWriteCloser, key *[32]byte) (outConn io.ReadWriteCloser) {
	block, err := AesCtr.NewCipher((*key)[:])
	if err != nil {
		panic(err)
	}
	return &symmetryConn{
		rwc:   conn,
		block: block,
	}
}

func NewSymmetryConnWithBlock(conn io.ReadWriteCloser, block cipher.Block) io.ReadWriteCloser {
	return &symmetryConn{
		rwc:   conn,
		block: block,
	}
}

type symmetryConn struct {
	rwc       io.ReadWriteCloser
	block     cipher.Block
	wBuf      udwBytes.BufWriter
	wCtr      cipher.Stream
	rCtr      cipher.Stream
	hasWrite  bool
	hasRead   bool
	readLock  sync.Mutex
	writeLock sync.Mutex
	rBuf      [20]byte
}

func (c *symmetryConn) Write(src []byte) (n int, err error) {
	c.writeLock.Lock()
	if !c.hasWrite {

		c.hasWrite = true
		buf := c.wBuf.GetHeadBuffer(len(src) + 20)
		_, err = io.ReadFull(rand.Reader, buf[:16])
		if err != nil {
			c.writeLock.Unlock()
			c.Close()
			return 0, err
		}
		ctr := AesCtr.PoolGetAesCtr(c.block, buf[:16])
		ctr.XORKeyStream(buf[16:20], gMagicBuf)
		ctr.XORKeyStream(buf[20:], src)
		c.wCtr = ctr
		n, err := c.rwc.Write(buf)
		n = n - 20
		if n < 0 {
			n = 0
		}
		c.writeLock.Unlock()
		if err != nil {
			c.Close()
		}
		return n, err
	}
	if c.wCtr == nil {
		c.writeLock.Unlock()

		c.Close()
		return 0, errors.New(udwNet.ErrMsgSocketCloseError + " e25qs67py8")
	}
	buf := c.wBuf.GetHeadBuffer(len(src))
	c.wCtr.XORKeyStream(buf, src)
	n, err = c.rwc.Write(buf)
	if n != len(src) {
		if err == nil {
			err = io.ErrShortWrite
		}
	}
	c.writeLock.Unlock()
	return n, err
}

func (c *symmetryConn) Read(dst []byte) (n int, err error) {
	c.readLock.Lock()
	if !c.hasRead {

		c.hasRead = true
		buf := c.rBuf[:]
		_, err := io.ReadFull(c.rwc, buf)
		if err != nil {
			c.readLock.Unlock()
			c.Close()
			return 0, err
		}
		ctr := AesCtr.PoolGetAesCtr(c.block, buf[:16])
		ctr.XORKeyStream(buf[16:20], buf[16:20])
		if !bytes.Equal(buf[16:20], gMagicBuf) {
			c.readLock.Unlock()
			c.Close()
			return 0, errors.New(ErrMsgDecryptKey)
		}
		c.rCtr = ctr
	}
	if c.rCtr == nil {
		c.readLock.Unlock()

		c.Close()
		return 0, errors.New(udwNet.ErrMsgSocketCloseError + " 4vg8b6g4rn")
	}
	n, err = c.rwc.Read(dst)
	c.rCtr.XORKeyStream(dst[:n], dst[:n])
	c.readLock.Unlock()
	return n, err
}

func (c *symmetryConn) Close() (err error) {
	err = c.rwc.Close()
	c.writeLock.Lock()
	if c.wCtr != nil {
		AesCtr.PoolPutAesCtr(c.wCtr)
		c.wCtr = nil
	}
	c.writeLock.Unlock()
	c.readLock.Lock()
	if c.rCtr != nil {
		AesCtr.PoolPutAesCtr(c.rCtr)
		c.rCtr = nil
	}
	c.readLock.Unlock()
	return err
}

var gNewSymmetryConnWithBlockPool = sync.Pool{}

func PoolGetSymmetryConnWithBlock(conn io.ReadWriteCloser, block cipher.Block) io.ReadWriteCloser {
	obj := gNewSymmetryConnWithBlockPool.Get()
	if obj == nil {
		return NewSymmetryConnWithBlock(conn, block)
	}
	obj2, ok := obj.(*symmetryConn)
	if !ok {
		return NewSymmetryConnWithBlock(conn, block)
	}
	obj2.rwc = conn
	obj2.block = block
	obj2.wBuf.Reset()
	obj2.hasWrite = false
	obj2.hasRead = false
	return obj2
}

func PoolPutSymmetryConnAndClose(rwc io.ReadWriteCloser) {
	if rwc == nil {
		return
	}

	rwc.Close()
	obj2, ok := rwc.(*symmetryConn)
	if !ok {
		return
	}
	obj2.rwc = nil
	obj2.block = nil
	gNewSymmetryConnWithBlockPool.Put(obj2)
}
