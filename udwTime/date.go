package udwTime

import (
	"time"
)

func ToDateString(t time.Time) string {
	return t.Format(FormatDateMysql)
}

func ToDateStringInDefaultTz(t time.Time) string {
	return t.In(GetDefaultTimeZone()).Format(FormatDateMysql)
}

func ToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func ToDateDefault(t time.Time) time.Time {
	y, m, d := t.In(GetDefaultTimeZone()).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, GetDefaultTimeZone())
}

func DateSub(t1 time.Time, t2 time.Time, loc *time.Location) time.Duration {
	return ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
}

func DateSubToDay(t1 time.Time, t2 time.Time, loc *time.Location) int {
	dur := ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
	return int(dur.Hours() / 24)
}

func DateSubToHour(t1 time.Time, t2 time.Time, loc *time.Location) int {
	dur := ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
	return int(dur.Hours())
}

func DateSubLocal(t1 time.Time, t2 time.Time) time.Duration {
	return DateSub(t1, t2, time.Local)
}

func IsSameDay(t1 time.Time, t2 time.Time, loc *time.Location) bool {
	return DateSub(t1, t2, loc) == 0
}

func IsSameHour(t1 time.Time, t2 time.Time, loc *time.Location) bool {
	return t1.In(loc).Format(FormatDateAndHour) == t2.In(loc).Format(FormatDateAndHour)
}

func IsSameMonth(t1 time.Time, t2 time.Time, loc *time.Location) bool {
	return ToMonth(t1, loc).Sub(ToMonth(t2, loc)) == 0
}

func ToMonth(t time.Time, loc *time.Location) time.Time {
	y, m, _ := t.In(loc).Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, loc)
}

func ToMonthWithOffset(t time.Time, loc *time.Location, offset int) time.Time {
	y, m, _ := t.In(loc).Date()
	return time.Date(y, m+time.Month(offset), 1, 0, 0, 0, 0, loc)
}

func CountMonthLeftDay(t time.Time, loc *time.Location) int {
	count := 0
	ot := t
	for {
		if !IsSameMonth(ot, t, loc) {
			break
		}
		count++
		ot = ot.Add(Day)
	}
	return count
}

func MonthLeftPercent(t time.Time, loc *time.Location) float64 {
	count := CountMonthLeftDay(t, loc)
	passed := t.In(loc).Day() - 1
	return float64(count) / float64(count+passed)
}
