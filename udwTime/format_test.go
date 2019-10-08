package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestFormat(ot *testing.T) {

	t, err := time.Parse(AppleJsonFormat, "2014-04-16 18:26:18 Etc/GMT")
	udwTest.Equal(err, nil)
	udwTest.Equal(t, MustFromMysqlFormatInLocation("2014-04-16 18:26:18", time.FixedZone("GMT", 0)))

	udwTest.Equal(DefaultFormat(t), "2014-04-17 02:26:18")
	udwTest.Equal(MonthAndDayFormat(t), "04-17")
	SetDefaultNowerToRealTime()
	udwTest.Equal(NowWithFileNameFormatV2(), time.Now().In(GetDefaultTimeZone()).Format(FormatFileNameV2))

	udwTest.Equal(DefaultFormatNsFixSize(t), "2014-04-17 02:26:18.000000000")
	t2 := MustDbTimeGetObjFromString("2014-04-16T18:26:18.010730669")
	udwTest.Equal(DefaultFormatNsFixSize(t2), "2014-04-17 02:26:18.010730669")

	udwTest.Equal(FormatMysqlMinuteInTz(t, GetDefaultTimeZone()), "2014-04-17 02:26 (UTC+08:00)")
}
