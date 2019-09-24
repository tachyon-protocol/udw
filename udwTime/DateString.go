package udwTime

import "time"

func MustDateStringAddDay(s string, dayNum int) string {
	t, err := time.ParseInLocation(FormatDateMysql, s, time.UTC)
	if err != nil {
		panic(err)
	}
	return t.Add(time.Duration(dayNum) * Day).Format(FormatDateMysql)
}

func MustDateStringSubToDay(s1 string, s2 string) int {
	t1, err := time.ParseInLocation(FormatDateMysql, s1, time.UTC)
	if err != nil {
		panic(err)
	}
	t2, err := time.ParseInLocation(FormatDateMysql, s2, time.UTC)
	if err != nil {
		panic(err)
	}
	dur := t1.Sub(t2)
	return int(dur.Hours() / 24)
}

func DateStringToTimeInDefault(s string) (t time.Time, err error) {
	return time.ParseInLocation(FormatDateMysql, s, GetDefaultTimeZone())
}

func MustDateStringToTime(s string, tz *time.Location) (t time.Time) {
	t, err := time.ParseInLocation(FormatDateMysql, s, tz)
	if err != nil {
		panic(err)
	}
	return t
}
