package udwStrconv_test

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"sort"
	"strconv"
	"testing"
)

func TestFormatUint64HexPadding(t *testing.T) {
	for _, i := range []uint64{
		5,
		(1 << 64) - 1,
		1 << 63,
	} {
		s := udwStrconv.FormatUint64HexPadding(i)
		udwTest.Equal(udwStrconv.MustParseUint64Hex(s), i)
	}
}

func TestFormatUint64Padding(ot *testing.T) {
	for _, i := range []uint64{
		0,
		5,
		1 << 32,
		(1 << 64) - 1,
		1 << 63,
	} {
		s := udwStrconv.FormatUint64Padding(i)
		udwTest.Equal(len(s), 20)
		i1, err := strconv.ParseUint(s, 10, 64)
		udwTest.Equal(err, nil)
		udwTest.Equal(i1, i)

	}
}

func TestFormatIntPaddingWithZeroPre(t *testing.T) {
	for _, cas := range []struct {
		i      int
		width  int
		expect string
	}{
		{0, 4, "0000"},
		{123, 4, "0123"},
		{1234, 4, "1234"},
		{12345, 4, "12345"},
	} {
		s := udwStrconv.FormatIntPaddingWithZeroPre(cas.i, cas.width)
		udwTest.Equal(s, cas.expect)
	}
	list := []string{}
	for i := 0; i < 100; i++ {
		s := udwStrconv.FormatIntPaddingWithZeroPre(i, 3)
		list = append(list, s)
	}
	sort.Strings(list)
	for i := 0; i < 99; i++ {
		a1 := udwStrconv.MustParseInt(list[i])
		a2 := udwStrconv.MustParseInt(list[i+1])
		udwTest.Ok(a1 < a2)
	}
	{
		list := []string{}
		for i := 0; i < 100; i++ {
			s := udwStrconv.FormatUint64HexPaddingWithZeroPrefix(uint64(i), 4)
			list = append(list, s)
		}
		sort.Strings(list)
		for i := 0; i < 99; i++ {
			a1 := udwStrconv.MustParseIntHex(list[i])
			a2 := udwStrconv.MustParseIntHex(list[i+1])
			udwTest.Ok(a1 < a2)
		}
	}
}

func TestGetPercent(t *testing.T) {
	udwTest.Ok(`0%` == udwStrconv.GetPercent(0, 0))
	udwTest.Ok(`10.00%` == udwStrconv.GetPercent(10, 100))
}

func TestFormatFloatPercentPaddingPrec4(t *testing.T) {
	for _, cas := range []float64{
		0.1, 1, 0, 0.1111,
	} {
		udwTest.Equal(len(udwStrconv.FormatFloatPercentPaddingPrec4(cas)), 9)
	}
}
