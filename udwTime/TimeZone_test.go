package udwTime

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestFixTimeZone(ot *testing.T) {
	t := MustDbTimeGetObjFromString("2000-01-01T00:00:00.000000000")
	udwTest.Equal(t.In(GetUtc8Zone()).Hour(), 8)
	udwTest.Equal(t.In(GetUtc8Zone()).Format(FormatMysql), "2000-01-01 08:00:00")
	udwTest.Equal(t.In(GetUtc2Zone()).Format(FormatMysql), "2000-01-01 02:00:00")

	t = MustParseFormatDateMysqlInDefaultTz("2000-01-01")
	udwTest.Equal(t.In(GetUtc8Zone()).Format(FormatMysql), "2000-01-01 00:00:00")
	udwTest.Equal(t.In(GetUtc2Zone()).Format(FormatMysql), "1999-12-31 18:00:00")

	t = time.Date(2000, 1, 1, 0, 0, 0, 0, GetUtc2Zone())
	udwTest.Equal(t.In(GetUtc8Zone()).Format(FormatMysql), "2000-01-01 06:00:00")

	t, err := time.ParseInLocation(FormatDateMysql, "2000-01-01", GetUtc2Zone())
	udwErr.PanicIfError(err)
	udwTest.Equal(t.In(GetUtc8Zone()).Format(FormatMysql), "2000-01-01 06:00:00")

	t = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	udwTest.Equal(t.In(GetUtcD8Zone()).Format(time.RFC3339Nano), "1999-12-31T16:00:00-08:00")
	udwTest.Equal(t.In(GetUtcD7Zone()).Format(time.RFC3339Nano), "1999-12-31T17:00:00-07:00")
	udwTest.Equal(t.In(GetUtcD5Zone()).Format(time.RFC3339Nano), "1999-12-31T19:00:00-05:00")
	udwTest.Equal(t.In(GetUtcD4Zone()).Format(time.RFC3339Nano), "1999-12-31T20:00:00-04:00")
	udwTest.Equal(t.In(GetUtc2Zone()).Format(time.RFC3339Nano), "2000-01-01T02:00:00+02:00")
	udwTest.Equal(t.In(GetUtc8Zone()).Format(time.RFC3339Nano), "2000-01-01T08:00:00+08:00")
}

func TestFormatTimeZone(ot *testing.T) {
	udwTest.Equal(FormatTimeZone(GetUtcD8Zone()), "UTC-8(-08:00)")
	udwTest.Equal(FormatTimeZone(GetUtcD7Zone()), "UTC-7(-07:00)")
	udwTest.Equal(FormatTimeZone(GetUtcD5Zone()), "UTC-5(-05:00)")
	udwTest.Equal(FormatTimeZone(GetUtcD4Zone()), "UTC-4(-04:00)")
	udwTest.Equal(FormatTimeZone(GetUtc2Zone()), "UTC2(+02:00)")
	udwTest.Equal(FormatTimeZone(GetUtc8Zone()), "UTC8(+08:00)")
}
