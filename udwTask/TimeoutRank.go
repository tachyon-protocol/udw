package udwTask

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"sync"
	"time"
)

type TimeoutRankTask struct {
	job    func() interface{}
	closer func()
}

type TimeoutRankManager struct {
	list []TimeoutRankTask
	req  NewTimeoutRankManagerRequest
}

type NewTimeoutRankManagerRequest struct {
	Timeout      time.Duration
	ThreadNumber int
	ResultLimit  int
}

func NewTimeoutRankManager(req NewTimeoutRankManagerRequest) *TimeoutRankManager {
	return &TimeoutRankManager{
		req:  req,
		list: make([]TimeoutRankTask, 0, req.ThreadNumber),
	}
}

func (manager *TimeoutRankManager) Add(job func() interface{}, closer func()) {
	manager.list = append(manager.list, TimeoutRankTask{
		job:    job,
		closer: closer,
	})
}

func (manager *TimeoutRankManager) Run() ([]interface{}, error) {
	wg := sync.WaitGroup{}
	chResultNumberDone := make(chan struct{})
	resultListLock := &sync.Mutex{}
	resultListReturn := false
	resultList := make([]interface{}, 0, manager.req.ResultLimit)
	errorListLock := &sync.Mutex{}
	errorListReturn := false
	errorList := make([]error, 0, manager.req.ThreadNumber)
	notCloseIndexListLock := &sync.Mutex{}
	notCloseIndexList := make([]int, 0, manager.req.ThreadNumber)
	for i, task := range manager.list {
		_i := i
		_job := task.job
		wg.Add(1)
		go func() {
			err := udwErr.PanicToError(func() {
				result := _job()
				resultListLock.Lock()
				if resultListReturn {
					resultListLock.Unlock()
					return
				}
				if len(resultList) >= manager.req.ResultLimit {
					close(chResultNumberDone)
					resultListLock.Unlock()
					return
				}
				resultList = append(resultList, result)
				resultListLock.Unlock()
				notCloseIndexListLock.Lock()
				notCloseIndexList = append(notCloseIndexList, _i)
				notCloseIndexListLock.Unlock()
			})
			wg.Done()
			if err == nil {
				return
			}
			errorListLock.Lock()
			if errorListReturn {
				errorListLock.Unlock()
				return
			}
			errorList = append(errorList, err)
			errorListLock.Unlock()
		}()
	}
	chAllDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(chAllDone)
	}()
	select {
	case <-chAllDone:
	case <-chResultNumberDone:
	case <-time.After(manager.req.Timeout):
	}
	for i, task := range manager.list {
		notClose := false
		notCloseIndexListLock.Lock()
		for _, notCloseI := range notCloseIndexList {
			if i == notCloseI {
				notClose = true
				break
			}
		}
		notCloseIndexListLock.Unlock()
		if notClose {
			continue
		}
		task.closer()
	}
	resultListLock.Lock()
	resultListReturn = true
	_resultList := resultList
	resultListLock.Unlock()
	if len(_resultList) > 0 {
		return _resultList, nil
	}
	errorListLock.Lock()
	errorListReturn = true
	_errString := ""
	for _, err := range errorList {
		_errString += err.Error() + "\n"
	}
	errorListLock.Unlock()
	return nil, errors.New(_errString)
}

func (manager *TimeoutRankManager) ForceCloseNotWaitForFinish() {
	for _, task := range manager.list {
		task.closer()
	}
}
