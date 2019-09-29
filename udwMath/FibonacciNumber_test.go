package udwMath

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

func TestNewFibonacciNumber(t *testing.T) {
	f := NewFibonacciNumber()
	expect := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765}
	for i := 0; i < len(expect); i++ {
		udwTest.Ok(f.Current() == expect[i])
		f.Next()
	}
}
