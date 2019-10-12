package udwBuild

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"sync"
)

var gTagSet map[string]struct{}
var gTagListLocker sync.Mutex

func TagAdd(tagList ...string) {
	gTagListLocker.Lock()
	if gTagSet == nil {
		gTagSet = map[string]struct{}{}
	}
	for _, tag := range tagList {
		gTagSet[tag] = struct{}{}
	}
	gTagListLocker.Unlock()
}
func TagHas(tag string) bool {
	gTagListLocker.Lock()
	if gTagSet == nil {
		gTagSet = map[string]struct{}{}
	}
	_, has := gTagSet[tag]
	gTagListLocker.Unlock()
	return has
}
func TagGetList() []string {
	gTagListLocker.Lock()
	if gTagSet == nil {
		gTagSet = map[string]struct{}{}
	}
	outList := udwMap.SetStringToStringListAes(gTagSet)
	gTagListLocker.Unlock()
	return outList
}
