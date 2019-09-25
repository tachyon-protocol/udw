package udwTime

import (
	"fmt"

	"strconv"
	"strings"
	"time"
)

const (
	Day        = 24 * time.Hour
	Month      = 30 * Day
	Year       = 365 * Day
	YearSecond = 365 * 24 * 60 * 60
	WeekSecond = 7 * 24 * 60 * 60
)

func ParseAutoInLocation(sTime string, loc *time.Location) (t time.Time, err error) {
	switch sTime {
	case "0000-00-00", "0000-00-00 00:00:00":
		return time.Time{}, nil
	}
	for _, format := range ParseFormatGuessList {
		t, err = time.ParseInLocation(format, sTime, loc)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("[ParseAutoInLocation] time: %s can not parse", sTime)
	return
}

func FixLocalTimeToOffsetSpecifiedZoneTime(timeOffset int, localTime string) string {
	clientTimeZone := time.FixedZone("ClientTimeZone", timeOffset)
	serverTime := MustParseAutoInDefault(localTime)
	return serverTime.In(clientTimeZone).Format(FormatZoneOffsetMysql)
}

func MustParseAutoInDefault(sTime string) (t time.Time) {
	t, err := ParseAutoInLocation(sTime, GetDefaultTimeZone())
	if err != nil {
		panic(err)
	}
	return t
}

func MustParseAutoInDefaultIgnoreEmpty(sTime string) (t time.Time) {
	if sTime == "" {
		return time.Time{}
	}
	return MustParseAutoInDefault(sTime)
}

func ParseAutoInDefault(sTime string) (t time.Time, err error) {
	return ParseAutoInLocation(sTime, GetDefaultTimeZone())
}

func MustParseAutoInUTC(sTime string) (t time.Time) {
	t, err := ParseAutoInLocation(sTime, time.UTC)
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromMysqlFormat(timeString string) time.Time {
	t, err := time.Parse(FormatMysql, timeString)
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromMysqlDateFormat(timeString string) time.Time {
	t, err := time.Parse(FormatDateMysql, timeString)
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromMysqlFormatInLocation(timeString string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation(FormatMysql, timeString, loc)
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromMysqlFormatDefaultTZ(timeString string) time.Time {
	if timeString == "0000-00-00 00:00:00" {
		return time.Time{}
	}
	t, err := time.ParseInLocation(FormatMysql, timeString, GetDefaultTimeZone())
	if err != nil {
		panic(err)
	}
	return t
}

func MustFromLocalMysqlFormat(timeString string) time.Time {
	t, err := time.ParseInLocation(FormatMysql, timeString, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func ToLocal(t time.Time) time.Time {
	return t.Local()
}

func ModBySecond(t1 time.Time) time.Time {
	return t1.Truncate(time.Second)
}

func MysqlUsFormat(t time.Time) string {
	s := t.In(GetDefaultTimeZone()).Format(FormatMysqlUs)
	if len(s) < 26 {
		s += strings.Repeat("0", 26-len(s))
	}
	return s
}

func TimeIntSecConvert(sec int) string {
	if sec < 60 {
		return strconv.Itoa(sec) + "s"
	}
	if sec > 60 && sec < 3600 {
		min := strconv.Itoa(sec/60) + "min"
		s := strconv.Itoa(sec%60) + "s"
		return min + s
	}
	if sec > 3600 {
		hour := strconv.Itoa(sec/3600) + "hour"
		min := strconv.Itoa((sec%3600)/60) + "min"
		s := strconv.Itoa((sec%3600)/60) + "s"
		return hour + min + s
	}
	return ""
}

func ParseMillStringToTime(mill string) (time.Time, error) {
	msInt, err := strconv.ParseInt(mill, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func GetSubSeconds(a time.Time, b time.Time) int64 {
	return a.Unix() - b.Unix()
}

func GetMaxTime() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 9999999, time.UTC)
}
