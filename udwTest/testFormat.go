package udwTest

import (
	"strconv"
	"time"
)

func durationFormatFloat64Ns(dur float64) string {
	const day = 24 * time.Hour
	const year = 365 * day
	if (dur >= float64(year)) || (dur <= float64(-year)) {
		return formatFloat64ToFInLen(float64(dur)/float64(year), 6) + "y"
	} else if dur >= float64(day) || dur <= float64(-day) {
		return formatFloat64ToFInLen(float64(dur)/float64(day), 6) + "d"
	} else if dur >= float64(time.Hour) || dur < float64(-time.Hour) {
		return formatFloat64ToFInLen(float64(dur)/float64(time.Hour), 6) + "h"
	} else if dur >= float64(time.Minute) || dur <= float64(-time.Minute) {
		return formatFloat64ToFInLen(float64(dur)/float64(time.Minute), 6) + "m"
	} else if dur >= float64(time.Second) || dur <= float64(-time.Second) {
		return formatFloat64ToFInLen(float64(dur)/float64(time.Second), 6) + "s"
	} else if dur >= float64(time.Millisecond) || dur <= float64(-time.Millisecond) {
		return formatFloat64ToFInLen(float64(dur)/float64(time.Millisecond), 5) + "ms"
	} else if dur >= float64(time.Microsecond) || dur <= float64(-time.Microsecond) {
		return formatFloat64ToFInLen(float64(dur)/float64(time.Microsecond), 5) + "us"
	} else {
		return formatFloat64ToFInLen(dur, 5) + "ns"
	}
}

func gbFromFloat64(byteNum float64) string {
	if byteNum >= 1e15 || byteNum <= -1e15 {
		return formatFloat64ToFInLen(byteNum/(1024*1024*1024*1024*1024), 5) + "PB"
	}
	if byteNum >= 1e12 || byteNum <= -1e12 {
		return formatFloat64ToFInLen(byteNum/(1024*1024*1024*1024), 5) + "TB"
	}
	if byteNum >= 1e9 || byteNum <= -1e9 {
		return formatFloat64ToFInLen(byteNum/(1024*1024*1024), 5) + "GB"
	}
	if byteNum >= 1e6 || byteNum <= -1e6 {
		return formatFloat64ToFInLen(byteNum/(1024*1024), 5) + "MB"
	}
	if byteNum >= 1e3 || byteNum <= -1e3 {
		return formatFloat64ToFInLen(byteNum/(1024), 5) + "KB"
	}
	return formatFloat64ToFInLen(byteNum, 6) + "B"
}

func formatFloat64ToFInLen(f float64, showLen int) string {
	s1 := strconv.FormatFloat(f, 'f', 0, 64)
	if len(s1)+1 >= showLen {
		if len(s1) == showLen {
			return s1
		} else {
			return "0" + s1
		}
	}
	return strconv.FormatFloat(f, 'f', showLen-len(s1)-1, 64)
}
