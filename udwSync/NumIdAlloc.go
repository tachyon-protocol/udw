package udwSync

import (
	"strconv"
	"sync"
)

type NumIdAlloc struct {
	locker sync.Mutex
	m      []bool
}

func (alloc *NumIdAlloc) AllocId() int {
	alloc.locker.Lock()
	for i, v := range alloc.m {
		if v == false {
			alloc.m[i] = true
			alloc.locker.Unlock()
			return i + 1
		}
	}
	alloc.m = append(alloc.m, true)
	thisId := len(alloc.m)
	alloc.locker.Unlock()
	return thisId
}

func (alloc *NumIdAlloc) MustFreeId(id int) {
	alloc.locker.Lock()
	if id > len(alloc.m) || id <= 0 {
		alloc.locker.Unlock()
		panic("62k6np9uvp " + strconv.Itoa(id))
	}
	alloc.m[id-1] = false
	alloc.locker.Unlock()
}
