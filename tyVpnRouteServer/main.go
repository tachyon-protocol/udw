package tyVpnRouteServer

import (
	"github.com/tachyon-protocol/udw/udwConsole"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"net"
	"net/http"
	"sync"
)

func RouteServerRunCmd() {
	initDb()
	initGcVpnNode()

	closer := Rpc_RunServer(":24587")
	udwConsole.WaitForExit()
	closer()
}

type serverRpcObj struct{}

var gSqlite3Db *udwSqlite3.Db
var gSqlite3DbOnce sync.Once

func initDb() {
	gSqlite3DbOnce.Do(func() {
		gSqlite3Db = udwSqlite3.MustNewDb(udwSqlite3.NewDbRequest{
			FilePath:                       "/usr/local/var/tachyonVpnInfoServer.sqlite3",
			EmptyDatabaseIfDatabaseCorrupt: true,
		})
	})
}

func getDb() *udwSqlite3.Db {
	return gSqlite3Db
}

func getClientIpStringIgnoreError(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return host
	}
	return ""
}
