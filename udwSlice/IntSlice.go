package udwSlice

func IntSliceRemoveAt(s *[]int, i int) {
	*s = append((*s)[:i], (*s)[i+1:]...)
}

func IntSliceRemove(s *[]int, v int) {
	thisLen := len(*s)
	for i := 0; i < thisLen; i++ {
		if (*s)[i] == v {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

func IsInIntSlice(s []int, t int) bool {
	for _, i := range s {
		if t == i {
			return true
		}
	}
	return false
}
