package udwProfileDelay

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwTime"
	"sync"
	"time"
)

const disable = true

var gDelay delayProfiler

func P() {
	if disable == false {
		PN(udwDebug.GetCurrentOneStackString(1))
	}
}

func PN(pos string) {
	if disable == false {
		t := time.Now()
		isFirst := false
		var thisLastTime time.Time
		var thisStartTime time.Time
		gDelay.locker.Lock()
		if gDelay.startTime.IsZero() {
			gDelay.startTime = t
			gDelay.lastTime = t
			isFirst = true
		} else {
			thisLastTime = gDelay.lastTime
			thisStartTime = gDelay.startTime
			gDelay.lastTime = t
		}
		gDelay.locker.Unlock()
		if isFirst {
			fmt.Println("u37ypwcdge first", udwTime.DefaultFormatNs(t), pos)
		} else {
			fmt.Println("u37ypwcdge", udwTime.DurationFormatPadding(t.Sub(thisLastTime)), udwTime.DurationFormatPadding(t.Sub(thisStartTime)), pos)
		}
	}
}

type delayProfiler struct {
	locker    sync.Mutex
	startTime time.Time
	lastTime  time.Time
}
