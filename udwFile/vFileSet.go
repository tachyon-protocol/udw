package udwFile

import (
	"path/filepath"
	"sort"
)

type VFileItem struct {
	Path    string
	Content []byte
	Hash    string
}

func VFileSetSort(vfileSet []VFileItem) {
	sort.Sort(vFileItemSortType(vfileSet))
}
func VFileSetAddDirPath(basePath string, vfileSet []VFileItem) {
	for i := range vfileSet {
		vfileSet[i].Path = filepath.Join(basePath, vfileSet[i].Path)
	}
	return
}

func MustCheckContentAndWriteVFileItemList(basePath string, vfileSet []VFileItem) {
	for _, item := range vfileSet {
		MustCheckContentAndWriteFileWithMkdir(filepath.Join(basePath, item.Path), item.Content)
	}
}

type vFileItemSortType []VFileItem

func (t vFileItemSortType) Len() int           { return len(t) }
func (t vFileItemSortType) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t vFileItemSortType) Less(i, j int) bool { return t[i].Path < t[j].Path }
