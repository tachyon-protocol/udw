package udwTask

import (
	"sync"
)

type LimitMaxThreadTaskManager struct {
	task_chan       chan func()
	wg              sync.WaitGroup
	maxThread       int
	threadNum       int
	threadNumLocker sync.RWMutex
}

func New(num_thread int) *LimitMaxThreadTaskManager {
	if num_thread == 0 || num_thread <= -2 {
		panic("gnpkdxm4bn")
	}
	tm := &LimitMaxThreadTaskManager{
		maxThread: num_thread,
	}
	return tm
}

func (tm *LimitMaxThreadTaskManager) AddFunc(f func()) {
	tm.threadNumLocker.Lock()
	if tm.task_chan == nil {
		tm.task_chan = make(chan func(), 0)
	}
	if tm.threadNum < tm.maxThread {
		go func() {
			task_chan := tm.task_chan
			for {
				task, ok := <-task_chan
				if ok == false {
					return
				}
				task()
				tm.wg.Done()
			}
		}()
		tm.threadNum++
	} else if tm.maxThread == -1 {
		tm.threadNumLocker.Unlock()
		tm.wg.Add(1)
		defer tm.wg.Done()
		f()
		return
	}
	tm.threadNumLocker.Unlock()

	tm.wg.Add(1)
	tm.task_chan <- f
}

func (t *LimitMaxThreadTaskManager) AddFuncSync(f func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	t.AddFunc(func() {
		defer wg.Done()
		f()
	})
	wg.Wait()
}

func (tm *LimitMaxThreadTaskManager) Close() {
	tm.wg.Wait()
	tm.threadNumLocker.Lock()
	if tm.task_chan == nil {
		tm.task_chan = make(chan func(), 0)
	}
	tm.threadNumLocker.Unlock()
	close(tm.task_chan)

}

func (t *LimitMaxThreadTaskManager) WaitAndNotClose() {
	t.wg.Wait()
}
