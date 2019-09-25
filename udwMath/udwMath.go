package udwMath

import (
	"math"
	"strconv"
)

func FloorToInt(x float64) int {
	return int(math.Floor(x))
}

func CeilToInt(x float64) int {
	return int(math.Ceil(x))
}

func RoundToInt(num float64) int {
	f := func(num float64) int {
		return int(num + math.Copysign(0.5, num))
	}
	output := math.Pow(10, float64(0))
	return (f(num * output)) / int(output)
}

func MustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
