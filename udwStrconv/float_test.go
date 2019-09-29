package udwStrconv_test

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestFloat(ot *testing.T) {
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 1), "1")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 2), "1")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 3), "1.0")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1, 6), "1.0000")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(1.3333334, 6), "1.3333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(0.3333334, 6), "0.3333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(-0.3333334, 6), "-0.333")
	udwTest.Equal(udwStrconv.FormatFloat64ToFInLen(-0, 6), "0.0000")
}
