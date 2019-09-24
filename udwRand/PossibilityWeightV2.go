package udwRand

import (
	"sort"
)

type PossibilityWeightRanderV2 struct {
	SearchList []float64
	Total      float64
	ResultList []string
}

func (p PossibilityWeightRanderV2) ChoiceOne(r *UdwRand) (result string) {
	if len(p.SearchList) == 0 {
		panic("[PossibilityWeightRanderV2.ChoiceOne] len(p.SearchList)==0")
	}
	floatR := r.Float64Between(0, p.Total)
	idx := sort.SearchFloat64s(p.SearchList, floatR) - 1
	return p.ResultList[idx]
}

func NewPossibilityWeightRanderV2(resultWeightMap map[string]float64) PossibilityWeightRanderV2 {
	pwl := PossibilityWeightRanderV2{
		SearchList: make([]float64, 0, len(resultWeightMap)),
		ResultList: make([]string, 0, len(resultWeightMap)),
	}
	for result, weight := range resultWeightMap {
		pwl.SearchList = append(pwl.SearchList, pwl.Total)
		pwl.ResultList = append(pwl.ResultList, result)
		pwl.Total += weight
	}
	return pwl
}
