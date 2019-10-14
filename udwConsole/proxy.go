package udwConsole

import (
	"github.com/tachyon-protocol/udw/udwCmd"
	"os"
)

func ProxyCommand(args ...string) func() {
	return func() {
		udwCmd.CmdSlice(append(args, os.Args...)).ProxyRun()
	}
}
