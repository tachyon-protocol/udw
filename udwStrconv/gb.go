package udwStrconv

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TB = 1024 * 1024 * 1024 * 1024
const GB = 1024 * 1024 * 1024
const MB = 1024 * 1024
const KB = 1024

func GbPaddingFromFloat64(byteNum float64) string {
	s := GbFromFloat64(byteNum)
	if len(s) < 8 {
		s = strings.Repeat(" ", 8-len(s)) + s
	}
	return s
}
func GbFromFloat64(byteNum float64) string {
	if byteNum >= 1e15 || byteNum <= -1e15 {
		return FormatFloat64ToFInLen(byteNum/(1024*1024*1024*1024*1024), 6) + "PB"
	}
	if byteNum >= 1e12 || byteNum <= -1e12 {
		return FormatFloat64ToFInLen(byteNum/(1024*1024*1024*1024), 6) + "TB"
	}
	if byteNum >= 1e9 || byteNum <= -1e9 {
		return FormatFloat64ToFInLen(byteNum/(1024*1024*1024), 6) + "GB"
	}
	if byteNum >= 1e6 || byteNum <= -1e6 {
		return FormatFloat64ToFInLen(byteNum/(1024*1024), 6) + "MB"
	}
	if byteNum >= 1e3 || byteNum <= -1e3 {
		return FormatFloat64ToFInLen(byteNum/(1024), 6) + "KB"
	}
	return FormatFloat64ToFInLen(byteNum, 7) + "B"
}

func GbFromFloat64WithUnit(byteNum float64) (string, string) {
	if byteNum >= 1e15 || byteNum <= -1e15 {
		return strconv.FormatFloat(byteNum/(1024*1024*1024*1024*1024), 'f', 2, 64), "PB"
	}
	if byteNum >= 1e12 || byteNum <= -1e12 {
		return strconv.FormatFloat(byteNum/(1024*1024*1024*1024), 'f', 2, 64), "TB"
	}
	if byteNum >= 1e9 || byteNum <= -1e9 {
		return strconv.FormatFloat(byteNum/(1024*1024*1024), 'f', 2, 64), "GB"
	}
	if byteNum >= 1e6 || byteNum <= -1e6 {
		return strconv.FormatFloat(byteNum/(1024*1024), 'f', 2, 64), "MB"
	}
	if byteNum >= 1e3 || byteNum <= -1e3 {
		return strconv.FormatFloat(byteNum/(1024), 'f', 2, 64), "KB"
	}
	if byteNum >= 1 || byteNum <= -1 {
		s1 := strconv.FormatFloat(byteNum, 'f', 2, 64)
		s2 := strconv.FormatFloat(byteNum, 'f', -1, 64)
		if len(s1) < len(s2) {
			return s1, "B"
		} else {
			return s2, "B"
		}
	}
	return strconv.FormatFloat(byteNum, 'f', -1, 64), "B"
}

func GbFromInt64(byteNum int64) string {
	return GbFromFloat64(float64(byteNum))
}
func GbFromUint64(byteNum uint64) string {
	return GbFromFloat64(float64(byteNum))
}

func GbFromInt(byteNum int) string {
	return GbFromFloat64(float64(byteNum))
}
func GbPaddingFromInt(byteNum int) string {
	return GbPaddingFromInt64(int64(byteNum))
}
func GbPaddingFromInt64(byteNum int64) string {
	return GbPaddingFromFloat64(float64(byteNum))
}

func GbSpeedFromInt(i int) string {
	return GbFromInt(i) + "/s"
}

func GbSpeedFromFloat64(f float64) string {
	return GbFromFloat64(f) + "/s"
}

func GbSpeedFromFloat64AndDuration(f float64, dur time.Duration) string {
	bytePerSecond := float64(f) / (float64(dur) / float64(time.Second))
	return GbFromFloat64(bytePerSecond) + "/s"
}

func SizeStringWithFloat64(bytePerSecond float64) string {
	return GbFromFloat64(bytePerSecond)
}

func SizeStringWithoutUnit(bytePerSecond float64) string {
	if bytePerSecond > 1e9 || bytePerSecond < -1e9 {
		return fmt.Sprintf("%.2f ", bytePerSecond/(1024*1024*1024))
	}
	if bytePerSecond > 1e6 || bytePerSecond < -1e6 {
		return fmt.Sprintf("%.2f ", bytePerSecond/(1024*1024))
	}
	if bytePerSecond > 1e3 || bytePerSecond < -1e3 {
		return fmt.Sprintf("%.2f ", bytePerSecond/1024)
	}
	return fmt.Sprintf("%.2f ", bytePerSecond)
}

func SizeString(byteNum int64) string {
	return GbFromInt64(byteNum)
}

func SizeStringWithPadding(byteNum int64) string {
	return GbPaddingFromInt64(byteNum)
}

func Float64ToFlowString(in float64) string {
	return GbFromFloat64(in)
}

func GbStringToFloat64Default0(in string) float64 {
	f, errMsg := GbstringToFloat64(in)
	if errMsg != "" {
		return 0
	}
	return f

}

func GbstringToFloat64(in string) (f float64, errMsg string) {
	a := float64(0)
	in = strings.TrimSpace(in)
	in = strings.ToUpper(in)
	if len(in) <= 1 {
		return 0, "[GbstringToFloat64] format error len(in)<=1"
	}
	b1 := in[len(in)-1]
	if b1 != 'B' {
		return 0, `[GbstringToFloat64] pe!="B"`
	}
	b2 := in[len(in)-2]
	var b3 string
	switch b2 {
	case 'K':
		a = 1 << 10
		b3 = in[:len(in)-2]
	case 'M':
		a = 1 << 20
		b3 = in[:len(in)-2]
	case 'G':
		a = 1 << 30
		b3 = in[:len(in)-2]
	case 'T':
		a = 1 << 40
		b3 = in[:len(in)-2]
	case 'P':
		a = 1 << 50
		b3 = in[:len(in)-2]
	default:
		a = 1
		b3 = in[:len(in)-1]
	}
	f, err := ParseFloat64(b3)
	if err != nil {
		return 0, err.Error()
	}
	return f * a, ""
}

func FormatGbAndIntFromInt64(i int64) string {
	return GbFromInt64(i) + "(" + FormatInt64(i) + ")"
}

func FormatGbAndIntFromInt(i int) string {
	return FormatGbAndIntFromInt64(int64(i))
}
