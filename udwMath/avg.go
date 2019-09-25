package udwMath

import (
	"math"
	"sort"
)

func FloatAverage(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	}
	total := 0.0
	for i := 0; i < len(a); i++ {
		total += a[i]
	}
	return total / float64(len(a))
}

func IntMax(a []int) int {
	if len(a) == 0 {
		return 0
	}
	max := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

func FloatMax(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	}
	max := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

func FloatMin(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	}
	min := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

func IntMin(a ...int) int {
	if len(a) == 0 {
		return 0
	}
	min := a[0]
	for i := range a {
		if a[i] < min {
			min = a[i]
		}
	}
	return min
}

func FloatStdDev(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	}
	avg := FloatAverage(a)
	t1 := 0.0
	for i := 0; i < len(a); i++ {
		t := a[i] - avg
		t1 += t * t
	}
	t1 = t1 / float64(len(a))
	t1 = math.Sqrt(t1)
	return t1
}

func FloatMid(a []float64) float64 {
	if len(a) == 0 {
		return 0.0
	}
	b := make([]float64, len(a))
	copy(b, a)
	sort.Float64s(b)
	return b[len(b)/2]
}
