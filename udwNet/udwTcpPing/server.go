package udwTcpPing

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwLog"
	"github.com/tachyon-protocol/udw/udwNet"
	"net"
	"strconv"
	"time"
)

func RunServer() (closer func() error) {
	fmt.Println("SERVER PORT", port)
	return udwNet.RunTCPServerListenAddr(":"+strconv.Itoa(port), func(conn net.Conn) {
		buf := make([]byte, 1)

		defer conn.Close()
		udwLog.Log("client", conn.RemoteAddr())
		for {
			udwErr.PanicIfError(conn.SetDeadline(time.Now().Add(timeout)))
			_, err := conn.Read(buf)
			if err != nil {

				udwLog.Log("[7ryrjhcqj4] read err", err, conn.RemoteAddr())
				return
			}
			_, err = conn.Write(buf)
			if err != nil {

				udwLog.Log("[35d4m5ustb] write err", err, conn.RemoteAddr())
				return
			}
		}
	})
}
