package udwTask

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwClose"
	"sync"
)

type LimitThreadTaskManagerV2 struct {
	task_chan        chan func()
	num_thread       int
	threadCloserList []*udwClose.Closer
	numThreadLocker  sync.Mutex
	wg               sync.WaitGroup
}

func NewLimitThreadTaskManagerV2(num_thread int) *LimitThreadTaskManagerV2 {
	tm := &LimitThreadTaskManagerV2{
		task_chan: make(chan func(), 100),
	}
	tm.SetThreadNum(num_thread)
	return tm
}

func NewLimitThreadTaskManagerV2NoBuffer(num_thread int) *LimitThreadTaskManagerV2 {
	tm := &LimitThreadTaskManagerV2{
		task_chan: make(chan func(), 0),
	}
	tm.SetThreadNum(num_thread)

	return tm
}

func (t *LimitThreadTaskManagerV2) SetThreadNum(num_thread int) {
	t.numThreadLocker.Lock()
	defer t.numThreadLocker.Unlock()
	if num_thread < 0 {
		panic(fmt.Errorf("[SetThreadNum] num_thread[%d]<0", num_thread))
	}
	if t.num_thread == num_thread {
		return
	} else if t.num_thread > num_thread {
		for i := num_thread; i < t.num_thread; i++ {
			t.threadCloserList[i].Close()
		}
		t.threadCloserList = t.threadCloserList[:num_thread]
	} else if t.num_thread < num_thread {
		diffNum := num_thread - t.num_thread
		t.threadCloserList = append(t.threadCloserList, make([]*udwClose.Closer, diffNum)...)
		task_chan := t.task_chan
		for i := t.num_thread; i < num_thread; i++ {
			threadId := i
			closer := &udwClose.Closer{}
			t.threadCloserList[threadId] = closer
			go func() {
				closerChan := closer.GetCloseChan()
				for {
					select {
					case task, ok := <-task_chan:
						if ok == false {
							return
						}
						task()
						t.wg.Done()
					case <-closerChan:
						return
					}
				}
			}()
		}
	}
	t.num_thread = num_thread
}

func (t *LimitThreadTaskManagerV2) AddFunc(f func()) {

	t.wg.Add(1)
	t.task_chan <- f
}

func (t *LimitThreadTaskManagerV2) Close() {
	t.wg.Wait()
	close(t.task_chan)
	t.numThreadLocker.Lock()
	for _, closer := range t.threadCloserList {
		closer.Close()
	}
	t.numThreadLocker.Unlock()
}

func (t *LimitThreadTaskManagerV2) CloseAsync() {
	go t.Close()
}

func (t *LimitThreadTaskManagerV2) WaitAndNotClose() {
	t.wg.Wait()
}
