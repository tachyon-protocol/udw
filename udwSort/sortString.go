package udwSort

func SortString(data []string) {
	string_quickSort(data, 0, len(data), string_maxDepth(len(data)))
}

func string_quickSort(data []string, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			string_heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := string_doPivot(data, a, b)

		if mlo-a < b-mhi {
			string_quickSort(data, a, mlo, maxDepth)
			a = mhi
		} else {
			string_quickSort(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {

		for i := a + 6; i < b; i++ {

			if data[i] < data[i-6] {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		string_insertionSort(data, a, b)
	}
}

func string_maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func string_insertionSort(data []string, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (data[j] < data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func string_medianOfThree(data []string, m1, m0, m2 int) {

	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}

	if data[m2] < data[m1] {
		data[m2], data[m1] = data[m1], data[m2]

		if data[m1] < data[m0] {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}

}

func string_doPivot(data []string, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {

		s := (hi - lo) / 8
		string_medianOfThree(data, lo, lo+s, lo+2*s)
		string_medianOfThree(data, m, m-s, m+s)
		string_medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	string_medianOfThree(data, lo, m, hi-1)

	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && (data[a] < data[pivot]); a++ {
	}
	b := a
	for {
		for ; b < c && !(data[pivot] < data[b]); b++ {
		}
		for ; b < c && (data[pivot] < data[c-1]); c-- {
		}
		if b >= c {
			break
		}

		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}

	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {

		dups := 0
		if !(data[pivot] < data[hi-1]) {
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !(data[b-1] < data[pivot]) {
			b--
			dups++
		}

		if !(data[m] < data[pivot]) {
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}

		protect = dups > 1
	}
	if protect {

		for {
			for ; a < b && !(data[b-1] < data[pivot]); b-- {
			}
			for ; a < b && (data[a] < data[pivot]); a++ {
			}
			if a >= b {
				break
			}

			data[a], data[b-1] = data[b-1], data[a]
			a++
			b--
		}
	}

	data[pivot], data[b-1] = data[b-1], data[pivot]
	return b - 1, c
}

func string_heapSort(data []string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- {
		string_siftDown(data, i, hi, first)
	}

	for i := hi - 1; i >= 0; i-- {

		data[first], data[first+i] = data[first+i], data[first]
		string_siftDown(data, lo, i, first)
	}
}

func string_siftDown(data []string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}

		if child+1 < hi && (data[first+child] < data[first+child+1]) {
			child++
		}

		if !(data[first+root] < data[first+child]) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}
