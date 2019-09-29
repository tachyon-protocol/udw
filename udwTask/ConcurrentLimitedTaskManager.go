package udwTask

import (
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwClose"
	"github.com/tachyon-protocol/udw/udwMath"
	"math"
	"sync"
	"time"
)

const debugUdwTask = false

type ConcurrentLimitedTaskManager struct {
	taskList    []innerTask
	maxRunCount int
	isStarted   bool
	locker      sync.Mutex
	canceller   udwClose.Closer
}

type taskDoFunc func() (out interface{}, err error)

type innerTask struct {
	doFunc taskDoFunc
	closer func()
}

type innerTaskResult struct {
	innerTask
	timeCost time.Duration
	out      interface{}
	err      error
	index    int
}

func MustNewConcurrentLimitedTaskManager(limited int) *ConcurrentLimitedTaskManager {
	if limited <= 0 {
		panic("illegal parameter")
	}
	return &ConcurrentLimitedTaskManager{
		maxRunCount: limited,
	}
}

func (tm *ConcurrentLimitedTaskManager) AddTask(doFunc taskDoFunc, closer func()) {
	if tm.canceller.IsClose() || doFunc == nil {
		return
	}

	tm.locker.Lock()
	if tm.isStarted {
		tm.locker.Unlock()
		return
	}
	tm.taskList = append(tm.taskList, innerTask{
		doFunc: doFunc,
		closer: closer,
	})
	tm.locker.Unlock()
}

func (tm *ConcurrentLimitedTaskManager) FairRunGetFastestAndClose(isCloseFastest bool) (out interface{}, err error) {
	err = tm.checkAndMarkStartedBeforeRun()
	if err != nil {
		return nil, errors.New("4r7qp38ba3 " + err.Error())
	}
	tm.locker.Lock()
	taskList := tm.taskList
	tm.locker.Unlock()
	totalCount := len(taskList)
	stepCount := udwMath.CeilToInt(float64(totalCount) / float64(tm.maxRunCount))
	if debugUdwTask {
		fmt.Println(">>> ConcurrentLimitedTaskManager FairRunAndClose totalCount", totalCount, "maxRunCount", tm.maxRunCount, "STEP_NUMBER", stepCount)
	}

	taskResultFastest := innerTaskResult{
		err: errors.New("e3huaj8uv4 all error"),
	}
	for i := 0; i < stepCount; i++ {
		if tm.canceller.IsClose() {
			if debugUdwTask {
				fmt.Println(">>> ConcurrentLimitedTaskManager FairRunAndClose outer canceled 1")
			}

			return nil, errors.New("j9ffxay7df cancelled")
		}
		startOffset := i * tm.maxRunCount
		endOffset := int(math.Min(float64((i+1)*tm.maxRunCount), float64(totalCount)))
		subTaskList := taskList[startOffset:endOffset]
		_taskResult := tm.getFastestRun(subTaskList, isCloseFastest)
		if debugUdwTask {
			fmt.Println(">>> ConcurrentLimitedTaskManager FairRunAndClose taskList over fast cost: ", _taskResult.timeCost.String())
		}
		if _taskResult.err != nil {

			if debugUdwTask {
				fmt.Println("error", "a6hfhgng4f", _taskResult.err)
			}
			continue
		}
		if taskResultFastest.err != nil {

			taskResultFastest = _taskResult
			continue
		}
		if _taskResult.timeCost > 0 && _taskResult.timeCost < taskResultFastest.timeCost {

			slower := taskResultFastest
			slower.close()
			taskResultFastest = _taskResult
		} else {

			_taskResult.close()
		}
	}
	if tm.canceller.IsClose() {
		if debugUdwTask {
			fmt.Println(">>> ConcurrentLimitedTaskManager FairRunAndClose outer canceled 2")
		}
		return nil, errors.New("3np35e9gyq cancelled")
	}
	if isCloseFastest {
		if taskResultFastest.err == nil {
			taskResultFastest.close()
		}
	}
	tm.clearTaskList()
	return taskResultFastest.out, taskResultFastest.err
}

