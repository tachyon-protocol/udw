package udwTask

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwErr"
)

type runnerResult struct {
	out   interface{}
	err   error
	index int
}

func GetFastestTaskResult(fnList []func() (out interface{}, err error)) (out interface{}, err error) {
	num := len(fnList)
	finishChan := make(chan runnerResult, num)
	for _, fn := range fnList {
		fn := fn
		go func() {
			outI, err := fn()
			finishChan <- runnerResult{
				out: outI,
				err: err,
			}
		}()
	}
	for i := 0; i < num; i++ {
		ret := <-finishChan
		err = ret.err
		if err == nil {
			return ret.out, nil
		} else {
			if debugGetFastestTaskResultRunnerV2 {
				fmt.Println("error", "GetFastestTaskResult fail", err.Error())
			}
		}
	}
	return nil, err
}

type GetFastestTaskResultRunner []func() (out interface{}, err error)

func NewGetFastestTaskResultRunner() *GetFastestTaskResultRunner {
	return &GetFastestTaskResultRunner{}
}

func (r *GetFastestTaskResultRunner) Add(fn func() (out interface{}, err error)) {
	*r = append(*r, fn)
}
func (r *GetFastestTaskResultRunner) AddReturnBool(fn func() bool) {
	r.Add(func() (out interface{}, err error) {
		ret := fn()
		if ret {
			return nil, nil
		}
		return nil, udwErr.New("AddReturnBool fail")
	})
}

func (r *GetFastestTaskResultRunner) Run() (out interface{}, err error) {
	return GetFastestTaskResult(*r)
}

func (r *GetFastestTaskResultRunner) RunReturnBool() bool {
	_, err := GetFastestTaskResult(*r)
	if err != nil {
		return false
	}
	return true
}
