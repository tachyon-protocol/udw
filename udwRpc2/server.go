package udwRpc2

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwNet"
	"github.com/tachyon-protocol/udw/udwShm"
	"net"
	"time"
)

type ServerReq struct {
	Addr        string
	MaxIdleTime time.Duration
	Handler     func(ctx *ReqCtx)
}
type ServerHub struct {
	req    ServerReq
	closer udwClose.Closer
}

func NewServerHub(req ServerReq) *ServerHub {
	if req.MaxIdleTime == 0 {
		req.MaxIdleTime = 4 * time.Minute
	}
	sh := &ServerHub{
		req: req,
	}
	closerFn := udwNet.TcpNewListener(req.Addr, func(conn net.Conn) {
		conn2 := &Conn{}
		conn2.wb = *udwShm.NewShmWriter(conn, 0)
		conn2.rb = *udwShm.NewShmReader(conn, 0)
		conn2.conn = conn
		for {
			conn2.conn.SetDeadline(time.Now().Add(req.MaxIdleTime))

			errMsg := conn2.rb.ReadArrayStart()
			if errMsg != "" {
				fmt.Println("dt883f58d6")
				conn2.Close()
				return
			}
			conn2.conn.SetDeadline(time.Time{})
			ctx := &ReqCtx{
				conn: conn2,
			}
			errMsg = udwErr.PanicToErrorMsgWithStack(func() {
				req.Handler(ctx)
			})
			if errMsg != "" {
				fmt.Println("det9aufaz4 ", errMsg)
				ctx.CloseConn()
				return
			}
			ctx.Close()

		}
	})
	sh.closer.AddOnClose(closerFn)
	return sh
}

func (sh *ServerHub) Close() {
	sh.closer.Close()
}
