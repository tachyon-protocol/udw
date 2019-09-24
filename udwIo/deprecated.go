package udwIo

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"time"
)

func FmtByteSpeed(byteNum int, dur time.Duration) string {
	bytePerSecond := float64(byteNum) / (float64(dur) / float64(time.Second))
	return udwStrconv.GbFromFloat64(bytePerSecond) + "/s"
}

func FmtByteNum(byteNum int) string {
	return udwStrconv.GbFromInt64(int64(byteNum))
}
