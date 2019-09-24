package udwTime

import (
	"time"
)

func GetDurToNextWeek(weekday time.Weekday, hour int) time.Duration {
	now := NowFromDefaultNower()
	n := now.Weekday() - weekday
	if n < 0 {
		n = -n
	}
	if now.Weekday() == weekday && now.Hour() >= hour {
		n = 7
	}
	_now := now.AddDate(0, 0, int(n))
	_now = ToDate(_now)
	_now = _now.Add(time.Hour * time.Duration(hour))
	return _now.Sub(now)
}

func GetDurationToNextWeek(weekday time.Weekday, offset time.Duration) time.Duration {
	now := NowFromDefaultNower()
	next := GetNextWeek(weekday, now)
	return next.Add(offset).Sub(now)
}

func GetThisWeek(weekday time.Weekday, sometime time.Time) time.Time {
	n := sometime.Weekday() - weekday
	t := ToDate(sometime.Add(-time.Hour * 24 * time.Duration(n)))
	return t
}

func GetLastWeek(weekday time.Weekday, sometime time.Time) time.Time {
	this := GetThisWeek(weekday, sometime)
	return this.Add(-time.Hour * 24 * 7)
}

func GetNextWeek(weekday time.Weekday, sometime time.Time) time.Time {
	this := GetThisWeek(weekday, sometime)
	return this.Add(time.Hour * 24 * 7)
}

func GetOneWeekPeriodArray(sometime time.Time) (periodArray [7]Period) {
	start := GetThisWeek(time.Monday, sometime)
	for i := range periodArray {
		periodArray[i] = Period{
			Start: start.Add(time.Hour * 24 * time.Duration(i)),
			End:   start.Add(time.Hour * 24 * time.Duration(i+1)).Add(-1),
		}
	}
	return periodArray
}
