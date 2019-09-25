package udwMap

import (
	"github.com/tachyon-protocol/udw/udwSort"
)

func MapStringBoolToSortedSlice(m map[string]bool) []string {
	output := make([]string, len(m))
	i := 0
	for s := range m {
		output[i] = s
		i++
	}
	udwSort.SortString(output)
	return output
}

func StringListToMapStringBool(sList []string) map[string]bool {
	m := make(map[string]bool, len(sList))
	for _, s := range sList {
		m[s] = true
	}
	return m
}

func StringListNoRepeat(sList []string) []string {
	if len(sList) <= 1 {
		return sList
	}

	if len(sList) <= 20 {
		output := make([]string, 1, len(sList))
		output[0] = sList[0]
		for i := 1; i < len(sList); i++ {
			hasFound := false
			for j := 0; j < i; j++ {
				if sList[i] == sList[j] {
					hasFound = true
					break
				}
			}
			if hasFound == false {
				output = append(output, sList[i])
			}
		}
		return output
	} else {
		m := make(map[string]struct{}, len(sList))

		output := make([]string, 0, len(sList))
		for _, s := range sList {
			_, ok := m[s]
			if ok == false {

				output = append(output, s)
				m[s] = struct{}{}
			}
		}

		return output
	}

}
