package udwMap

import (
	"github.com/tachyon-protocol/udw/udwSort"
)

func SetStringToStringListAes(m map[string]struct{}) []string {
	if len(m) == 0 {
		return nil
	}
	pairList := make([]string, len(m))
	index := 0
	for s := range m {
		pairList[index] = s
		index++
	}
	udwSort.SortString(pairList)
	return pairList
}

func StringListToSetString(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for _, s := range list {
		m[s] = struct{}{}
	}
	return m
}

func SetStringAddStringList(m map[string]struct{}, list []string) {
	for _, s := range list {
		m[s] = struct{}{}
	}
}

func SetStringMapEqualStringList(setStringMap map[string]struct{}, stringList []string) bool {
	if len(setStringMap) != len(stringList) {
		return false
	}
	for _, s := range stringList {
		_, ok := setStringMap[s]
		if ok == false {
			return false
		}
	}
	m := StringListToSetString(stringList)
	for s := range setStringMap {
		_, ok := m[s]
		if ok == false {
			return false
		}
	}
	return true
}
