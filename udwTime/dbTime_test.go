package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"testing"
	"time"
)

func TestDbTime(ot *testing.T) {
	var t time.Time
	t = time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	udwTest.Equal(MustDbTimeGetStringFromObj(t), "2000-01-01T01:01:01.000000001")
	t = time.Date(2000, 1, 1, 1, 1, 1, 1, GetBeijingZone())
	udwTest.Equal(MustDbTimeGetStringFromObj(t), "1999-12-31T17:01:01.000000001")
	t = time.Date(2000, 1, 1, 1, 1, 1, 10, time.UTC)
	udwTest.Equal(MustDbTimeGetStringFromObj(t), "2000-01-01T01:01:01.000000010")

	t = MustDbTimeGetObjFromString("2000-01-01T01:01:01.000000010")
	udwTest.Equal(t.Location().String(), "UTC")
	udwTest.Equal(t.Nanosecond(), 10)

	udwTest.Equal(MustDbTimeGetStringFromObj(time.Time{}), "0001-01-01T00:00:00.000000000")

	udwTest.AssertPanic(func() {
		MustDbTimeGetObjFromString("")
	})
	udwTest.Ok(MustDbTimeGetObjFromStringIgnoreEmpty("").IsZero())

	udwTest.Equal(MustDbTimeGetStringFromObj(time.Date(10000, 1, 1, 1, 1, 1, 10, time.UTC)), "10000-01-01T01:01:01.000000010")
}

func TestDbTimeDesc(ot *testing.T) {
	convFromDbTime := func(s string) string {
		return MustDbTimeDescGetStringFromObj(MustDbTimeGetObjFromString(s))
	}
	udwTest.Equal(convFromDbTime("2000-01-01T01:01:01.000000010"), "f2dcafdc49223df5")
	udwTest.Equal(convFromDbTime("2000-01-01T01:01:01.000000011"), "f2dcafdc49223df4")
	udwTest.Ok(convFromDbTime("2000-01-01T01:01:01.000000011") < convFromDbTime("2000-01-01T01:01:01.000000010"))
	udwTest.Equal(convFromDbTime("0001-01-01T00:00:00.000000000"), "ffffffffffffffff")
	udwTest.Equal(convFromDbTime("1969-12-31T23:59:59.999999999"), "ffffffffffffffff")
	udwTest.Equal(convFromDbTime("1970-01-01T00:00:00.000000000"), "ffffffffffffffff")
	udwTest.Equal(convFromDbTime("1970-01-01T00:00:00.000000001"), "fffffffffffffffe")
	udwTest.Equal(convFromDbTime("2261-01-01T00:00:00.000000000"), "808f09bed2cdffff")
	udwTest.Equal(convFromDbTime("2261-12-31T23:59:59.999999999"), "801effeba52b0000")
	udwTest.Equal(convFromDbTime("2262-01-01T00:00:00.000000000"), "801effeba52affff")
	udwTest.Equal(convFromDbTime("2262-04-11T23:47:16.854775807"), "8000000000000000")
	udwTest.Equal(convFromDbTime("2262-04-11T23:47:16.854775806"), "8000000000000001")
	udwTest.Equal(convFromDbTime("2553-01-01T00:00:00.000000000"), "8000000000000000")
	udwTest.Equal(convFromDbTime("2554-01-01T00:00:00.000000000"), "8000000000000000")

	convToDbTime := func(s string) string {
		return MustDbTimeGetStringFromObj(MustDbTimeDescGetObjFromString(s))
	}
	udwTest.Equal(convToDbTime("ffffffffffffffff"), "1970-01-01T00:00:00.000000000")
	udwTest.Equal(convToDbTime("f2dcafdc49223df4"), "2000-01-01T01:01:01.000000011")
	udwTest.Equal(convToDbTime("8000000000000000"), "2262-04-11T23:47:16.854775807")
	udwTest.Equal(convToDbTime("8000000000000001"), "2262-04-11T23:47:16.854775806")
	udwTest.Equal(convToDbTime("0000000000000000"), "2262-04-11T23:47:16.854775807")

}

func TestDbTimeGetUint64FromObjOrMax(ot *testing.T) {
	convFromDbTimeS := func(s string) uint64 {
		return DbTimeGetUint64FromObjOrMax(MustDbTimeGetObjFromString(s))
	}
	convToDbTimeS := func(u uint64) string {
		return MustDbTimeGetStringFromObj(DbTimeGetObjFromUint64(u))
	}
	type tCas struct {
		s string
		u uint64
	}
	for _, cas := range []tCas{
		{"2000-01-01T01:01:01.000000010", uint64(946688461000000010)},
		{"1970-01-01T00:00:00.000000000", uint64(0)},
		{"1970-01-01T00:00:00.000000001", uint64(1)},
		{"2262-04-11T23:47:16.854775807", uint64(9223372036854775807)},
		{"2262-04-11T23:47:16.854775808", uint64(9223372036854775808)},
		{"2554-01-01T00:00:00.000000000", uint64(18429292800000000000)},
		{"2554-07-21T23:34:33.709551615", uint64(math.MaxUint64)},
	} {
		udwTest.Equal(convFromDbTimeS(cas.s), cas.u)
		udwTest.Equal(convToDbTimeS(cas.u), cas.s)
	}
	udwTest.Equal(convFromDbTimeS("1969-12-31T23:59:59.999999999"), uint64(0))
	udwTest.Equal(convFromDbTimeS("2555-07-21T23:34:33.709551615"), uint64(math.MaxUint64))
	udwTest.Equal(convFromDbTimeS("2554-07-21T23:34:33.709551616"), uint64(math.MaxUint64))
}

func TestDbTimeGetUint64SecondFromObj(ot *testing.T) {
	convFromDbTimeS := func(s string) uint64 {
		return DbTimeGetUint64SecondFromObj(MustDbTimeGetObjFromString(s))
	}
	convToDbTimeS := func(u uint64) string {
		return MustDbTimeGetStringFromObj(DbTimeGetObjFromUint64Second(u))
	}
	type tCas struct {
		s string
		u uint64
	}
	for _, cas := range []tCas{
		{"2000-01-01T01:01:01.000000000", uint64(946688461)},
		{"1970-01-01T00:00:00.000000000", uint64(0)},
		{"2262-04-11T23:47:16.000000000", uint64(9223372036)},
		{"2554-01-01T00:00:00.000000000", uint64(18429292800)},
		{"2554-07-21T23:34:33.000000000", uint64(18446744073)},
		{"9999-12-31T23:59:59.000000000", uint64(253402300799)},
	} {
		udwTest.Equal(convFromDbTimeS(cas.s), cas.u)
		udwTest.Equal(convToDbTimeS(cas.u), cas.s)
	}
	udwTest.Equal(convFromDbTimeS("1969-12-31T23:59:59.999999999"), uint64(0))
	udwTest.Equal(convToDbTimeS(253402300899), "9999-12-31T23:59:59.000000000")
	udwTest.Equal(convToDbTimeS(math.MaxInt64), "9999-12-31T23:59:59.000000000")
	udwTest.Equal(convToDbTimeS(math.MaxUint64), "9999-12-31T23:59:59.000000000")
}
