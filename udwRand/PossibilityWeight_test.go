package udwRand

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestPossibilityWeight(t *testing.T) {
	r := NewInt64SeedUdwRand(0)
	rander := NewPossibilityWeightRander([]float64{1e-20})
	for i := 0; i < 100; i++ {
		udwTest.Equal(rander.ChoiceOne(r), 0)
	}
	rander = NewPossibilityWeightRander([]float64{1, 2, 3, 4})
	for i := 0; i < 100; i++ {
		ret := rander.ChoiceOne(r)
		udwTest.Ok(ret >= 0)
		udwTest.Ok(ret <= 3)
	}
	rander = NewPossibilityWeightRander([]float64{1, 0, 3, 0, 1})
	for i := 0; i < 100; i++ {
		ret := rander.ChoiceOne(r)
		udwTest.Ok(ret >= 0)
		udwTest.Ok(ret <= 4)
	}
}
