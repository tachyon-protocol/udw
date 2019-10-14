package kkcflate

func sortbyFreq(data []literalNode) {
	byFreq_quickSort(data, 0, len(data), byFreq_maxDepth(len(data)))
}

func byFreq_quickSort(data []literalNode, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			byFreq_heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := byFreq_doPivot(data, a, b)

		if mlo-a < b-mhi {
			byFreq_quickSort(data, a, mlo, maxDepth)
			a = mhi
		} else {
			byFreq_quickSort(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {

		for i := a + 6; i < b; i++ {

			if (data[i].freq == data[i-6].freq && data[i].literal < data[i-6].literal) || (data[i].freq < data[i-6].freq) {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		byFreq_insertionSort(data, a, b)
	}
}

func byFreq_maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func byFreq_insertionSort(data []literalNode, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && ((data[j].freq == data[j-1].freq && data[j].literal < data[j-1].literal) || (data[j].freq < data[j-1].freq)); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func byFreq_medianOfThree(data []literalNode, m1, m0, m2 int) {

	if (data[m1].freq == data[m0].freq && data[m1].literal < data[m0].literal) || (data[m1].freq < data[m0].freq) {
		data[m1], data[m0] = data[m0], data[m1]
	}

	if (data[m2].freq == data[m1].freq && data[m2].literal < data[m1].literal) || (data[m2].freq < data[m1].freq) {
		data[m2], data[m1] = data[m1], data[m2]

		if (data[m1].freq == data[m0].freq && data[m1].literal < data[m0].literal) || (data[m1].freq < data[m0].freq) {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}

}

func byFreq_doPivot(data []literalNode, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {

		s := (hi - lo) / 8
		byFreq_medianOfThree(data, lo, lo+s, lo+2*s)
		byFreq_medianOfThree(data, m, m-s, m+s)
		byFreq_medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	byFreq_medianOfThree(data, lo, m, hi-1)

	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && ((data[a].freq == data[pivot].freq && data[a].literal < data[pivot].literal) || (data[a].freq < data[pivot].freq)); a++ {
	}
	b := a
	for {
		for ; b < c && !((data[pivot].freq == data[b].freq && data[pivot].literal < data[b].literal) || (data[pivot].freq < data[b].freq)); b++ {
		}
		for ; b < c && ((data[pivot].freq == data[c-1].freq && data[pivot].literal < data[c-1].literal) || (data[pivot].freq < data[c-1].freq)); c-- {
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
		if !((data[pivot].freq == data[hi-1].freq && data[pivot].literal < data[hi-1].literal) || (data[pivot].freq < data[hi-1].freq)) {
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !((data[b-1].freq == data[pivot].freq && data[b-1].literal < data[pivot].literal) || (data[b-1].freq < data[pivot].freq)) {
			b--
			dups++
		}

		if !((data[m].freq == data[pivot].freq && data[m].literal < data[pivot].literal) || (data[m].freq < data[pivot].freq)) {
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}

		protect = dups > 1
	}
	if protect {

		for {
			for ; a < b && !((data[b-1].freq == data[pivot].freq && data[b-1].literal < data[pivot].literal) || (data[b-1].freq < data[pivot].freq)); b-- {
			}
			for ; a < b && ((data[a].freq == data[pivot].freq && data[a].literal < data[pivot].literal) || (data[a].freq < data[pivot].freq)); a++ {
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

func byFreq_heapSort(data []literalNode, a, b int) {
	first := a
	lo := 0
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- {
		byFreq_siftDown(data, i, hi, first)
	}

	for i := hi - 1; i >= 0; i-- {

		data[first], data[first+i] = data[first+i], data[first]
		byFreq_siftDown(data, lo, i, first)
	}
}

func byFreq_siftDown(data []literalNode, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}

		if child+1 < hi && ((data[first+child].freq == data[first+child+1].freq && data[first+child].literal < data[first+child+1].literal) || (data[first+child].freq < data[first+child+1].freq)) {
			child++
		}

		if !((data[first+root].freq == data[first+child].freq && data[first+root].literal < data[first+child].literal) || (data[first+root].freq < data[first+child].freq)) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}
