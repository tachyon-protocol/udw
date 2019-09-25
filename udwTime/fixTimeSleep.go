package udwTime

import (
	"github.com/tachyon-protocol/udw/udwRand"
	"time"
)

func FixClockTimeSleep(dur time.Duration) {
	now := time.Now()
	thisSleepDur := now.Truncate(dur).Add(dur).Sub(now)
	time.Sleep(thisSleepDur)

}

func FixClockTimeSleepWithOffset(dur time.Duration, offset time.Duration) {
	now := time.Now()
	thisSleepDur := getFixTimeSleepTimeWithOffset(now, dur, offset)
	time.Sleep(thisSleepDur)
}

func FixClockTimeSleepToTimeInDayInDefaultTimeZone(offset time.Duration) {
	FixClockTimeSleepToTimeInDay(offset, GetDefaultTimeZone())
}

func FixClockTimeSleepToTimeInDay(offset time.Duration, zone *time.Location) {
	now := time.Now().In(zone)
	thisSleepDur := ToDate(now).Add(offset).Sub(now)
	if thisSleepDur < 0 {
		thisSleepDur = ToDate(now).Add(Day + offset).Sub(now)
	}
	time.Sleep(thisSleepDur)
}

func sleepToTime(now time.Time, t time.Time) {
	dur := t.Sub(now)
	if dur <= 0 {
		return
	}
	time.Sleep(dur)
}

func SleepToTodayTimePoint(now time.Time, d time.Duration) {
	dayBegin := ToDate(now)
	dur := d - now.Sub(dayBegin)
	if dur <= 0 {
		return
	}
	time.Sleep(dur)
}

func GetNextTimeDuration(hour int, minute int) time.Duration {
	left := NowFromDefaultNower()
	right := left.Add(time.Hour * 24)
	today := ToDateDefault(left)
	tomorrow := ToDateDefault(right)
	delta := time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute)
	next := today.Add(delta)
	if next.Before(left) {
		next = tomorrow.Add(delta)
	}
	return next.Sub(left)
}

type FixTimeDurationLoopRequest struct {
	FixSleepDur          time.Duration
	RandomDurRange       time.Duration
	SleepDurAfterTimeout time.Duration
	Fn                   func()
}

func SleepLoopWithFixDuration(req FixTimeDurationLoopRequest) {
	changeDur := time.Duration(udwRand.Float64Between(-0.5, 0.5) * float64(req.RandomDurRange))
	time.Sleep(req.FixSleepDur + changeDur)
	for {
		startTime := time.Now()
		req.Fn()
		dur := time.Since(startTime)
		changeDur := time.Duration(udwRand.Float64Between(-0.5, 0.5) * float64(req.RandomDurRange))
		if dur < req.FixSleepDur {
			time.Sleep(req.FixSleepDur - dur + changeDur)
			continue
		}
		if req.SleepDurAfterTimeout > 0 {
			time.Sleep(req.SleepDurAfterTimeout + changeDur)
		}
	}
}

func getFixTimeSleepTimeWithOffset(now time.Time, dur time.Duration, offset time.Duration) time.Duration {
	nextTime := now.Truncate(dur).Add(offset)
	if nextTime.After(now) == false {
		nextTime = nextTime.Add(dur)
	}
	thisSleepDur := nextTime.Sub(now)
	return thisSleepDur
}
