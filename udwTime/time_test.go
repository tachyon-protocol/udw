package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustFromMysqlFormatDefaultTZ(ot *testing.T) {
	t := MustFromMysqlFormatDefaultTZ("2001-01-01 00:00:00")
	udwTest.Equal(t.Hour(), 0)
	udwTest.Equal(t.Day(), 1)

	t = MustFromMysqlFormatDefaultTZ("0000-00-00 00:00:00")
	udwTest.Equal(t.IsZero(), true)
}

func TestFixLocalTimeToOffsetSpecifiedZoneTime(ot *testing.T) {
	localTime := "2015-11-12 14:15:55"
	otherZoneTime := FixLocalTimeToOffsetSpecifiedZoneTime(3600, localTime)
	udwTest.Equal(otherZoneTime, "2015-11-12 07:15:55 +0100")
	otherZoneTime = FixLocalTimeToOffsetSpecifiedZoneTime(7200, localTime)
	udwTest.Equal(otherZoneTime, "2015-11-12 08:15:55 +0200")
	otherZoneTime = FixLocalTimeToOffsetSpecifiedZoneTime(-18000, localTime)
	udwTest.Equal(otherZoneTime, "2015-11-12 01:15:55 -0500")
	otherZoneTime = FixLocalTimeToOffsetSpecifiedZoneTime(43200, localTime)
	udwTest.Equal(otherZoneTime, "2015-11-12 18:15:55 +1200")
}
