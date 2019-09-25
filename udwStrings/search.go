package udwStrings

import "sort"

func IsInSliceBSearch(ss []string, s string) bool {
	idx := sort.SearchStrings(ss, s)
	if idx < 0 || idx >= len(ss) {
		return false
	}
	return ss[idx] == s
}
