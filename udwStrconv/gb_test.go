package udwStrconv_test

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"testing"
)

func TestGbFromFloat64(t *testing.T) {
	udwTest.Equal(udwStrconv.GbFromFloat64(0.5), "0.5000B")
	udwTest.Equal(udwStrconv.GbFromFloat64(158.33333333333334), "158.33B")
	udwTest.Equal(udwStrconv.GbFromFloat64(1), "1.0000B")
	udwTest.Equal(udwStrconv.GbFromFloat64(999), "999.00B")
	udwTest.Equal(udwStrconv.GbFromFloat64(1000), "0.977KB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1001), "0.978KB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1024), "1.000KB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1000*1000-1), "976.6KB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1000*1000), "0.954MB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1000*1000+1), "0.954MB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1024*1024), "1.000MB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1024*1024*1024), "1.000GB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1024*1024*1024*1024), "1.000TB")
	udwTest.Equal(udwStrconv.GbFromFloat64(1024*1024*1024*1024*1024), "1.000PB")
}

func TestGbPaddingFromInt64(t *testing.T) {
	udwTest.Equal(udwStrconv.GbPaddingFromInt64(1), "1.0000B")
	udwTest.Equal(udwStrconv.GbPaddingFromInt64(1000*1024), "0.977MB")
	udwTest.Equal(udwStrconv.GbPaddingFromInt64(1000*1000-1), "976.6KB")
	udwTest.Equal(udwStrconv.GbPaddingFromInt64(1000*1000), "0.954MB")
}

func TestGbStringToFloat64(t *testing.T) {
	f, errMsg := udwStrconv.GbstringToFloat64("16.00B")
	udwTest.Equal(errMsg, "")
	udwTest.Equal(f, float64(16))

	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00B"), float64(16))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00KB"), float64(16<<10))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00MB"), float64(16<<20))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00GB"), float64(16<<30))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00TB"), float64(16<<40))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("16.00PB"), float64(16<<50))
	udwTest.Equal(udwStrconv.GbStringToFloat64Default0("0.50GB"), float64(512<<20))

	for _, f := range []float64{
		1024 * 1024 * 1024,
		1000*1000 + 1,
		1024 * 1024 * 1024 * 1024 * 1024,
		0.5,
		1,
		1001,
		1000*1000 - 1,
	} {
		s := udwStrconv.GbFromFloat64(f)
		f1 := udwStrconv.GbStringToFloat64Default0(s)
		udwTest.Ok(math.Abs(f1-f) < 0.01*f)
	}
}
