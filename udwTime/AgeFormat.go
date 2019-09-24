package udwTime

import "time"

func DefaultAgeFormat(targetTime time.Time, now time.Time) string {
	return DurationFormat(now.Sub(targetTime))
}

func DefaultTimeAndAgeFormat(targetTime time.Time, now time.Time) string {
	return DefaultFormat(targetTime) + "(" + DurationFormat(now.Sub(targetTime)) + ")"
}
