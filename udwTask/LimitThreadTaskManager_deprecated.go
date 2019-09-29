package udwTask

import "sync"

type LimitThreadTaskManager struct {
	task_chan       chan Task
	num_thread      int
	wg              sync.WaitGroup
	oneThreadLocker sync.Mutex
}

func NewLimitThreadTaskManager(num_thread int) *LimitThreadTaskManager {
	bufSize := 100
	if num_thread == 1 {
		bufSize = 0
	}
	tm := &LimitThreadTaskManager{
		task_chan:  make(chan Task, bufSize),
		num_thread: num_thread,
	}
	tm.run()
	return tm
}

func (t *LimitThreadTaskManager) run() {
	if t.num_thread == 1 {
		return
	}
	for i := 0; i < t.num_thread; i++ {
		go func() {
			for {
				task, ok := <-t.task_chan
				if ok == false {
					return
				}
				task.Run()
				t.wg.Done()
			}
		}()
	}
}

func (t *LimitThreadTaskManager) AddTask(task Task) {
	t.wg.Add(1)
	if t.num_thread == 1 {
		defer t.wg.Done()
		t.oneThreadLocker.Lock()
		defer t.oneThreadLocker.Unlock()
		task.Run()
		return
	}
	t.task_chan <- task
}

func (t *LimitThreadTaskManager) AddFunc(f func()) {
	t.AddTask(TaskFunc(f))
}

func (t *LimitThreadTaskManager) Wait() {
	t.wg.Wait()
}

func (t *LimitThreadTaskManager) WaitAndClose() {
	t.Close()
}

func (t *LimitThreadTaskManager) Close() {
	defer close(t.task_chan)
	t.Wait()
}

func (t *LimitThreadTaskManager) AddTaskNewThread(task Task) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		task.Run()
	}()
}
