package udwTime

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwStrconv"
	"time"
)

func DurationFormatTimeMysql(dur time.Duration) string {
	isNeg := (dur <= -time.Second)
	if dur < 0 {
		dur = -dur
	}
	s := fmt.Sprintf("%02d:%02d:%02d",
		int(dur/time.Hour),
		int((dur%time.Hour)/time.Minute),
		int((dur%time.Minute)/time.Second),
	)
	if isNeg {
		s = "-" + s
	}
	return s
}

func DurationFormat(dur time.Duration) string {
	return DurationFormatFloat64Ns(float64(dur))
}

func DurationFormatPadding(dur time.Duration) string {
	return DurationFormatFloat64Ns(float64(dur))
}

func DurationFormatFloat64Seconds(dur float64) string {
	return DurationFormatFloat64Ns(dur * 1e9)
}

func DurationFormatFloat64Ns(dur float64) string {
	const day = 24 * time.Hour
	const year = 365 * day
	if (dur >= float64(year)) || (dur <= float64(-year)) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(year), 6) + "y"
	} else if dur >= float64(day) || dur <= float64(-day) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(day), 6) + "d"
	} else if dur >= float64(time.Hour) || dur < float64(-time.Hour) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(time.Hour), 6) + "h"
	} else if dur >= float64(time.Minute) || dur <= float64(-time.Minute) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(time.Minute), 6) + "m"
	} else if dur >= float64(time.Second) || dur <= float64(-time.Second) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(time.Second), 6) + "s"
	} else if dur >= float64(time.Millisecond) || dur <= float64(-time.Millisecond) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(time.Millisecond), 5) + "ms"
	} else if dur >= float64(time.Microsecond) || dur <= float64(-time.Microsecond) {
		return udwStrconv.FormatFloat64ToFInLen(float64(dur)/float64(time.Microsecond), 5) + "us"
	} else {
		return udwStrconv.FormatFloat64ToFInLen(dur, 5) + "ns"
	}
}

func DurationFormatBefore(dur time.Duration) string {
	durInt := int(dur.Seconds())
	if durInt < 120 {
		return "1 minute ago"
	}
	if durInt >= 120 && durInt < 3600 {
		return fmt.Sprintf("%s minutes ago", udwStrconv.FormatInt(durInt/60))
	}
	if durInt >= 3600 && durInt < 24*3600 {
		h := durInt / 3600
		if h == 1 {
			return "1 hour ago"
		} else {
			return fmt.Sprintf("%s hours ago", udwStrconv.FormatInt(h))
		}
	}
	if durInt >= 24*3600 {
		d := durInt / (24 * 3600)
		if d == 1 {
			return "1 day ago"
		} else {
			return fmt.Sprintf("%s days ago", udwStrconv.FormatInt(d))
		}
	}
	return ""
}

func DurationFormatByHourMin(dur time.Duration) string {
	mins := dur.Minutes()
	hours := int(mins) / 60
	remainMins := int(mins) % 60
	return fmt.Sprintf("%02d:%02d", hours, remainMins)
}
