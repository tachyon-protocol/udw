package udwMath

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"math"
	"testing"
)

func TestSafeDivZeroFloat64(ot *testing.T) {
	udwTest.Equal(SafeDivZeroFloat64(1, 0), math.Inf(1))
	udwTest.Equal(SafeDivZeroFloat64(-1, 0), math.Inf(-1))
	udwTest.Ok(math.IsNaN(SafeDivZeroFloat64(0, 0)))
}

func TestSafeDivZeroInt(ot *testing.T) {
	udwTest.Equal(SafeDivZeroInt(1, 0), -1)
	udwTest.Equal(SafeDivZeroInt(-1, 0), -1)
	udwTest.Equal(SafeDivZeroInt(0, 0), -1)
}
