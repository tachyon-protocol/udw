package udwMap

import (
	"github.com/tachyon-protocol/udw/udwSort"
)

type StringListSet struct {
	data []string
}

func (this *StringListSet) AppendString(s string) {
	if this.Contains(s) {
		return
	}

	this.data = append(this.data, s)
}

func (this *StringListSet) AppendStringSlice(ss []string) {
	for _, s := range ss {
		this.AppendString(s)
	}
}

func (this *StringListSet) Sort() {
	udwSort.SortString(this.data)
}

func (this *StringListSet) GetSliceCopy() []string {
	return append([]string{}, this.data...)
}

func (this *StringListSet) Contains(s string) bool {
	for _, thisS := range this.data {
		if thisS == s {
			return true
		}
	}
	return false
}

func (this StringListSet) Intersection(other StringListSet) StringListSet {
	result := []string{}
	for _, o := range other.data {
		if this.Contains(o) {
			result = append(result, o)
		}
	}
	return StringListSet{
		data: result,
	}
}

func (this StringListSet) Len() int {
	return len(this.data)
}

func (this StringListSet) At(idx int) (v string, exists bool) {
	if idx < 0 || idx > this.Len() {
		return "", false
	}
	return this.data[idx], true
}
