package udwTask

import "sync"

type SimpleTaskManager struct {
	wg sync.WaitGroup
}

func (t *SimpleTaskManager) AddTask(task Task) {
	t.wg.Add(1)
	go func() {

		task.Run()

		t.wg.Done()
	}()
}

func (t *SimpleTaskManager) AddFunc(task func()) {
	t.wg.Add(1)
	go func() {

		task()

		t.wg.Done()
	}()
}

func (t *SimpleTaskManager) Wait() {
	t.wg.Wait()
}

func (t *SimpleTaskManager) Close() {
	t.Wait()
}

func NewSimpleTaskManager() *SimpleTaskManager {
	return &SimpleTaskManager{}
}

func RunTaskRepeat(f func(), num int) {
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}

func RunTaskRepeatWithLimitThread(f func(), taskNum int, threadNum int) {
	var wg sync.WaitGroup
	taskChan := make(chan func())
	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				task, ok := <-taskChan
				if ok == false {
					return
				}
				task()
				wg.Done()
			}
		}()
	}
	wg.Add(taskNum)
	for i := 0; i < taskNum; i++ {
		taskChan <- f
	}
	close(taskChan)
	wg.Wait()
}

func RunTask(threadNum int, funcList ...func()) {
	if threadNum > len(funcList) {
		threadNum = len(funcList)
	}
	var wg sync.WaitGroup
	taskChan := make(chan func())
	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				task, ok := <-taskChan
				if ok == false {
					return
				}
				task()
				wg.Done()
			}
		}()
	}
	wg.Add(len(funcList))
	for i := range funcList {
		taskChan <- funcList[i]
	}
	close(taskChan)
	wg.Wait()
}
