package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestGetFixTimeSleepTimeWithOffset(ot *testing.T) {
	now := time.Date(2010, 1, 1, 0, 0, 12, 0, time.UTC)
	dur := time.Second * 10
	offset := time.Second * 5
	thisDur := getFixTimeSleepTimeWithOffset(now, dur, offset)
	udwTest.Equal(thisDur, 3*time.Second)
	now = time.Date(2010, 1, 1, 0, 0, 15, 0, time.UTC)
	thisDur = getFixTimeSleepTimeWithOffset(now, dur, offset)
	udwTest.Equal(thisDur, 10*time.Second)
	now = time.Date(2010, 1, 1, 0, 0, 17, 0, time.UTC)
	thisDur = getFixTimeSleepTimeWithOffset(now, dur, offset)
	udwTest.Equal(thisDur, 8*time.Second)
}

func TestSleepLoopWithFixDuration(ot *testing.T) {

}
