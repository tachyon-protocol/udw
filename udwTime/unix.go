package udwTime

import (
	"math"
	"time"
)

func FromUnixNano(f uint64) time.Time {
	return time.Unix(int64(f/1e9), int64(f%1e9)).In(GetDefaultTimeZone())
}

func FromUnixNanoInt64(f int64) time.Time {
	return time.Unix(int64(f/1e9), int64(f%1e9)).In(GetDefaultTimeZone())
}

func FromUnixMillisecondsInt64(f int64) time.Time {
	return time.Unix(0, f*int64(time.Millisecond)).In(GetDefaultTimeZone())
}

func UnixNanoNow() int64 {
	return time.Now().UnixNano()
}

func GetUnixFloat(t1 time.Time) float64 {
	return (float64(t1.Nanosecond()) / 1e9) + float64(t1.Unix())
}

func FromUnixFloat(f float64) time.Time {
	s, ns := math.Modf(f)
	return time.Unix(int64(s), int64(ns*1e9)).In(GetDefaultTimeZone())
}

func FromUnixInt64(f int64) time.Time {
	return time.Unix(int64(f), 0).In(GetDefaultTimeZone())
}