func (tm *ConcurrentLimitedTaskManager) FairRunAllAndClose(isCloseAll bool) error {
	err := tm.checkAndMarkStartedBeforeRun()
	if err != nil {

		return errors.New("5uq8ztuqf4 " + err.Error())
	}
	tm.locker.Lock()
	taskList := tm.taskList
	tm.locker.Unlock()
	totalCount := len(taskList)
	stepCount := udwMath.CeilToInt(float64(totalCount) / float64(tm.maxRunCount))
	for i := 0; i < stepCount; i++ {
		if tm.canceller.IsClose() {
			return errors.New("RunAndCloseWithResolveFunc outer canceled")
		}
		startOffset := i * tm.maxRunCount
		endOffset := int(math.Min(float64((i+1)*tm.maxRunCount), float64(totalCount)))
		subTaskList := taskList[startOffset:endOffset]

		tm.runWaitAll(subTaskList, isCloseAll)
	}
	if tm.canceller.IsClose() {
		return errors.New("outer canceled")
	}
	tm.clearTaskList()
	return nil
}

func (tm *ConcurrentLimitedTaskManager) UnfairRunAndClose(isCloseFastest bool) (out interface{}, err error) {
	err = tm.checkAndMarkStartedBeforeRun()
	if err != nil {
		return nil, errors.New("6r7cmy2ukz " + err.Error())
	}
	taskCount := len(tm.taskList)
	if debugUdwTask {
		fmt.Println(">>> ConcurrentLimitedTaskManager UnfairRunAndClose taskCount", taskCount, "maxRunCount", tm.maxRunCount)
	}
	taskLimitedCount := int(math.Min(float64(taskCount), float64(tm.maxRunCount)))
	if taskLimitedCount <= 0 {
		return nil, errors.New("dsxgd37225 taskLimitedCount <= 0")
	}
	tm.locker.Lock()
	taskList := tm.taskList
	tm.locker.Unlock()
	out, err = tm.getUnfairFastestAndClose(taskList, isCloseFastest)
	tm.clearTaskList()
	return
}

func (tm *ConcurrentLimitedTaskManager) getUnfairFastestAndClose(taskList []innerTask, isCloseFastest bool) (out interface{}, err error) {
	taskCount := len(taskList)
	taskLimitedCount := int(math.Min(float64(taskCount), float64(tm.maxRunCount)))
	taskBarrier := make(chan int, taskLimitedCount)
	finishChan := make(chan innerTaskResult, taskLimitedCount)
	closer := udwClose.NewCloser()
	closer.AddUpperCloser(&tm.canceller)
	go func() {
		for i, t := range tm.taskList {
			task := t
			idx := i
			select {
			case taskBarrier <- idx:
				go func() {
					if closer.IsClose() {
						if debugUdwTask {
							fmt.Println(">>> ConcurrentLimitedTaskManager getUnfairFastestAndClose task", idx, "canceled pwgrdeh84r")
						}
						return
					}
					out, err := task.doFunc()
					taskRet := innerTaskResult{
						index:     idx,
						innerTask: task,
						out:       out,
						err:       err,
					}
					if closer.IsClose() {
						if debugUdwTask {
							fmt.Println(">>> ConcurrentLimitedTaskManager getUnfairFastestAndClose task", idx, "cancel send result z9dncggw6b")
						}
						return
					}
					select {
					case finishChan <- taskRet:
						return
					case <-closer.GetCloseChan():
						if debugUdwTask {
							fmt.Println(">>> ConcurrentLimitedTaskManager getUnfairFastestAndClose task", idx, "canceled beause outer find success.")
						}
						return
					}
				}()
			case <-closer.GetCloseChan():
				if debugUdwTask {
					fmt.Println(">>> ConcurrentLimitedTaskManager getUnfairFastestAndClose we may find the fastest, so cancel")
				}
				return
			}
		}
	}()

	for i := 0; i < taskCount; i++ {
		select {
		case taskResult := <-finishChan:
			if taskResult.err == nil {

				closer.Close()
				for idx, task := range tm.taskList {
					if !isCloseFastest && idx == taskResult.index {
						continue
					}
					task.close()
				}
				return taskResult.out, nil
			} else {
				<-taskBarrier
			}
		case <-tm.canceller.GetCloseChan():

		}
	}

	return nil, errors.New("2cx9frkpa5 all task failed")
}

