package udwGoImport

import (
	"github.com/tachyon-protocol/udw/udwMap"
	"sort"
)

type MustMulitGetPackageImportResponse struct {
	absImportPathMapMap   map[string]map[string]struct{}
	absImportPathToDirMap map[string]string
}

func (resp MustMulitGetPackageImportResponse) GetAllLevelImportPathList(absImportPath string) (output []string) {
	m := resp.absImportPathMapMap[absImportPath]
	if len(m) == 0 {
		return nil
	}
	var visitor func(absImportPath string)
	seenImportPathSet := map[string]struct{}{}
	visitor = func(absImportPath string) {
		_, ok := seenImportPathSet[absImportPath]
		if ok {
			return
		}
		seenImportPathSet[absImportPath] = struct{}{}
		m := resp.absImportPathMapMap[absImportPath]
		for pkg := range m {
			visitor(pkg)
		}
	}
	visitor(absImportPath)
	return udwMap.SetStringToStringListAes(seenImportPathSet)
}

func (resp MustMulitGetPackageImportResponse) GetAllIncludeImportPathList() (output []string) {
	seenImportPathSet := map[string]struct{}{}
	for p1, m1 := range resp.absImportPathMapMap {
		seenImportPathSet[p1] = struct{}{}
		for p2 := range m1 {
			seenImportPathSet[p2] = struct{}{}
		}
	}
	return udwMap.SetStringToStringListAes(seenImportPathSet)
}
func (resp MustMulitGetPackageImportResponse) GetAllIncludeDirList() (output []string) {
	dirImportPathSet := map[string]struct{}{}
	for p1, m1 := range resp.absImportPathMapMap {
		dirImportPathSet[resp.absImportPathToDirMap[p1]] = struct{}{}
		for p2 := range m1 {
			dirImportPathSet[resp.absImportPathToDirMap[p2]] = struct{}{}
		}
	}
	return udwMap.SetStringToStringListAes(dirImportPathSet)
}

func (resp MustMulitGetPackageImportResponse) GetImportPathTreeMapSet() map[string]map[string]struct{} {
	return resp.absImportPathMapMap
}
func (resp MustMulitGetPackageImportResponse) GetDirectReferToPkgList(pkg string) []string {
	level1ReferPkgList := []string{}
	for referPkg, m1 := range resp.absImportPathMapMap {
		_, ok := m1[pkg]
		if ok {
			level1ReferPkgList = append(level1ReferPkgList, referPkg)
		}
	}
	sort.Strings(level1ReferPkgList)
	return level1ReferPkgList
}
