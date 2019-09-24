package udwTime

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"math"
	"strings"
	"time"
)

func MustDbTimeGetObjFromString(s string) time.Time {
	t, err := time.Parse(FormatMysqlNs, s)
	if err != nil {
		panic(err)
	}
	return t
}

func DbTimeGetObjFromString(s string) (t time.Time, err error) {
	return time.Parse(FormatMysqlNs, s)
}

func MustDbTimeGetObjFromStringIgnoreEmpty(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(FormatMysqlNs, s)
	if err != nil {
		panic(err)
	}
	return t
}

func DbTimeGetObjFromStringIgnoreError(s string) time.Time {
	t, err := time.ParseInLocation(FormatMysqlNs, s, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}

func MustDbTimeGetStringFromObj(t time.Time) string {
	return DbTimeGetStringFromObj(t)
}

func DbTimeGetStringFromObj(t time.Time) string {
	t = t.UTC()
	s := t.Format(FormatMysqlNs)
	index := strings.IndexByte(s, '.')

	if index == -1 {
		s += ".000000000"
	} else if index+10 > len(s) {
		s += strings.Repeat("0", index+10-len(s))
	}

	return s
}

func DbTimeGetStringFromObjV2(t time.Time) string {
	t = t.UTC()
	s := t.Format(FormatMysqlNsV3)
	index := strings.IndexByte(s, '.')
	if index == -1 {
		s += ".000000000"
	} else if index+10 > len(s) {
		s += strings.Repeat("0", index+10-len(s))
	}
	return s
}

func MustDbTimeSecondGetStringFromObj(t time.Time) string {
	s := t.UTC().Format(FormatDbTimeSecond)
	return s
}

func DbTimeGetMinValue() string {
	return "0001-01-01T00:00:00.000000000"
}

func DbTimeGetMaxValue() string {
	return "9999-12-31T23:59:59.999999999"
}

func MustDbTimeHourGetStringFromObj(t time.Time) string {
	s := t.UTC().Format(Iso3339Hour)
	return s
}

func MustDbTimeHourGetObjFromString(s string) time.Time {
	t, err := time.ParseInLocation(Iso3339Hour, s, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}

func MustDbTimeDescGetStringFromObj(t time.Time) string {
	t = t.UTC()
	y := t.Year()
	if y < 1970 {
		return "ffffffffffffffff"
	}
	if y > 2261 && t.After(getDbTimeDescMaxTime()) {
		return "8000000000000000"
	}
	value := t.UnixNano()
	i := uint64(math.MaxUint64 - uint64(value))
	return udwStrconv.FormatUint64HexPaddingWithZeroPrefix(i, 16)
}

func MustDbTimeDescGetObjFromString(s string) (t time.Time) {
	i := udwStrconv.MustParseUint64Hex(s)
	i64 := int64(math.MaxUint64 - uint64(i))
	if i64 < 0 {
		i64 = math.MaxInt64
	}
	return time.Unix(i64/1e9, i64%1e9).UTC()
}

func getDbTimeDescMaxTime() time.Time {
	return MustDbTimeGetObjFromString("2262-04-11T23:47:16.854775807")
}

func DbTimeGetUint64FromObjOrMax(t time.Time) uint64 {
	t = t.UTC()
	y := t.Year()
	if y < 1970 {
		return 0
	}
	if y >= 2554 && t.After(getDbTimeUint64MaxTime()) {
		return math.MaxUint64
	}
	return uint64(t.UnixNano())
}

func getDbTimeUint64MaxTime() time.Time {
	return MustDbTimeGetObjFromString("2554-07-21T23:34:33.709551615")
}

func DbTimeGetObjFromUint64(u uint64) (t time.Time) {
	return time.Unix(int64(u/1e9), int64(u%1e9)).UTC()
}

func DbTimeGetUint64SecondFromObj(t time.Time) uint64 {
	t = t.UTC()
	y := t.Year()
	if y < 1970 {
		return 0
	}
	u := uint64(t.Unix())
	if u >= uint64(253402300799) {
		u = 253402300799
	}
	return u
}

func DbTimeGetObjFromUint64Second(u uint64) (t time.Time) {
	if u > uint64(253402300799) {
		u = 253402300799
	}
	return time.Unix(int64(u), 0).UTC()
}
