package udwStrconv

import (
	"fmt"
	"strings"
	"time"
)

func PrintProgress(get int64, total int64, dur time.Duration, lastBytes int64) {
	percent := 0.0
	if total <= 0 {
		percent = 0.0
	} else if total < get {
		percent = 1.0
	} else {
		percent = float64(get) / float64(total)
	}
	showNum := int(percent * 40)
	notShowNum := 40 - showNum
	fmt.Printf("%s%s %.2f%% %s/%s %s     \r",
		strings.Repeat("#", showNum), strings.Repeat(" ", notShowNum), percent*100,
		GbFromInt64(get), GbFromInt64(total), GbSpeedFromFloat64AndDuration(float64(lastBytes), dur))
}
