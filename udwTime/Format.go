package udwTime

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"strconv"
	"strings"
	"time"
)

const (
	FormatMysqlZero  = "0000-00-00 00:00:00"
	FormatMysql      = "2006-01-02 15:04:05"
	FormatMysqlUs    = "2006-01-02 15:04:05.999999"
	FormatMysqlNs    = "2006-01-02T15:04:05.999999999"
	FormatMysqlNsV2  = "2006-01-02 15:04:05.999999999"
	FormatMysqlNsV3  = "2006-01-02T15-04-05.999999999"
	FormatFileName   = "2006-01-02_15-04-05"
	FormatFileNameV2 = "2006-01-02-15-04-05"

	FormatDateMysql = "2006-01-02"
	FormatMysqlDate = "2006-01-02"

	FormatZoneOffsetMysql   = "2006-01-02 15:04:05 -0700"
	FormatZoneOffsetMysqlV2 = "2006-01-02 15:04:05 -0700 CST"
	Iso3339Hour             = "2006-01-02T15"
	Iso3339Minute           = "2006-01-02T15:04"
	Iso3339Second           = "2006-01-02T15:04:05"
	FormatDbTimeSecond      = Iso3339Second
	AppleJsonFormat         = "2006-01-02 15:04:05 Etc/MST"
	AppleJsonFormatV2       = "2006-01-02 15:04:05 MST"
	Iso8601                 = "2006-01-02T15:04Z"
	Iso8601GMT              = "2006-01-02T15:04:05Z"

	FormatMysqlMinute           = "2006-01-02 15:04"
	FormatMysqlMouthAndDay      = "01-02"
	FormatMysqlYearAndMoney     = "2006-01"
	FormatInternational         = "Monday, 02 January 2006"
	FormatHourAndMinute         = "15:04"
	FormatMouthDayHourAndMinute = "01-02 15:04"

	FormatDateAndHour           = "2006-01-02 15"
	FormatUdwLog                = "20060102T15:04:05.000000-07"
	FormatDateTimeDigitalSecond = "20060102150405"
)

var ParseFormatGuessList = []string{
	FormatFileNameV2,
	FormatMysqlZero,
	FormatMysql,
	FormatDateMysql,
	FormatZoneOffsetMysqlV2,
	Iso3339Hour,
	Iso3339Minute,
	Iso3339Second,
	time.RFC3339,
	time.RFC3339Nano,
	Iso8601,
	Iso8601GMT,
}

var MysqlStart = "0000-01-01 00:00:00"
var MysqlEnd = "9999-12-31 23:59:59"

func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format(FormatMysql)
}

func DefaultFormatSecondV2(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format("20060102 15:04:05-07")
}

func DefaultMysqlFormat(t time.Time) string {
	return DefaultFormat(t)
}

func DefaultFormatNs(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format(FormatMysqlNsV2)
}

func DefaultFormatLocal(t time.Time) string {
	return t.In(time.Local).Format(time.RFC3339)
}

func DefaultDateFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format(FormatDateMysql)
}

func PstTimeZoneDateFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetUtcD8Zone()).Format(FormatDateMysql)
}

func MonthAndDayFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format(FormatMysqlMouthAndDay)
}

func YearMonthFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetDefaultTimeZone()).Format(FormatMysqlYearAndMoney)
}

func HourAndMinuteFormatLocal(t time.Time) string {
	return t.In(time.Local).Format(FormatHourAndMinute)
}

func NowWithFileNameFormatV2() string {
	return NowFromDefaultNower().Format(FormatFileNameV2)
}

func MonthDayYearFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	s := t.In(GetDefaultTimeZone()).Format(FormatDateMysql)
	list := strings.Split(s, "-")
	return list[1] + "/" + list[2] + "/" + list[0]
}

func MustSplitMysqlDateFormatPrefix(s string) (datePart string, remainPart string) {
	if len(s) < 10 {
		panic("[MustSplitMysqlDateFormatPrefix] len(s)<10 " + strconv.Itoa(len(s)))
	}
	datePart = s[:10]
	_, err := time.Parse(FormatDateMysql, datePart)
	if err != nil {
		panic("[MustSplitMysqlDateFormatPrefix] time.Parse " + err.Error())
	}
	return datePart, s[10:]
}

func DefaultFormatNsFixSize(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	s := t.In(GetDefaultTimeZone()).Format(FormatMysqlNsV2)
	if len(s) < 29 {
		if len(s) == 19 {
			s += "."
		}
		s += strings.Repeat("0", 29-len(s))
	}
	return s
}

func MustDateMysqlFormat(t string) {
	_, err := time.Parse(FormatDateMysql, t)
	if err != nil {
		panic("tfmhzkxx64 date format should like " + FormatDateMysql + " " + t)
	}
}

func ParseFormatDateMysqlInDefaultTz(tS string) (t time.Time, errMsg string) {
	t, err := time.ParseInLocation(FormatDateMysql, tS, GetDefaultTimeZone())
	if err != nil {
		return t, "jqbnayunze date format should like " + FormatDateMysql + " [" + tS + "]"
	}
	return t, ""
}

func MustParseFormatDateMysqlInDefaultTz(tS string) (t time.Time) {
	t, errMsg := ParseFormatDateMysqlInDefaultTz(tS)
	if errMsg != "" {
		panic(errMsg)
	}
	return t
}

func MustParseFormatDateMysqlInTz(tS string, tz *time.Location) (t time.Time) {
	t, err := time.ParseInLocation(FormatDateMysql, tS, tz)
	if err != nil {
		panic("jqbnayunze date format should like " + FormatDateMysql + " [" + tS + "]")
	}
	return t
}

func FormatRfc3339NanoNoTz(t time.Time) string {
	s := t.Format("2006-01-02T15:04:05.999999999")
	if len(s) < 29 {
		if len(s) == 19 {
			s += "."
		}
		s += strings.Repeat("0", 29-len(s))
	}
	return s
}

func FormatDefaultRfc3339(t time.Time) string {
	return t.In(GetDefaultTimeZone()).Format(time.RFC3339)
}

func FormatMysqlMinuteInTz(t time.Time, tz *time.Location) string {
	t = t.In(time.Local)
	s := t.Format(FormatMysqlMinute)
	_, offset := t.Zone()
	s += ` (UTC`
	if offset == 0 {
		return s + ")"
	} else if offset >= 0 {
		s += "+"
	} else {
		offset = -offset
		s += "-"
	}
	tzHour := offset / 3600
	tzMin := (offset - tzHour*3600) / 60
	s += udwStrconv.FormatIntPaddingWithZeroPre(tzHour, 2) + ":" + udwStrconv.FormatIntPaddingWithZeroPre(tzMin, 2) + ")"
	return s
}

func DateFormatInUtc2(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(GetUtc2Zone()).Format(FormatDateMysql)
}

func UdwLogFormat(t time.Time) string {
	return t.In(GetDefaultTimeZone()).Format(FormatUdwLog)
}
