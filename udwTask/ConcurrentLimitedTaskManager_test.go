package udwTask

import (
	"errors"
	"fmt"
	"github.com/tachyon-protocol/udw/udwSync"
	"github.com/tachyon-protocol/udw/udwTest"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTaskManager_New(t *testing.T) {
	udwTest.AssertPanic(func() {
		MustNewConcurrentLimitedTaskManager(0)
	})

	udwTest.AssertPanic(func() {
		MustNewConcurrentLimitedTaskManager(-1)
	})

	tm := MustNewConcurrentLimitedTaskManager(1)
	udwTest.Equal(tm != nil, true)
}

func TestTaskManager_AddNilTask(t *testing.T) {
	for i := 0; i < 3; i++ {
		tm := MustNewConcurrentLimitedTaskManager(1)
		tm.AddTask(nil, nil)

		var out interface{}
		var err error

		if i == 0 {
			out, err = tm.FairRunGetFastestAndClose(true)
			udwTest.Equal(out, nil)
			udwTest.Equal(err != nil, true)
		} else if i == 1 {
			out, err = tm.UnfairRunAndClose(true)
			udwTest.Equal(out, nil)
			udwTest.Equal(err != nil, true)
		} else if i == 2 {
			err = tm.FairRunAllAndClose(true)
			udwTest.Equal(err != nil, true)
		}
	}
}

func TestTaskManager_AddTask(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(1)
	tm.AddTask(func() (out interface{}, err error) {
		return 5, nil
	}, nil)
	out, err := tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(err, nil)
	udwTest.Equal(out, 5)
	udwTest.Equal(len(tm.taskList), 0)

	tm.AddTask(func() (out interface{}, err error) {
		return 5, nil
	}, nil)
	udwTest.Equal(len(tm.taskList), 0)
	out, err = tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(out, nil)
	udwTest.Equal(err != nil, true)
}

func TestTaskManager_FairRunGetFastestAndClose(t *testing.T) {

	tm := MustNewConcurrentLimitedTaskManager(1)
	count := 0
	for i := 0; i < 1000; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			count++
			return nil, errors.New("default error")
		}, nil)
	}
	_, err := tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(err != nil, true)
	udwTest.Equal(count, 1000)

	tm = MustNewConcurrentLimitedTaskManager(10)
	countLocker := sync.Mutex{}
	count = 0
	for i := 0; i < 1000; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			countLocker.Lock()
			count++
			countLocker.Unlock()
			return nil, errors.New("default error")
		}, nil)
	}
	_, err = tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(err != nil, true)
	udwTest.Equal(count, 1000)

	tm = MustNewConcurrentLimitedTaskManager(100)
	countLocker = sync.Mutex{}
	count = 0
	for i := 0; i < 100000; i++ {
		i := i
		tm.AddTask(func() (out interface{}, err error) {
			countLocker.Lock()
			count++
			countLocker.Unlock()
			if i == 2 {
				return 2, nil
			}
			return nil, errors.New("default error")
		}, nil)
	}
	out, err := tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(out, 2)
	udwTest.Equal(err, nil)
	udwTest.Equal(count, 100000)

	tm = MustNewConcurrentLimitedTaskManager(5)
	closeCount := 0
	closeCountLock := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		i := i
		tm.AddTask(func() (out interface{}, err error) {
			return nil, errors.New("error, " + strconv.Itoa(i))
		}, func() {
			closeCountLock.Lock()
			closeCount++
			closeCountLock.Unlock()
		})
	}
	_, err = tm.FairRunGetFastestAndClose(true)
	udwTest.Equal(err != nil, true)
	udwTest.Equal(closeCount, 0)
}

func TestConcurrentLimitedTaskManager_CancelFairRunGetFastest(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	runCount := udwSync.NewInt(0)
	cancelCount := udwSync.NewInt(0)
	for i := 0; i < 100; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			time.Sleep(time.Millisecond * 100)
			runCount.Add(1)

			return nil, nil
		}, func() {
			cancelCount.Add(1)

		})
	}
	doneWg := sync.WaitGroup{}
	doneWg.Add(1)
	go func() {
		tm.FairRunGetFastestAndClose(true)
		doneWg.Done()
	}()
	time.Sleep(time.Millisecond * 50)
	go func() {
		tm.Cancel()
	}()

	doneWg.Wait()
	fmt.Println("done runCount", runCount.Get(), "cancelCount", cancelCount.Get())
	udwTest.Equal(runCount.Get() < 100, true)
	udwTest.Equal(cancelCount.Get() >= 100, true)
}

func TestConcurrentLimitedTaskManager_UnfairRunAndClose(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	taskRunCount := udwSync.NewInt(0)
	for i := 0; i < 5; i++ {
		i := i
		tm.AddTask(func() (out interface{}, err error) {
			taskRunCount.Add(1)

			ret := i%2 + 1
			time.Sleep(time.Duration(ret) * time.Millisecond * 10)
			return ret, nil
		}, nil)
	}

	_, err := tm.UnfairRunAndClose(true)
	udwTest.Equal(err, nil)
	udwTest.Equal(taskRunCount.Get(), 2)
}

func TestConcurrentLimitedTaskManager_UnfairRunAndClose3(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	taskRunCount := 0
	locker := sync.Mutex{}
	for i := 0; i < 5; i++ {
		i := i
		tm.AddTask(func() (out interface{}, err error) {
			locker.Lock()
			taskRunCount++
			locker.Unlock()
			ret := i%2 + 1
			time.Sleep(time.Duration(ret) * time.Millisecond * 50)
			err = errors.New(fmt.Sprintln("task", i, "error", out))
			return ret, err
		}, nil)
	}

	go func() {
		tm.UnfairRunAndClose(true)
	}()
	time.Sleep(time.Millisecond * 100)
	tm.Cancel()
}

func TestConcurrentLimitedTaskManager_FairRunAllAndClose(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	taskRunCount := 0
	cancelCount := 0
	locker := sync.Mutex{}
	for i := 0; i < 5; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			locker.Lock()
			taskRunCount++
			locker.Unlock()
			return
		}, func() {
			locker.Lock()
			cancelCount++
			locker.Unlock()
		})
	}
	tm.FairRunAllAndClose(true)
	udwTest.Equal(taskRunCount, 5)
	udwTest.Equal(cancelCount, 5)
}

func TestConcurrentLimitedTaskManager_FairRunAllAndClose2(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	taskRunCount := 0
	cancelCount := 0
	locker := sync.Mutex{}
	for i := 0; i < 5; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			locker.Lock()
			taskRunCount++
			locker.Unlock()
			time.Sleep(time.Millisecond * 50)
			return
		}, func() {
			locker.Lock()
			cancelCount++
			locker.Unlock()
		})
	}
	go func() {
		tm.FairRunAllAndClose(true)
	}()
	time.Sleep(time.Millisecond * 10)
	tm.Cancel()

	udwTest.Equal(taskRunCount < 5, true)
}

func TestConcurrentLimitedTaskManager_Cancel(t *testing.T) {
	tm := MustNewConcurrentLimitedTaskManager(2)
	for i := 0; i < 5; i++ {
		tm.AddTask(func() (out interface{}, err error) {
			return
		}, nil)
	}

	udwTest.Equal(len(tm.taskList), 5)
	tm.FairRunAllAndClose(true)
	udwTest.Equal(tm.taskList, nil)
	udwTest.Equal(len(tm.taskList), 0)
	tm.Cancel()

	for _, item := range tm.taskList {
		item.close()
	}

}
