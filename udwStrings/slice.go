package udwStrings

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwMath"
	"github.com/tachyon-protocol/udw/udwRand"
	"strings"
)

func ContainsAnyInSlice(toSearch string, KeywordSlice []string) bool {
	return IsInKeywordSlice(KeywordSlice, toSearch)
}

func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}

func IsInKeywordSlice(KeywordSlice []string, toSearch string) bool {
	for _, s := range KeywordSlice {
		if strings.Contains(toSearch, s) {
			return true
		}
	}
	return false
}

func StringSliceMerge(objList ...interface{}) []string {
	out := []string{}
	for _, objI := range objList {
		switch obj := objI.(type) {
		case string:
			out = append(out, obj)
		case []string:
			out = append(out, obj...)
		case nil:
		default:
			panic(fmt.Errorf("[StringSliceMerge] can only passed in string or []string got[%T]", objI))
		}
	}
	return out
}

func MergeNoRepeat(lists ...[]string) []string {
	size := 0
	for _, list := range lists {
		size += len(list)
	}
	out := make([]string, 0, size)
	noRepeatMap := map[string]struct{}{}
	for _, list := range lists {
		for _, s := range list {
			_, ok := noRepeatMap[s]
			if ok {
				continue
			}
			noRepeatMap[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

func StringSliceAverageSplit(slice []string, partNum int) (subSlice [][]string) {
	if partNum <= 0 {
		return [][]string{slice}
	}
	size := len(slice)
	if size == 0 {
		return nil
	}
	if partNum > size {
		return [][]string{slice}
	}
	step := udwMath.FloorToInt(float64(size) / float64(partNum))
	remainder := size % partNum
	start := 0
	for i := 1; i <= partNum; i++ {
		end := start + step
		if remainder > 0 {
			end++
			remainder--
		}
		subSlice = append(subSlice, slice[start:end])
		start = end
	}
	return subSlice
}

func AverageRandomIpSlice(ipSlice []string, partNum int, perPartCount int) (result []string) {
	if len(ipSlice) <= partNum*perPartCount {
		return ipSlice
	}
	for _, subSlice := range StringSliceAverageSplit(ipSlice, partNum) {
		if len(subSlice) <= perPartCount {
			result = append(result, subSlice...)
			continue
		}
		indexSlice := udwRand.Perm(len(subSlice))
		for _, index := range indexSlice[:perPartCount] {
			result = append(result, subSlice[index])
		}
	}
	return result
}

func SliceNoRepeatMerge(s1 []string, s2 []string) []string {
	result := append([]string{}, s1...)

	for _, s := range s2 {
		if !IsInSlice(result, s) {
			result = append(result, s)
		}
	}
	return result
}

func SliceNoRepeatAdd(slice []string, s string) []string {
	if IsInSlice(slice, s) {
		return slice
	}
	slice = append(slice, s)
	return slice
}

func SliceNoRepeat(slice []string) []string {
	outSlice := make([]string, 0, len(slice))
	m := map[string]struct{}{}
	for _, s := range slice {
		_, ok := m[s]
		if ok {
			continue
		}
		m[s] = struct{}{}
		outSlice = append(outSlice, s)
	}
	return outSlice
}

func StringSliceLastIndex(slice []string, search string) (pos int) {
	for pos = len(slice) - 1; pos >= 0; pos-- {
		if slice[pos] == search {
			return pos
		}
	}
	return -1
}

func StringSliceExcept(inList []string, exceptList []string) []string {
	outputList := []string{}
	for _, s := range inList {
		if IsInSlice(exceptList, s) == false {
			outputList = append(outputList, s)
		}
	}
	return outputList
}

func StringSliceClone(s []string) []string {
	if s == nil {
		return nil
	}
	outSlice := make([]string, len(s))
	for i := range s {
		outSlice[i] = s[i]
	}
	return outSlice
}

func ContainsIgnoreCase(s string, sub string) (ok bool) {
	lowS := strings.ToLower(s)
	lowSub := strings.ToLower(sub)
	return strings.Contains(lowS, lowSub)
}

func ContainsAnyIgnoreCase(s string, subList []string) (result map[string]int) {
	result = map[string]int{}
	lowS := strings.ToLower(s)
	for _, sub := range subList {
		lowSub := strings.ToLower(sub)
		n := strings.Count(lowS, lowSub)
		if n == 0 {
			continue
		}
		result[sub] = n
	}
	return result
}

func SliceRemove(slice []string, s string) []string {
	if !IsInSlice(slice, s) {
		return slice
	}
	if len(slice) == 0 {
		return []string{}
	}
	slice = append(slice, s)
	return slice
}
