package udwSync

import (
	"github.com/tachyon-protocol/udw/udwErr"
	"runtime"
	"sync"
)

const debug = false

func IsChannelClosed(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "send on closed channel"
}

func SendToChannel(send func()) {
	err := udwErr.PanicToError(send)
	if err != nil && !IsChannelClosed(err) {
		panic(err)
	}
}

type MemoryObjectPool struct {
	req          MemoryObjectPoolRequest
	objContainer chan interface{}

	lock               sync.Mutex
	allocedObjectCount int
	idleCount          int
}

type MemoryObjectPoolRequest struct {
	MaxEntryNumber int
	NewObjectFn    func() interface{}
}

func NewMemoryObjectPool(req MemoryObjectPoolRequest) *MemoryObjectPool {
	return &MemoryObjectPool{
		req:          req,
		objContainer: make(chan interface{}, req.MaxEntryNumber),
	}
}

func (pool *MemoryObjectPool) Go(job func(obj interface{})) {
	pool.allocateRoutine()
	obj := <-pool.objContainer
	job(obj)
	pool.objContainer <- obj
	pool.lock.Lock()
	pool.idleCount++
	pool.lock.Unlock()
}

func (pool *MemoryObjectPool) allocateRoutine() {
	pool.lock.Lock()
	if pool.idleCount > 0 {
		pool.idleCount--
		pool.lock.Unlock()
		return
	}
	if pool.allocedObjectCount >= pool.req.MaxEntryNumber {
		pool.lock.Unlock()
		return
	}
	pool.allocedObjectCount++
	pool.lock.Unlock()
	pool.objContainer <- pool.req.NewObjectFn()
}

type MemoryObjectPoolEntry struct {
	Obj interface{}
}

type MemoryObjectPoolWithCpuNumber struct {
	pool *MemoryObjectPool
	once sync.Once
}

func (pool2 *MemoryObjectPoolWithCpuNumber) Go(job func(entry *MemoryObjectPoolEntry)) {
	pool2.once.Do(func() {
		pool2.pool = NewMemoryObjectPool(MemoryObjectPoolRequest{
			MaxEntryNumber: runtime.NumCPU(),
			NewObjectFn: func() interface{} {
				return &MemoryObjectPoolEntry{}
			},
		})
	})
	pool2.pool.Go(func(obj interface{}) {
		job(obj.(*MemoryObjectPoolEntry))
	})
}
