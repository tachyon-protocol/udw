package udwMath

func SafeDivZeroFloat64(a float64, b float64) float64 {
	return a / b
}

func SafeDivZeroInt(a int, b int) int {
	if b == 0 {
		return -1
	}
	return a / b
}
