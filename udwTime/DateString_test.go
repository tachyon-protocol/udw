package udwTime

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestMustDateStringAddDay(t *testing.T) {
	udwTest.Equal(MustDateStringAddDay("2017-07-26", 1), "2017-07-27")
	udwTest.Equal(MustDateStringAddDay("2017-07-26", -1), "2017-07-25")

	udwTest.Equal(MustDateStringSubToDay("2017-07-26", "2017-07-25"), 1)
	udwTest.Equal(MustDateStringSubToDay("2017-07-26", "2017-07-27"), -1)
	udwTest.Equal(MustDateStringSubToDay("2017-07-26", "2017-07-26"), 0)
	udwTest.Equal(MustDateStringSubToDay("2017-07-26", "2017-07-01"), 25)
	udwTest.Equal(MustDateStringSubToDay("2017-07-26", "2017-06-30"), 26)
}
