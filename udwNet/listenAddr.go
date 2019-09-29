// +build !js

package udwNet

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwLog"
	"strconv"
	"strings"
)

func GetPossiableConnectAddrStringList(listenAddr string) (outList []string) {
	if !strings.HasPrefix(listenAddr, ":") {
		return []string{listenAddr}
	}
	port, err := strconv.Atoi(listenAddr[1:])
	if err != nil {
		return nil
	}
	err = udwErr.PanicToError(func() {
		outList = MustGetCurrentIpWithPortList(uint16(port))
	})
	if err != nil {
		udwLog.Log("error", "[GetPossiableConnectAddrStringList]", err.Error())
	}
	return outList
}

func GetPossiableConnectAddrDebugString(listenAddr string) string {
	return strings.Join(GetPossiableConnectAddrStringList(listenAddr), ", ")
}
