package udwStrconv

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type StringPadding struct {
	max    int
	locker sync.Mutex
}

func (np *StringPadding) SetMaxBySlice(slice interface{}) {
	rv := reflect.ValueOf(slice)
	np.locker.Lock()
	np.max = len(strconv.Itoa(rv.Len()))
	np.locker.Unlock()
}

func (np *StringPadding) UpdateMax(name string) {
	np.locker.Lock()
	if len(name) > np.max {
		np.max = len(name)
	}
	np.locker.Unlock()
}

func (np *StringPadding) GetInt(i int) string {
	return np.Get(strconv.Itoa(i))
}

func (np *StringPadding) Get(name string) string {
	np.locker.Lock()
	delta := np.max - len(name)
	np.locker.Unlock()
	if delta <= 0 {
		return name
	}
	return name + strings.Repeat(" ", delta)
}
