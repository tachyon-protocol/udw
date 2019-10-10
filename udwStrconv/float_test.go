package udwStrconv_test

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestFloat(ot *testing.T) {
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 1), "1")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 2), "01")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 3), "1.0")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 6), "1.0000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1.3333334, 6), "1.3333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(0.3333334, 6), "0.3333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(-0.3333334, 6), "-0.333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(-0, 6), "0.0000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(-1, 6), "-1.000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(10000.1, 6), "010000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(010000.1, 6), "010000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1.6, 2), "02")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1.66, 3), "1.7")
}
