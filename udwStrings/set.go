package udwStrings

import (
	"sort"
)

func GetIntersection(lists ...[]string) []string {
	if len(lists) < 2 {
		return nil
	}
	for _, l := range lists {
		if len(l) == 0 {
			return nil
		}
	}
	intersection := lists[0]
	for i := 1; i < len(lists); i++ {
		_intersection := []string{}
		for _, a := range intersection {
			for _, b := range lists[i] {
				if a == b {
					_intersection = append(_intersection, a)
				}
			}
		}
		if len(_intersection) == 0 {
			return nil
		}
		intersection = _intersection
	}
	return intersection
}

func GetIntersectionV3(lists ...[]string) []string {
	if len(lists) < 2 {
		return nil
	}
	for _, l := range lists {
		if len(l) == 0 {
			return nil
		}
	}
	if len(lists) >= 3 {
		sort.Slice(lists, func(i, j int) bool {
			return len(lists[i]) < len(lists[j])
		})
	}
	intersection := lists[0]
	for i := 1; i < len(lists); i++ {
		_intersection := []string{}
		for _, a := range intersection {
			for _, b := range lists[i] {
				if a == b {
					_intersection = append(_intersection, a)
				}
			}
		}
		if len(_intersection) == 0 {
			return nil
		}
		intersection = _intersection
	}
	return intersection
}

func GetIntersectionV4(lists ...[]string) []string {
	if len(lists) < 2 {
		return nil
	}
	for _, l := range lists {
		if len(l) == 0 {
			return nil
		}
	}
	noRepeat := map[string]bool{}
	for i, list := range lists {
		if i == 0 {
			for _, v := range list {
				noRepeat[v] = true
			}
			continue
		}
		for _, v := range list {
			if noRepeat[v] {
				continue
			}
			delete(noRepeat, v)
		}
	}
	intersection := make([]string, 0, len(noRepeat))
	for v := range noRepeat {
		intersection = append(intersection, v)
	}
	return intersection
}

func GetIntersectionV2(lists ...[]string) []string {
	if len(lists) < 2 {
		return nil
	}
	for _, l := range lists {
		if len(l) == 0 {
			return nil
		}
	}
	intersection := []string{}
	for _, a := range lists[0] {
		for _, b := range lists[1] {
			if a == b {
				intersection = append(intersection, a)
			}
		}
	}
	if len(intersection) == 0 {
		return nil
	}
	if len(lists[2:]) == 0 {
		return intersection
	}
	next := make([][]string, 0, len(lists)-1)
	next = append(next, intersection)
	for _, list := range lists[2:] {
		next = append(next, list)
	}
	return GetIntersectionV2(next...)
}

func IsEqualCheckOrder(lists ...[]string) bool {
	return isEqual(true, lists...)
}

func IsEqualIgnoreOrder(lists ...[]string) bool {
	return isEqual(false, lists...)
}

func isEqual(checkOrder bool, lists ...[]string) bool {
	if len(lists) < 2 {
		return true
	}
	first := lists[0]
	size := len(first)
	if !checkOrder {
		sort.Strings(first)
	}
	for _, l := range lists {
		if len(l) != size {
			return false
		}
	}
	for i, list := range lists {
		if i == 0 {
			continue
		}
		if !checkOrder {
			sort.Strings(list)
		}
		for j := 0; j < size; j++ {
			if list[j] != first[j] {
				return false
			}
		}
	}
	return true
}
