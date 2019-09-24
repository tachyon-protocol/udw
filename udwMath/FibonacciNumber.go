package udwMath

type FibonacciNumber struct {
	before  int
	current int
}

func NewFibonacciNumber() *FibonacciNumber {
	return &FibonacciNumber{
		before:  0,
		current: 0,
	}
}

func (f *FibonacciNumber) Current() int {
	return f.current
}

func (f *FibonacciNumber) Next() int {
	if f.before == 0 && f.current == 0 {
		f.current = 1
		return 1
	}
	f.before, f.current = f.current, f.before+f.current
	return f.current
}
