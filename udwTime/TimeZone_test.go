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
}
