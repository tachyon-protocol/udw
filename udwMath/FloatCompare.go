package udwMath

import (
	"math"
	"math/big"
)

const Epsilon = 1e-10

func Float64LessThan(x float64, y float64) bool {
	return x < (y - Epsilon)
}
func Float64LessEqualThan(x float64, y float64) bool {
	return x < (y + Epsilon)
}
func Float64GreaterThan(x float64, y float64) bool {
	return x > (y + Epsilon)
}
func Float64GreaterEqualThan(x float64, y float64) bool {
	return x > (y - Epsilon)
}
func Float64Equal(x float64, y float64) bool {
	diff := x - y
	return math.Abs(diff) < Epsilon
}

func Float64Compare(x float64, y float64) int {
	var bigFloat1 = big.NewFloat(x)
	var bigFloat2 = big.NewFloat(y)
	return bigFloat1.Cmp(bigFloat2)
}
