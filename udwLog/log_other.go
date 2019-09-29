// +build !ios,!android,!macAppStore

package udwLog

import (
	"github.com/tachyon-protocol/udw/udwBytes"
	"os"
)

func log(bufW *udwBytes.BufWriter) {
	bufW.WriteByte_('\n')
	os.Stdout.Write(bufW.GetBytes())
}
