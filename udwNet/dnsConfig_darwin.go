// +build !js

package udwNet

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"strconv"
	"strings"
)

func MustSetDnsServerAddr(ip string) {
	udwCmd.CmdSlice([]string{"networksetup", "-setdnsservers", GetDefaultNetworksetupName(), ip}).MustCombinedOutput()
	udwCmd.MustCombinedOutputAndNotExitStatusCheck("killall -HUP mDNSResponder")
}

func MustSetDnsServerToDefault() {
	udwCmd.CmdSlice([]string{"networksetup", "-setdnsservers", GetDefaultNetworksetupName(), "Empty"}).MustCombinedOutput()
}

func GetDefaultNetworksetupName() string {
	routeRule := MustGetDefaultRouteRule()
	devName := routeRule.GetOutInterface().GetName()
	output := udwCmd.MustCombinedOutput("networksetup listnetworkserviceorder")
	devServiceMap := parseNetworksetupListnetworkserviceorder(string(output))
	name, ok := devServiceMap[devName]
	if !ok {
		panic(fmt.Errorf("[GetDefaultNetworksetupName] can not found networkservice with devName %s", devName))
	}
	return name
}

func parseNetworksetupListnetworkserviceorder(content string) map[string]string {
	out := map[string]string{}
	lastNetworkserviceName := ""

	for _, line := range strings.Split(string(content), "\n") {

		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "(") {
			continue
		}

		rightIndex := strings.LastIndex(line, ")")
		if rightIndex == -1 {
			continue
		}
		_, err := strconv.Atoi(line[1:rightIndex])
		if err == nil {
			lastNetworkserviceName = line[rightIndex+2:]
			continue
		}

		line = strings.TrimSpace(line)
		DeviceIndex := strings.Index(line, "Device:")
		if DeviceIndex == -1 || DeviceIndex+len("Device:") >= rightIndex {
			continue
		}
		devName := strings.TrimSpace(line[DeviceIndex+len("Device:") : rightIndex])
		if devName != "" {
			out[devName] = lastNetworkserviceName
		}

	}
	return out
}
