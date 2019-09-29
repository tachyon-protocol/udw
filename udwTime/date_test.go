package udwTime

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestToMonth(ot *testing.T) {
	t := time.Date(2000, 1, 10, 0, 0, 0, 0, GetDefaultTimeZone())
	outT := ToMonth(t, GetDefaultTimeZone())
	y, m, d := outT.Date()
	udwTest.Equal(y, 2000)
	udwTest.Equal(int(m), 1)
	udwTest.Equal(d, 1)
	udwTest.Equal(outT.Hour(), 0)
	udwTest.Equal(outT.Minute(), 0)
	udwTest.Equal(outT.Second(), 0)
	udwTest.Equal(outT.Nanosecond(), 0)

	t = time.Date(2000, 1, 10, 0, 0, 0, 0, GetDefaultTimeZone())
	outT = ToMonthWithOffset(t, GetDefaultTimeZone(), 1)
	udwTest.Equal(outT.Year(), 2000)
	udwTest.Equal(int(outT.Month()), 2)
	udwTest.Equal(outT.Day(), 1)

	outT = ToMonthWithOffset(t, GetDefaultTimeZone(), 12)
	udwTest.Equal(outT.Year(), 2001)
	udwTest.Equal(int(outT.Month()), 1)
	udwTest.Equal(outT.Day(), 1)

	outT = ToMonthWithOffset(t, GetDefaultTimeZone(), 25)
	udwTest.Equal(outT.Year(), 2002)
	udwTest.Equal(int(outT.Month()), 2)
	udwTest.Equal(outT.Day(), 1)

	outT = ToMonthWithOffset(t, GetDefaultTimeZone(), -1)
	udwTest.Equal(outT.Year(), 1999)
	udwTest.Equal(int(outT.Month()), 12)
	udwTest.Equal(outT.Day(), 1)

	outT = ToMonthWithOffset(t, GetDefaultTimeZone(), -12)
	udwTest.Equal(outT.Year(), 1999)
	udwTest.Equal(int(outT.Month()), 1)
	udwTest.Equal(outT.Day(), 1)
}

func TestNextTime(t *testing.T) {
	now := "2017-07-26 22:21:50"
	SetDefaultNowerToFixTimeString(now)
	far := GetNextTimeDuration(22, 22)
	udwTest.Ok(far == 10*time.Second)
	far = GetNextTimeDuration(22, 21)
	fmt.Println(far)
	udwTest.Ok(far == (23*time.Hour + 59*time.Minute + 10*time.Second))
}

func TestSameMonth(ot *testing.T) {
	t := time.Date(2000, 8, 25, 0, 0, 0, 0, GetDefaultTimeZone())
	t1 := time.Date(2010, 8, 25, 0, 0, 0, 0, GetDefaultTimeZone())
	t2 := time.Date(2000, 8, 10, 0, 0, 0, 0, GetDefaultTimeZone())
	sameMonth := IsSameMonth(t, time.Now(), GetDefaultTimeZone())
	udwTest.Equal(sameMonth, false)
	sameMonth = IsSameMonth(t, t1, GetDefaultTimeZone())
	udwTest.Equal(sameMonth, false)
	sameMonth = IsSameMonth(t, t2, GetDefaultTimeZone())
	udwTest.Equal(sameMonth, true)
}

func TestIsSameHour(t *testing.T) {
	t1 := MustParseAutoInDefault(`2018-10-11 01:12:12`)
	t2 := MustParseAutoInDefault(`2018-10-11 01:15:12`)
	udwTest.Ok(`2018-10-11 01` == t1.Format(FormatDateAndHour))
	udwTest.Ok(`2018-10-11 01` == t2.Format(FormatDateAndHour))
	udwTest.Ok(IsSameHour(t1, t2, time.UTC))
}

func TestMonthLeftPercent(t *testing.T) {
	t1 := MustParseAutoInDefault(`2018-10-11 01:12:12`)
	udwTest.Equal(CountMonthLeftDay(t1, GetDefaultTimeZone()), 21)
	udwTest.Equal(MonthLeftPercent(t1, GetDefaultTimeZone()), 21.0/31.0)
	t1 = MustParseAutoInDefault(`2018-09-11 01:12:12`)
	udwTest.Equal(CountMonthLeftDay(t1, GetDefaultTimeZone()), 20)
	udwTest.Equal(MonthLeftPercent(t1, GetDefaultTimeZone()), 2.0/3.0)
}
