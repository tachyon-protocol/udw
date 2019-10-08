package udwSync

import (
	"github.com/tachyon-protocol/udw/udwTest"
	"sync"
	"testing"
)

func TestNumIdAlloc(t *testing.T) {
	var alloc NumIdAlloc
	var id int
	id = alloc.AllocId()
	udwTest.Equal(id, 1)
	id = alloc.AllocId()
	udwTest.Equal(id, 2)
	alloc.MustFreeId(1)
	alloc.MustFreeId(2)
	id = alloc.AllocId()
	udwTest.Equal(id, 1)
	alloc.MustFreeId(1)

	wg := sync.WaitGroup{}
	wg.Add(10)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	wg3 := sync.WaitGroup{}
	wg3.Add(10)
	m := map[int]int{}
	var mLocker sync.Mutex
	for i := 0; i < 10; i++ {
		go func() {
			id := alloc.AllocId()
			mLocker.Lock()
			m[id]++
			mLocker.Unlock()
			wg.Done()
			wg2.Wait()
			alloc.MustFreeId(id)
			wg3.Done()
		}()
	}
	wg.Wait()
	wg2.Done()
	wg3.Wait()
	udwTest.Equal(len(m), 10)
	{
		var alloc NumIdAlloc
		var id int
		id = alloc.AllocId()
		udwTest.Equal(id, 1)
		id = alloc.AllocId()
		udwTest.Equal(id, 2)
		alloc.MustFreeId(1)
		id = alloc.AllocId()
		udwTest.Equal(id, 1)
		id = alloc.AllocId()
		udwTest.Equal(id, 3)
	}
}
