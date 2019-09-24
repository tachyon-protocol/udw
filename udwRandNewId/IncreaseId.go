package udwRandNewId

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"strconv"
	"sync/atomic"
)

var intId uint64

func NextUint64IdString() string {
	var idInt = NextUint64Id()
	return strconv.FormatUint(idInt, 10)
}

func NextUint64IdStringWithPadding() string {
	var idInt = NextUint64Id()
	return udwStrconv.FormatUint64Padding(idInt)
}

func NextUint64Id() uint64 {
	var idInt = atomic.AddUint64(&intId, 1)
	return idInt
}
