package udwRand

import "sort"

type PossibilityWeightRander struct {
	SearchList []float64
	Total      float64
}

func (p PossibilityWeightRander) ChoiceOne(r *UdwRand) (Index int) {
	if len(p.SearchList) == 0 {
		panic("[PossibilityWeightRander.ChoiceOne] len(p.SearchList)==0")
	}
	floatR := r.Float64Between(0, p.Total)
	return sort.SearchFloat64s(p.SearchList, floatR) - 1
}

func NewPossibilityWeightRander(weightList []float64) PossibilityWeightRander {
	pwl := PossibilityWeightRander{
		SearchList: make([]float64, len(weightList)),
	}
	for i, weight := range weightList {
		pwl.SearchList[i] = pwl.Total
		pwl.Total += weight
	}
	return pwl
}
