package kkcflate

func sortbyLiteral(data []literalNode) {
	byLiteral_quickSort(data, 0, len(data), byLiteral_maxDepth(len(data)))
}

func byLiteral_quickSort(data []literalNode, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			byLiteral_heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := byLiteral_doPivot(data, a, b)

		if mlo-a < b-mhi {
			byLiteral_quickSort(data, a, mlo, maxDepth)
			a = mhi
		} else {
			byLiteral_quickSort(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {

		for i := a + 6; i < b; i++ {

			if data[i].literal < data[i-6].literal {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		byLiteral_insertionSort(data, a, b)
	}
}

func byLiteral_maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func byLiteral_insertionSort(data []literalNode, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (data[j].literal < data[j-1].literal); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func byLiteral_medianOfThree(data []literalNode, m1, m0, m2 int) {

	if data[m1].literal < data[m0].literal {
		data[m1], data[m0] = data[m0], data[m1]
	}

	if data[m2].literal < data[m1].literal {
		data[m2], data[m1] = data[m1], data[m2]

		if data[m1].literal < data[m0].literal {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}

}

func byLiteral_doPivot(data []literalNode, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {

		s := (hi - lo) / 8
		byLiteral_medianOfThree(data, lo, lo+s, lo+2*s)
		byLiteral_medianOfThree(data, m, m-s, m+s)
		byLiteral_medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	byLiteral_medianOfThree(data, lo, m, hi-1)

	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && (data[a].literal < data[pivot].literal); a++ {
	}
	b := a
	for {
		for ; b < c && !(data[pivot].literal < data[b].literal); b++ {
		}
		for ; b < c && (data[pivot].literal < data[c-1].literal); c-- {
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
		if !(data[pivot].literal < data[hi-1].literal) {
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !(data[b-1].literal < data[pivot].literal) {
			b--
			dups++
		}

		if !(data[m].literal < data[pivot].literal) {
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}

		protect = dups > 1
	}
	if protect {

		for {
			for ; a < b && !(data[b-1].literal < data[pivot].literal); b-- {
			}
			for ; a < b && (data[a].literal < data[pivot].literal); a++ {
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

func byLiteral_heapSort(data []literalNode, a, b int) {
	first := a
	lo := 0
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- {
		byLiteral_siftDown(data, i, hi, first)
	}

	for i := hi - 1; i >= 0; i-- {

		data[first], data[first+i] = data[first+i], data[first]
		byLiteral_siftDown(data, lo, i, first)
	}
}

func byLiteral_siftDown(data []literalNode, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}

		if child+1 < hi && (data[first+child].literal < data[first+child+1].literal) {
			child++
		}

		if !(data[first+root].literal < data[first+child].literal) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}
