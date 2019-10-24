// +build darwin linux

package udwSys

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwCmd"
)

func mustIsRoot() bool {
	return bytes.Equal(udwCmd.MustCombinedOutput("whoami"), []byte("root\n"))
}
