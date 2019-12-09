package udwRpc2

import (
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwShm"
	"net"
	"time"
)

type Conn struct {
	wb            udwShm.ShmWriter
	rb            udwShm.ShmReader
	conn          net.Conn
	closer        udwClose.Closer
	idleStartTime time.Time
}

func (conn *Conn) Close() {
	conn.closer.CloseWithCallback(func() {
		conn.conn.Close()
	})
}

type ReqCtx struct {
	conn *Conn
}

func (ctx *ReqCtx) GetWriter() *udwShm.ShmWriter {
	return &ctx.conn.wb
}
func (ctx *ReqCtx) GetReader() *udwShm.ShmReader {
	return &ctx.conn.rb
}
func (ctx *ReqCtx) GetPeerIp() string {
	ip := udwNet.GetIpStringFromtNetAddrIgnoreNotExist(ctx.conn.conn.RemoteAddr())
	return ip
}

func (ctx *ReqCtx) Close() {
	ctx.conn.Close()
}
