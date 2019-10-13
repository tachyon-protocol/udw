package udwTcpPing

import "github.com/tachyon-protocol/udw/udwConsole"

func AddCommand() {
	udwConsole.AddCommandWithName("TcpPing", func(req struct {
		Ip string
	}) {
		Ping(req.Ip)
	})
	udwConsole.AddCommandWithName("TcpPingServer", func() {
		RunServer()
		udwConsole.WaitForExit()
	})
}
