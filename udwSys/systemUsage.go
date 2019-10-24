package udwSys

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwCmd"
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwPlatform"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"strings"
)

func IsNotSupport() bool {
	return udwPlatform.IsLinux() == false
}

func UsagePreInstall() {
	if IsNotSupport() {
		return
	}
	if udwCmd.Exist("mpstat") && udwCmd.Exist("netstat") {
		return
	}
	udwCmd.CmdBash("sudo apt-get update").MustRunWithFullError()
	udwCmd.CmdBash("sudo apt-get install -y sysstat").MustRunWithFullError()
}

func NetworkConnection() (connectionCount int) {
	if IsNotSupport() {
		return
	}
	return networkConnection(string(udwCmd.CmdSlice([]string{"bash", "-c", "netstat -na | grep ESTABLISHED | wc -l"}).MustCombinedOutput()))
}

func networkConnection(output string) (connectionCount int) {
	output = strings.TrimSpace(output)
	return udwStrconv.AtoIDefault0(output)
}

func IKEUserCount() (out int) {
	if IsNotSupport() {
		return
	}
	if !udwCmd.Exist("swanctl") {
		return 0
	}
	err := udwErr.PanicToError(func() {
		out = ikeUserCount(string(udwCmd.MustCombinedOutputAndNotExitStatusCheck("swanctl -S")))
	})
	if err != nil {
		fmt.Println("error", "IKEUserCount", err.Error())
	}
	return out

}

func ikeUserCount(output string) int {
	lines := strings.Split(output, "\n")
	c := ""
	if !strings.Contains(output, "IKE_SAs") {
		return 0
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "IKE_SAs") {
			_line := strings.Split(line, "total")
			c = _line[0]
			c = strings.Trim(c, "IKE_SAs:")
			c = strings.TrimSpace(c)
			break
		}
	}
	return udwStrconv.AtoIDefault0(c)
}
