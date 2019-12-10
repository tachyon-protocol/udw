package tyVpnRouteServer

import (
	"github.com/tachyon-protocol/udw/tyVpnClient"
	"github.com/tachyon-protocol/udw/tyVpnRouteServer/tyVpnRouteClient"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwRpc2"
	"github.com/tachyon-protocol/udw/udwSqlite3"
	"time"
)

func (serverRpcObj) VpnNodeRegister(clientIp udwRpc2.PeerIp, thisNode tyVpnRouteClient.VpnNode) (errMsg string) {
	startTime := time.Now().UTC()
	thisNode.UpdateTime = startTime.Truncate(time.Second)
	if thisNode.Ip == "" {
		thisNode.Ip = clientIp.Ip
	}
	err := tyVpnClient.Ping(tyVpnClient.PingReq{
		Ip:        thisNode.Ip,
		ServerChk: thisNode.ServerChk,
	})
	if err != nil {
		return "f3pbhbjveg " + err.Error()
	}
	getDb().MustSet(k1VpnNodeIp, thisNode.Ip, udwJson.MustMarshalToString(thisNode))
	return ""
}

func (serverRpcObj) VpnNodeList() []tyVpnRouteClient.VpnNode {
	outList := []tyVpnRouteClient.VpnNode{}
	getDb().MustGetRangeCallback(udwSqlite3.GetRangeReq{
		K1: k1VpnNodeIp,
	}, func(key string, value string) {
		var thisNode tyVpnRouteClient.VpnNode
		udwJson.MustUnmarshalFromString(value, &thisNode)
		if isNodeTimeout(thisNode) == false {
			outList = append(outList, thisNode)
		} else {
			getDb().MustDeleteWithKv(k1VpnNodeIp, key, value)
		}
	})
	return outList
}

func (serverRpcObj) Ping() {}

func initGcVpnNode() {
	go func() {
		for {
			time.Sleep(k1VpnNodeTtl)
			getDb().MustGetRangeCallback(udwSqlite3.GetRangeReq{
				K1: k1VpnNodeIp,
			}, func(key string, value string) {
				var thisNode tyVpnRouteClient.VpnNode
				udwJson.MustUnmarshalFromString(value, &thisNode)
				if isNodeTimeout(thisNode) {
					getDb().MustDeleteWithKv(k1VpnNodeIp, key, value)
				}
			})
		}
	}()
}

func isNodeTimeout(thisNode tyVpnRouteClient.VpnNode) bool {
	return time.Now().Add(-k1VpnNodeTtl).After(thisNode.UpdateTime)
}

const k1VpnNodeIp = "k1VpnNodeIp2"
const k1VpnNodeTtl = time.Minute