func (t *innerTask) close() {
	if t.closer != nil {
		t.closer()
	}
}

func (tm *ConcurrentLimitedTaskManager) Cancel() {
	if tm.canceller.IsClose() {
		return
	}
	tm.canceller.CloseWithCallback(func() {
		tm.locker.Lock()
		if tm.taskList == nil {
			tm.locker.Unlock()
			return
		}
		taskList := tm.taskList
		tm.taskList = nil
		tm.locker.Unlock()
		for _, task := range taskList {
			task.close()
		}
	})
}

func (tm *ConcurrentLimitedTaskManager) clearTaskList() {
	tm.locker.Lock()
	if tm.taskList != nil {
		tm.taskList = nil
		if debugUdwTask {
			fmt.Println("clearTaskList")
		}
	}
	tm.locker.Unlock()
}

func (tm *ConcurrentLimitedTaskManager) IsStarted() bool {
	tm.locker.Lock()
	is := tm.isStarted
	tm.locker.Unlock()
	return is
}

func (tm *ConcurrentLimitedTaskManager) checkAndMarkStartedBeforeRun() (err error) {
	if tm.canceller.IsClose() {
		return errors.New("aen4vu7ecv cancelled")
	}
	tm.locker.Lock()
	if tm.isStarted {
		tm.locker.Unlock()

		return errors.New("rfhkeuk57g started")
	}
	tm.isStarted = true
	tm.locker.Unlock()
	if len(tm.taskList) <= 0 {

		return errors.New("gbs5fzgfvz len(tm.taskList) <= 0")
	}
	return nil
}

func (tm *ConcurrentLimitedTaskManager) runWaitAll(taskList []innerTask, isCloseAll bool) {
	taskCount := len(taskList)
	taskChan := make(chan innerTaskResult, taskCount)
	for i, t := range taskList {
		if tm.canceller.IsClose() {
			return
		}
		idx := i
		_task := t
		go func() {
			mark := time.Now()
			out, err := _task.doFunc()
			timeCost := time.Now().Sub(mark)
			taskChan <- innerTaskResult{
				index:     idx,
				innerTask: _task,
				out:       out,
				err:       err,
				timeCost:  timeCost,
			}
		}()
	}
	for i := 0; i < taskCount; i++ {
		if tm.canceller.IsClose() {
			return
		}
		_runResult := <-taskChan
		if isCloseAll {
			_runResult.close()
		}
	}
}

func (tm *ConcurrentLimitedTaskManager) getFastestRun(taskList []innerTask, isCloseFastest bool) (_taskResult innerTaskResult) {
	taskCount := len(taskList)
	if debugUdwTask {
		fmt.Println("getFastestRun len:", len(taskList), isCloseFastest)
	}
	taskChan := make(chan innerTaskResult, taskCount)
	for i, t := range taskList {
		idx := i
		task := t
		go func() {
			mark := time.Now()
			out, err := task.doFunc()
			timeCost := time.Now().Sub(mark)
			taskChan <- innerTaskResult{
				index:     idx,
				innerTask: task,
				out:       out,
				err:       err,
				timeCost:  timeCost,
			}
		}()
	}
	for i := 0; i < taskCount; i++ {
		if tm.canceller.IsClose() {
			_taskResult.err = errors.New("getFastestRun outer canceled 2")
			return _taskResult
		}
		_runResult := <-taskChan
		if _runResult.err == nil {
			for idx, task := range taskList {
				if !isCloseFastest && idx == _runResult.index {
					continue
				}
				task.close()
			}
			return _runResult
		} else {

			_taskResult = _runResult
		}
	}
	return _taskResult
}
