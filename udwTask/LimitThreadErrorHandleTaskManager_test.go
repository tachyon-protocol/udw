package udwTask

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwTest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewLimitThreadErrorHandleTaskManager1(ot *testing.T) {
	tasker := NewLimitThreadErrorHandleTaskManager(1, 1)
	run_num := uint32(0)
	tasker.AddTask(func() (err error) {
		atomic.AddUint32(&run_num, 1)
		return fmt.Errorf("[TestNewLimitThreadErrorHandleTaskManager] fail")
	})
	tasker.Wait()
	udwTest.Equal(run_num, uint32(1))
	err := tasker.GetError()
	udwTest.Ok(err != nil)
	udwTest.Ok(strings.Contains(err.Error(), "[TestNewLimitThreadErrorHandleTaskManager]"))
}

func TestNewLimitThreadErrorHandleTaskManager2(ot *testing.T) {
	tasker := NewLimitThreadErrorHandleTaskManager(2, 2)
	run_num_locker := sync.Mutex{}
	run_num := 0
	tasker.AddTask(func() (err error) {
		run_num_locker.Lock()
		thisNum := run_num
		run_num++
		run_num_locker.Unlock()
		if thisNum < 1 {
			return fmt.Errorf("[TestNewLimitThreadErrorHandleTaskManager] fail")
		} else {
			return nil
		}
	})
	tasker.Wait()
	udwTest.Equal(run_num, 2)
	err := tasker.GetError()
	udwTest.Equal(err, nil)
}

func TestNewLimitThreadErrorHandleTaskManager3(ot *testing.T) {
	tasker := NewLimitThreadErrorHandleTaskManager(2, 2)
	run_num_locker := sync.Mutex{}
	run_num := 0
	fn := func() (err error) {
		run_num_locker.Lock()
		run_num++
		run_num_locker.Unlock()
		return nil
	}
	tasker.AddTask(fn)
	tasker.AddTask(fn)
	tasker.Wait()
	udwTest.Equal(run_num, 2)
	err := tasker.GetError()
	udwTest.Equal(err, nil)
}

func TestNewLimitThreadErrorHandleTaskManager4(ot *testing.T) {
	tasker := NewLimitThreadErrorHandleTaskManager(2, 2)
	run_num_locker := sync.Mutex{}
	run1_num := 0
	run2_num := 0
	tasker.AddTask(func() (err error) {
		run_num_locker.Lock()
		run1_num++
		run_num_locker.Unlock()
		return fmt.Errorf("[TestNewLimitThreadErrorHandleTaskManager] fail")
	})
	tasker.AddTask(func() (err error) {
		run_num_locker.Lock()
		run2_num++
		run_num_locker.Unlock()
		return fmt.Errorf("[TestNewLimitThreadErrorHandleTaskManager] fail")
	})
	tasker.Wait()
	udwTest.Equal(run1_num, 2)
	udwTest.Equal(run2_num, 2)
	err := tasker.GetError()
	udwTest.Ok(err != nil)
	udwTest.Ok(strings.Contains(err.Error(), "[TestNewLimitThreadErrorHandleTaskManager]"))
}
