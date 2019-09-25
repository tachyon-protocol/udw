package udwNet

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwLog"
	"strings"
)

func MustIsVpnConnected() bool {
	for _, device := range MustGetNetDeviceList() {
		if device.GetFirstIpv4IP() == nil {
			continue
		}
		name := device.GetName()

		if strings.Contains(name, "utun") || strings.Contains(name, "ipsec") {
			return true
		}
	}
	return false
}

func IsVpnConnectedIgnorePanic() bool {
	ret := false
	err := udwErr.PanicToError(func() {
		ret = MustIsVpnConnected()
	})
	if err != nil {
		udwLog.Log("error", "[IsVpnConnectedIgnorePanic]", err.Error())
	}
	return ret
}
