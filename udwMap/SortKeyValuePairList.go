package udwMap

func SortKeyValuePairList(data []KeyValuePair) {
	keyValuePair_quickSort(data, 0, len(data), keyValuePair_maxDepth(len(data)))
}

func keyValuePair_quickSort(data []KeyValuePair, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			keyValuePair_heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := keyValuePair_doPivot(data, a, b)

		if mlo-a < b-mhi {
			keyValuePair_quickSort(data, a, mlo, maxDepth)
			a = mhi
		} else {
			keyValuePair_quickSort(data, mhi, b, maxDepth)
			b = mlo
		}
	}
	if b-a > 1 {

		for i := a + 6; i < b; i++ {

			if data[i].Key < data[i-6].Key {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		keyValuePair_insertionSort(data, a, b)
	}
}

func keyValuePair_maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func keyValuePair_insertionSort(data []KeyValuePair, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (data[j].Key < data[j-1].Key); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func keyValuePair_medianOfThree(data []KeyValuePair, m1, m0, m2 int) {

	if data[m1].Key < data[m0].Key {
		data[m1], data[m0] = data[m0], data[m1]
	}

	if data[m2].Key < data[m1].Key {
		data[m2], data[m1] = data[m1], data[m2]

		if data[m1].Key < data[m0].Key {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}

}

func keyValuePair_doPivot(data []KeyValuePair, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {

		s := (hi - lo) / 8
		keyValuePair_medianOfThree(data, lo, lo+s, lo+2*s)
		keyValuePair_medianOfThree(data, m, m-s, m+s)
		keyValuePair_medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	keyValuePair_medianOfThree(data, lo, m, hi-1)

	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && (data[a].Key < data[pivot].Key); a++ {
	}
	b := a
	for {
		for ; b < c && !(data[pivot].Key < data[b].Key); b++ {
		}
		for ; b < c && (data[pivot].Key < data[c-1].Key); c-- {
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
		if !(data[pivot].Key < data[hi-1].Key) {
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !(data[b-1].Key < data[pivot].Key) {
			b--
			dups++
		}

		if !(data[m].Key < data[pivot].Key) {
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}

		protect = dups > 1
	}
	if protect {

		for {
			for ; a < b && !(data[b-1].Key < data[pivot].Key); b-- {
			}
			for ; a < b && (data[a].Key < data[pivot].Key); a++ {
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

func keyValuePair_heapSort(data []KeyValuePair, a, b int) {
	first := a
	lo := 0
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- {
		keyValuePair_siftDown(data, i, hi, first)
	}

	for i := hi - 1; i >= 0; i-- {

		data[first], data[first+i] = data[first+i], data[first]
		keyValuePair_siftDown(data, lo, i, first)
	}
}

func keyValuePair_siftDown(data []KeyValuePair, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}

		if child+1 < hi && (data[first+child].Key < data[first+child+1].Key) {
			child++
		}

		if !(data[first+root].Key < data[first+child].Key) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}
