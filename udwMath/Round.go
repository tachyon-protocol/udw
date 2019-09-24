package udwMath

import (
	"math"
	"strconv"
)

func Float64RoundToRelativePrec(f float64, prec int) float64 {
	s := strconv.FormatFloat(f, 'e', prec, 64)
	o, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return o
}

func Float64RoundToRelativePrecOnOneFloat(f float64, prec int, precBaseF float64) float64 {
	if precBaseF == 0 {
		return 0
	}
	absPrec := math.Floor(math.Log10(precBaseF)) - float64(prec)
	return f - math.Mod(f, math.Pow10(int(absPrec)))
}
