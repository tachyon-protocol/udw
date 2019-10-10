package udwRandNewId

import (
	"github.com/tachyon-protocol/udw/udwTime"
	"time"
)

func IncreaseStringId() string {
	return udwTime.MustDbTimeGetStringFromObj(time.Now()) + "_" + NextUint64IdStringWithPadding()
}
