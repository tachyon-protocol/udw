package udwSys

import (
	"github.com/tachyon-protocol/udw/udwSys/udwSysEnv"
)

func RecoverPath() {
	udwSysEnv.RecoverPath()
}

func GetBinPathList() []string {
	return udwSysEnv.GetBinPathList()
}
