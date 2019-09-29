package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
	"time"
)

func TestGetOneWeekPeriodArray(t *testing.T) {
	pa := GetOneWeekPeriodArray(MustParseAutoInDefault("2018-06-01 12:12:12"))
	for _, p := range pa {
		udwTest.Ok(!p.Start.IsZero())
		udwTest.Ok(!p.End.IsZero())
	}
	target := ToDate(MustParseAutoInDefault("2018-05-30 00:00:00")).Add(-1)
	pa[1].End.Equal(target)
}

func TestGetDurationToNextWeek(t *testing.T) {
	SetDefaultNowerToFixTimeString("2018-07-18 00:00:00")
	dur := GetDurationToNextWeek(time.Monday, 3*time.Hour)
	udwTest.Ok(dur == time.Hour*24*5+time.Hour*3)
}

func TestGetThisWeek(t *testing.T) {
	thisMonday := GetThisWeek(time.Monday, MustParseAutoInDefault("2018-07-18 12:12:12"))
	udwTest.Ok(thisMonday.Weekday() == time.Monday)
	udwTest.Ok(DefaultFormat(thisMonday) == "2018-07-16 00:00:00")
}

func TestGetLastWeek(t *testing.T) {
	thisMonday := GetLastWeek(time.Monday, MustParseAutoInDefault("2018-07-18 12:12:12"))
	udwTest.Ok(thisMonday.Weekday() == time.Monday)
	udwTest.Ok(DefaultFormat(thisMonday) == "2018-07-09 00:00:00")
}
