package udwTask

import (
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwClose"
	"sync"
)

const debugGetFastestTaskResultRunnerV2 = false

type getFastestTaskResultRunnerV2Element struct {
	doFn   func() (out interface{}, err error)
	closer func()
}
type GetFastestTaskResultRunnerV2 struct {
	list       []getFastestTaskResultRunnerV2Element
	locker     sync.Mutex
	isStartRun bool
	closer     udwClose.Closer
}

func NewGetFastestTaskResultRunnerV2() *GetFastestTaskResultRunnerV2 {
	tasker := &GetFastestTaskResultRunnerV2{}
	tasker.closer.AddOnClose(func() {
		tasker.locker.Lock()
		list := tasker.list
		tasker.list = nil
		tasker.locker.Unlock()
		for _, elem := range list {
			if elem.closer != nil {
				elem.closer()
			}
		}
		tasker.list = nil
	})
	return tasker
}

func (r *GetFastestTaskResultRunnerV2) GrowToSize(size int) {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.isStartRun {
		panic("[GetFastestTaskResultRunnerV2] GrowToSize after start run")
	}
	if r.closer.IsClose() {
		return
	}
	if len(r.list) == 0 && cap(r.list) < size {
		r.list = make([]getFastestTaskResultRunnerV2Element, 0, size)
	}
}

func (r *GetFastestTaskResultRunnerV2) Add(fn func() (out interface{}, err error)) {
	r.AddWithCloser(fn, nil)
}

func (r *GetFastestTaskResultRunnerV2) AddWithCloser(doFn func() (out interface{}, err error), closer func()) {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.isStartRun {
		panic("[GetFastestTaskResultRunnerV2] AddWithCloser after start run")
	}
	if r.closer.IsClose() {
		return
	}
	r.list = append(r.list, getFastestTaskResultRunnerV2Element{
		doFn:   doFn,
		closer: closer,
	})
}

func (r *GetFastestTaskResultRunnerV2) Close() {
	r.closer.Close()
}

func (r *GetFastestTaskResultRunnerV2) AddUpperCloser(closer *udwClose.Closer) {
	r.closer.AddUpperCloser(closer)
}

func (r *GetFastestTaskResultRunnerV2) RunIsNotCloseFastest(isNotCloseFastest bool) (out interface{}, err error) {
	r.locker.Lock()
	elementList := r.list
	if r.closer.IsClose() {
		r.locker.Unlock()
		return nil, errors.New("[GetFastestTaskResultRunnerV2] RunIsNotCloseFastest closed before start")
	}
	r.isStartRun = true
	r.locker.Unlock()
	elementNum := len(elementList)
	if elementNum == 0 {
		return nil, errors.New("[GetFastestTaskResultRunnerV2] list count 0")
	}
	finishChan := make(chan runnerResult, elementNum)
	for i, elem := range elementList {
		i := i
		elem := elem
		go func() {
			outI, err := elem.doFn()
			finishChan <- runnerResult{
				out:   outI,
				err:   err,
				index: i,
			}
		}()
	}
	for i := 0; i < elementNum; i++ {
		ret := <-finishChan
		err = ret.err
		if err == nil {
			for i, elem := range elementList {
				if isNotCloseFastest && i == ret.index {
					continue
				}
				if elem.closer != nil {
					elem.closer()
				}
			}
			r.locker.Lock()
			r.list = nil
			r.locker.Unlock()
			r.closer.Close()
			return ret.out, nil
		} else {
			if debugGetFastestTaskResultRunnerV2 {
				fmt.Println("error", "GetFastestTaskResult fail", err.Error())
			}
		}
	}
	return nil, err
}

func (r *GetFastestTaskResultRunnerV2) Run() (out interface{}, err error) {
	return r.RunIsNotCloseFastest(false)
}

func (r *GetFastestTaskResultRunnerV2) RunReturnBool() bool {
	_, err := r.Run()
	if err != nil {
		return false
	}
	return true
}
