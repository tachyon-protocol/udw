package udwTask

import (
	"errors"
	"github.com/tachyon-protocol/udw/udwErr"
	"strings"
	"sync"
)

type limitThreadTaskManagerV3Job struct {
	fn     func()
	closer func()
}

type LimitThreadTaskManagerV3 []limitThreadTaskManagerV3Job

func NewLimitThreadTaskManagerV3(threadNumber int) *LimitThreadTaskManagerV3 {
	manager := &LimitThreadTaskManagerV3{}
	*manager = make([]limitThreadTaskManagerV3Job, 0, threadNumber)
	return manager
}

func (m *LimitThreadTaskManagerV3) Add(job func(), closer func()) {
	*m = append(*m, limitThreadTaskManagerV3Job{
		fn:     job,
		closer: closer,
	})
}

func (m *LimitThreadTaskManagerV3) ForceCloseNotWaitForFinish() {
	for _, job := range *m {
		if job.closer != nil {
			job.closer()
		}
	}
}

func (m *LimitThreadTaskManagerV3) Run() error {
	err := m.RunAndNotClose()
	m.ForceCloseNotWaitForFinish()
	return err
}

func (m *LimitThreadTaskManagerV3) RunAndNotClose() error {
	wg := &sync.WaitGroup{}
	errMsgSliceLocker := &sync.Mutex{}
	errMsgSlice := make([]string, 0, len(*m))
	for _, job := range *m {
		wg.Add(1)
		_job := job
		go func() {
			err := udwErr.PanicToError(_job.fn)
			if err != nil {
				errMsgSliceLocker.Lock()
				errMsgSlice = append(errMsgSlice, err.Error())
				errMsgSliceLocker.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if len(errMsgSlice) == 0 {
		return nil
	}
	return errors.New(strings.Join(errMsgSlice, "\n"))
}
