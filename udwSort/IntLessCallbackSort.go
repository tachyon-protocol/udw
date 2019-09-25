package udwSort

import "sort"

type IntLessCallbackSortT struct {
	Data     []int
	LessFunc func(a int, b int) bool
}

func (s *IntLessCallbackSortT) Len() int {
	return len(s.Data)
}

func (s *IntLessCallbackSortT) Less(i, j int) bool {
	return s.LessFunc(i, j)
}

func (s *IntLessCallbackSortT) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func (s *IntLessCallbackSortT) Sort() {
	sort.Sort(s)
}

func IntLessCallbackSort(Data []int, LessFunc func(a int, b int) bool) {
	sort.Sort(&IntLessCallbackSortT{Data: Data, LessFunc: LessFunc})
}

func ReverseStringSort(sList []string) {
	sort.Sort(ReverseStringSlice(sList))
}

type ReverseStringSlice []string

func (p ReverseStringSlice) Len() int           { return len(p) }
func (p ReverseStringSlice) Less(i, j int) bool { return p[i] > p[j] }
func (p ReverseStringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type StringIntPair struct {
	String string
	Int    int
}

type StringIntPairSortedByIntDesc []StringIntPair

func (l StringIntPairSortedByIntDesc) Len() int {
	return len(l)
}
func (l StringIntPairSortedByIntDesc) Less(i int, j int) bool {
	if l[i].Int != l[j].Int {
		return l[i].Int > l[j].Int
	}
	return l[i].String < l[j].String
}
func (l StringIntPairSortedByIntDesc) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}

func ReverseStringList(sList []string) {
	if len(sList) <= 1 {
		return
	}
	for i := 0; i < len(sList)/2; i++ {
		sList[i], sList[len(sList)-i-1] = sList[len(sList)-i-1], sList[i]
	}
}
