package udwMath

import "testing"
import (
	"github.com/tachyon-protocol/udw/udwTest"
)

func TestFloatCompare(ot *testing.T) {
	for i, testCase := range []struct {
		f      func(x float64, y float64) bool
		x      float64
		y      float64
		result bool
	}{
		{Float64LessThan, 1.0, 2.0, true},
		{Float64LessThan, 2.0, 1.0, false},
		{Float64LessThan, 2.0, 2.0, false},
		{Float64LessEqualThan, 1.0, 2.0, true},
		{Float64LessEqualThan, 1.0, 1.0, true},
		{Float64LessEqualThan, 2.0, 1.0, false},

		{Float64GreaterThan, 1.0, 2.0, false},
		{Float64GreaterThan, 2.0, 1.0, true},
		{Float64GreaterThan, 1.0, 1.0, false},
		{Float64Equal, 1.0, 1.0, true},
		{Float64Equal, 1.0, 2.0, false},

		{Float64Equal, 1.0 / 3.0 * 3.0, 1.0, true},
		{Float64GreaterEqualThan, 1.0, 2.0, false},
		{Float64GreaterEqualThan, 2.0, 1.0, true},
		{Float64GreaterEqualThan, 1.0, 1.0, true},
	} {
		udwTest.Equal(testCase.f(testCase.x, testCase.y), testCase.result,
			"fail at %d", i)
	}
}

func TestFloat64Compare(t *testing.T) {
	for i, testCase := range []struct {
		f      func(x float64, y float64) int
		x      float64
		y      float64
		result int
	}{
		{Float64Compare, 1.0, 2.0, -1},
		{Float64Compare, 2.0, 1.0, 1},
		{Float64Compare, 2.0, 2.0, 0},
		{Float64Compare, 0.0000003, 0.0000002, 1},
		{Float64Compare, 0.999999999, 1.0, -1},
	} {
		udwTest.Equal(testCase.f(testCase.x, testCase.y), testCase.result,
			"fail at %d", i)
	}
}
