package udwTask

import (
	"fmt"
	"sync"
)

type ErrorTaskFn func() (err error)

type LimitThreadErrorHandleTaskManager struct {
	ErrorArr   []error
	errorMutex sync.Mutex
	threadNum  int
	retryNum   int
	taskChan   chan ErrorTaskFn
	wg         sync.WaitGroup
}

func NewLimitThreadErrorHandleTaskManager(threadNum int, retryNum int) *LimitThreadErrorHandleTaskManager {
	tm := &LimitThreadErrorHandleTaskManager{
		taskChan:  make(chan ErrorTaskFn),
		threadNum: threadNum,
		retryNum:  retryNum,
	}
	tm.run()
	return tm
}

func (m *LimitThreadErrorHandleTaskManager) run() {
	for i := 0; i < m.threadNum; i++ {
		go func() {
			for task := range m.taskChan {
				m.runOneTask(task)
			}
		}()
	}
}
func (m *LimitThreadErrorHandleTaskManager) runOneTask(task ErrorTaskFn) {
	defer m.wg.Done()
	retryNum := m.retryNum
Retry:
	err := task()
	if err != nil {
		fmt.Println("[LimitThreadErrorHandleTaskManager]", err)
		if retryNum >= 2 {
			retryNum--
			goto Retry
		}
	}
	if err != nil {
		m.errorMutex.Lock()
		m.ErrorArr = append(m.ErrorArr, err)
		m.errorMutex.Unlock()
	}
}
func (m *LimitThreadErrorHandleTaskManager) AddTask(task ErrorTaskFn) {
	m.wg.Add(1)
	m.taskChan <- task
}

func (m *LimitThreadErrorHandleTaskManager) Wait() {
	m.wg.Wait()
}

func (m *LimitThreadErrorHandleTaskManager) GetError() error {
	totalErrStr := ""
	for _, err := range m.ErrorArr {
		totalErrStr += err.Error() + "\n"
	}
	if totalErrStr != "" {
		return fmt.Errorf("%s", totalErrStr)
	}
	return nil
}
func (m *LimitThreadErrorHandleTaskManager) Close() {
	close(m.taskChan)
}
